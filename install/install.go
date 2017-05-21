package install

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/serkansipahi/app-decorators-cli/helper"
	"github.com/serkansipahi/app-decorators-cli/util/exec"
	osx "github.com/serkansipahi/app-decorators-cli/util/os"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

var (
	ErrNoModuleName    = errors.New("Failed: Please set module name e.g. 'appdec init --name=mymodule'")
	ErrAppPathExists   = errors.New("Failed: Apppath exists")
	ErrSrcPath         = errors.New("Failed: cant create src path")
	ErrCantInstallDeps = errors.New("Failed: someting gone wrong while installing dependencies")
	ErrWhileCleanup    = errors.New("Failed: someting gone wrong while cleaning")
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

func (i Install) CopyCoreFiles(os CopyCore, ignore string) error {

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

		if file.Name() == ignore {
			continue
		}

		err := os.CopyFile(src, dist)
		if err != nil {
			return err
		}
	}

	return nil
}

func (i Install) PrepareDepsPkg(appPath string, cliDepName string, name string) error {

	cliDepsPath := filepath.Join(appPath, "node_modules", cliDepName, "appdec.json")
	tplData, err := ioutil.ReadFile(cliDepsPath)
	if err != nil {
		return err
	}

	tplData = bytes.Replace(tplData, []byte("{{name}}"), []byte(name), -1)
	tplData = bytes.Replace(tplData, []byte("{{version}}"), []byte(i.Version), -1)

	appPkgJson := filepath.Join(appPath, "package.json")
	appPkgFile, err := os.Create(appPkgJson)
	if err != nil {
		return err
	}

	if _, err = appPkgFile.Write(tplData); err != nil {
		return err
	}
	appPkgFile.Sync()
	appPkgFile.Close()

	return nil
}

// @todo: refactor, to it with go templates
func (i Install) CreateIndexTpl(appPath string, name string) error {

	//load template file
	fmt.Println("Run: create index.html...")
	tplPath := filepath.Join(appPath, "html.tpl")
	tplData, err := ioutil.ReadFile(tplPath)
	if err != nil {
		return err
	}

	indexHTMLPath := filepath.Join(appPath, "index.html")
	indexHTMLFile, err := os.Create(indexHTMLPath)
	if err != nil {
		return err
	}

	tplByte := bytes.Replace(tplData, []byte("{{name}}"), []byte(name), -1)
	if _, err = indexHTMLFile.Write(tplByte); err != nil {
		return err
	}
	indexHTMLFile.Sync()
	indexHTMLFile.Close()

	return nil
}

// @todo: refactor, to it with go templates
func (i Install) CreateComTpl(appPath string, name string) error {

	name = strings.Title(name)

	//load template file
	fmt.Println("Run: create src/index.js")
	tplPath := filepath.Join(appPath, "component.tpl")
	tplData, err := ioutil.ReadFile(tplPath)
	if err != nil {
		return err
	}

	indexHTMLPath := filepath.Join(appPath, "src", "index.js")
	indexHTMLFile, err := os.Create(indexHTMLPath)
	if err != nil {
		return err
	}

	tplByte := bytes.Replace(tplData, []byte("{{name}}"), []byte(name), -1)
	if _, err = indexHTMLFile.Write(tplByte); err != nil {
		return err
	}
	indexHTMLFile.Sync()
	indexHTMLFile.Close()

	return nil
}

func (i Install) Install(exec exec.Execer) error {

	/**
	 * @Todo: Implement lighthouse, too ! for measuring (use app-decorators-cli-deps)
	 */
	var (
		err        error
		name       string = i.Name
		appPath    string = filepath.Join(i.RootPath, name)
		cliDepName string = "app-decorators-cli-deps"
		cliDeps    string = cliDepName + "@" + i.Version
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

	fmt.Println("Run: install...")

	// Create "appPath" if not exists
	if err := i.CreateAppPath(appPath); err != nil {
		return err
	}

	// create src directory
	srcPath := filepath.Join(appPath, "src")
	if err = os.Mkdir(srcPath, 0755); err != nil {
		return ErrSrcPath
	}

	//////////////////////////////////////////////////
	// change directory to module dir e.g. collapsible
	//////////////////////////////////////////////////
	if err = os.Chdir(appPath); err != nil {
		return errors.New("Cant change to: " + appPath)
	}

	// Get package configuration template
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
	err = i.PrepareDepsPkg(appPath, cliDepName, name)
	if err != nil {
		log.Fatalln(err)
	}

	// Install prepared dependencies
	err = exec.Run([]string{
		"npm install",
	})
	if err != nil {
		return err
	}

	// Copy core files
	fmt.Println("Run: create core files...")
	if err = i.CopyCoreFiles(osx.Os{}, "package.json"); err != nil {
		return err
	}

	err = i.CreateIndexTpl(appPath, name)
	if err != nil {
		log.Fatalln(err)
	}

	// create src/index.js
	err = i.CreateComTpl(appPath, name)
	if err != nil {
		log.Fatalln(err)
	}

	// Done
	fmt.Println("Run: done!")
	return nil
}
