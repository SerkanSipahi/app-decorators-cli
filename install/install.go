package install

import (
	"errors"
	"fmt"
	"github.com/serkansipahi/app-decorators-cli/helper"
	"github.com/serkansipahi/app-decorators-cli/util/exec"
	"github.com/serkansipahi/app-decorators-cli/util/json"
	osx "github.com/serkansipahi/app-decorators-cli/util/os"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

var (
	ErrNoModuleName        = errors.New("Failed: Please set module name e.g. 'appdec init --name=mymodule'")
	ErrAppPathExists       = errors.New("Failed: Apppath exists")
	ErrSrcPath             = errors.New("Failed: cant create src path")
	ErrCantChangeToApppath = errors.New("Failed: cant change to apppath")
	ErrCantInstallDeps     = errors.New("Failed: someting gone wrong while installing dependencies")
	ErrWhileCleanup        = errors.New("Failed: someting gone wrong while cleaning")
)

type Config struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type CopyCore interface {
	osx.ReadFiler
	osx.CopyFiler
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
		exec.New(false, i.Debug),
	)
}

func (i Install) CreateAppPath(appPath string) error {

	var err error

	if err = os.Mkdir(appPath, 0755); err != nil {
		return ErrAppPathExists
	}
	if err = os.Chdir(appPath); err != nil {
		return ErrCantChangeToApppath
	}

	return nil
}

func (i Install) Cleanup(appPath string) error {

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
	files, err := os.ReadFiles(appDecoratorPath)
	if err != nil {
		return err
	}

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

func (i Install) Install(exec exec.Execer) error {

	var (
		err           error
		name          string = i.Name
		appPath       string = filepath.Join(i.RootPath, name)
		cliDepName    string = "app-decorators-cli-deps"
		cliDeps       string = cliDepName + "@" + i.Version
		cliDepsString string
	)

	if name == "" {
		return ErrNoModuleName
	}

	// Return when "appPath" exists
	if err = helper.ModuleExists(appPath); err == nil {
		err = errors.New(fmt.Sprintf("\n"+
			"Run: '%s' module already created\n"+
			"You can delete it with 'appdec delete --name=%s\n"+
			"", name, name),
		)
		return err
	}

	// Create "appPath" if not exists
	if err := i.CreateAppPath(appPath); err != nil {
		return err
	}

	// create src directory
	srcPath := filepath.Join(appPath, "src")
	if err = os.Mkdir(srcPath, 0755); err != nil {
		return ErrSrcPath
	}

	// Install dependencies
	err = exec.Run([]string{
		"npm init -y",
		"npm install " + cliDeps,
	})
	if err != nil {
		return ErrCantInstallDeps
	}

	// Cleanup
	if err = i.Cleanup(appPath); err != nil {
		return err
	}

	//prepare dependency package.json
	cliDepsPath := filepath.Join(appPath, "node_modules", cliDepName, "appdec.json")
	data, err := ioutil.ReadFile(cliDepsPath)
	if err != nil {
		return err
	}
	cliDepsString = string(data)
	cliDepsString = strings.Replace(cliDepsString, "{{name}}", name, -1)
	cliDepsString = strings.Replace(cliDepsString, "{{version}}", i.Version, -1)

	appPkgJson := filepath.Join(appPath, "package.json")
	appPkgFile, err := os.Create(appPkgJson)
	defer appPkgFile.Close()
	if err != nil {
		return err
	}

	if _, err = appPkgFile.WriteString(cliDepsString); err != nil {
		return err
	}
	appPkgFile.Sync()

	// write into
	// ....

	// close file

	// Create app specific json file
	/*
		fmt.Println("Run: create " + i.CliName + ".json...")
		if err = i.CreateAppJson(appPath, json.New()); err != nil {
			return err
		}
	*/

	// Copy core files
	/*
		fmt.Println("Run: create core files...")
		if err = i.CopyCoreFiles(osx.Os{}); err != nil {
			return err
		}
	*/
	fmt.Println("Run: done!")
	return nil
}
