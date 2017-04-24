package install

import (
	"errors"
	"fmt"
	"github.com/serkansipahi/app-decorators-cli/util/exec"
	"github.com/serkansipahi/app-decorators-cli/util/json"
	"github.com/serkansipahi/app-decorators-cli/util/os"
	"io/ioutil"
	"path"
	"path/filepath"
)

var (
	ErrNoModuleName        = errors.New("Failed: Please set module name e.g. 'appdec init --name=mymodule'")
	ErrAppPathExists       = errors.New("Failed: Apppath exists")
	ErrCantChangeToApppath = errors.New("Failed: cant change to apppath")
	ErrCantInstallDeps     = errors.New("Failed: someting gone wrong while installing dependencies")
	ErrWhileCleanup        = errors.New("Failed: someting gone wrong while cleaning")
)

type Config struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type AppExists interface {
	os.Chdir
}

type CreateApp interface {
	os.Mkdir
	os.Chdir
}

type Deps interface {
	exec.Execers
}

type Clean interface {
	os.Remover
}

type CopyCore interface {
	os.Remover
	os.ReadFiler
	os.CopyFiler
}

func New(name string, rootPath string, version string, cliName string, debug bool) *Install {

	return &Install{
		name,
		filepath.Join(rootPath, version),
		rootPath,
		cliName,
		version,
		debug,
	}
}

type Install struct {
	Name     string
	AppPath  string
	RootPath string
	CliName  string
	Version  string
	Debug    bool
}

func (i Install) Run() error {

	return i.Install(
		os.Os{},
		json.New(),
		exec.New(false, i.Debug),
	)
}

func (i Install) AppPathExists(appPath string, os AppExists) error {

	_, name := path.Split(appPath)

	if err := os.Chdir(appPath); err == nil {
		err = errors.New(fmt.Sprintf("\n"+
			"Run: '%s' already created\n"+
			"You can delete it with 'appdec delete --name=%s\n"+
			"", name, name))

		return err
	}

	return nil
}

func (i Install) CreateAppPath(appPath string, os CreateApp) error {

	var err error

	if err = os.Mkdir(appPath, 0755); err != nil {
		return ErrAppPathExists
	}
	if err = os.Chdir(appPath); err != nil {
		return ErrCantChangeToApppath
	}

	return nil
}

func (i Install) InstallDeps(commands []string, exec Deps) error {

	var err error

	err = exec.Run(commands)
	if err != nil {
		return ErrCantInstallDeps
	}

	return nil
}

func (i Install) Cleanup(appPath string, os Clean) error {

	var err error

	pkgJsonPath := filepath.Join(appPath, "package.json")
	if err = os.Remove(pkgJsonPath); err != nil {
		return ErrWhileCleanup
	}

	return nil
}

func (i Install) CreateAppJson(appPath string, json json.Stringifyer) error {

	config := Config{
		Name:    i.Name,
		Version: i.Version,
	}
	jsonData, err := json.Stringify(config)
	if err != nil {
		return err
	}

	appDecJsonPath := filepath.Join(appPath, i.CliName+".json")
	if err = ioutil.WriteFile(appDecJsonPath, jsonData, 0755); err != nil {
		return err
	}

	return nil
}

func (i Install) CopyCoreFiles(os CopyCore) error {

	appDecoratorPath := filepath.Clean(
		filepath.Join(i.RootPath, i.Name, "node_modules", "app-decorators"),
	)
	files, _ := os.ReadFiles(appDecoratorPath)

	for _, file := range files {

		src := filepath.Join(appDecoratorPath, file.Name())
		dist := filepath.Join(i.RootPath, i.Name, file.Name())

		err := os.CopyFile(src, dist)
		if err != nil {
			return err
		}
	}

	return nil
}

func (i Install) Install(os os.Os, json json.Stringifyer, exec Deps) error {

	var (
		err     error
		name    string = i.Name
		appPath string = i.RootPath + "/" + name
		// npm packages
		appDecPkg   string = "app-decorators@" + i.Version
		babelCliPkg string = "babel-cli@6.24.1"
	)

	if name == "" {
		return ErrNoModuleName
	}

	// Todo: create first app specific json

	// Return when  "appPath" exists
	if err = i.AppPathExists(appPath, os); err != nil {
		return err
	}

	// Create "appPath" if not exists
	if err := i.CreateAppPath(appPath, os); err != nil {
		return err
	}

	// Install dependencies
	err = i.InstallDeps([]string{
		"npm init -y",
		"npm install " + appDecPkg,
		"npm install " + babelCliPkg,
	}, exec)
	if err != nil {
		return ErrCantInstallDeps
	}

	// Cleanup
	if err = i.Cleanup(appPath, os); err != nil {
		return ErrWhileCleanup
	}

	// Create app specific json file
	fmt.Println("Run: create " + i.CliName + ".json...")
	if err = i.CreateAppJson(appPath, json); err != nil {
		return ErrWhileCleanup
	}

	// Copy core files
	fmt.Println("Run: create core files...")
	if err = i.CopyCoreFiles(os); err != nil {
		return ErrWhileCleanup
	}

	fmt.Println("Run: done!")
	return nil
}
