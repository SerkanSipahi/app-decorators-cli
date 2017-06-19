package install

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/serkansipahi/app-decorators-cli/helper"
	osx "github.com/serkansipahi/app-decorators-cli/util/os"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
)

var (
	ErrNoModuleName  = errors.New("Failed: Please set module name e.g. 'appdec init --name=mymodule'")
	ErrAppPathExists = errors.New("Failed: Apppath exists")
	ErrSrcPath       = errors.New("Failed: cant create src path")
	ErrWhileCleanup  = errors.New("Failed: someting gone wrong while cleaning")
)

type Config struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type CopyCore interface {
	osx.ReadFiler
	osx.CopyFiler
}

func NpmRun(cmd *exec.Cmd, sigs chan<- os.Signal, appPath string, debug bool) {

	if debug {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}
	cmd.Run()
	status := cmd.ProcessState.Sys().(syscall.WaitStatus)
	exitStatus := status.ExitStatus()
	if exitStatus == -1 {
		if err := os.RemoveAll(appPath); err != nil {
			log.Fatalln(err)
		}
		fmt.Println("\nStopped: installing...")
		close(sigs)
		os.Exit(1)
	}
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

	return i.Install()
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

	name = strings.Title(name)
	srcPath := filepath.Join(appPath, "node_modules", cliDepName, "appdec.json")
	destPath := filepath.Join(appPath, "package.json")

	return i.TplFromTo(srcPath, destPath, map[string]string{
		"{{name}}":    name,
		"{{version}}": i.Version,
	})

}

func (i Install) CreateIndexTpl(appPath string, name string) error {

	fmt.Println("Run: create index.html")
	srcPath := filepath.Join(appPath, "html.tpl")
	destPath := filepath.Join(appPath, "index.html")

	return i.TplFromTo(srcPath, destPath, map[string]string{
		"{{name}}": name,
	})
}

func (i Install) CreateComTpl(appPath string, name string) error {

	fmt.Println("Run: create src/index.js")
	name = strings.Title(name)
	srcPath := filepath.Join(appPath, "component.tpl")
	destPath := filepath.Join(appPath, "src", "index.js")

	return i.TplFromTo(srcPath, destPath, map[string]string{
		"{{name}}": name,
	})

}

func (i Install) TplFromTo(srcTplPath string, destPath string, data map[string]string) error {

	//load template file
	tplData, err := ioutil.ReadFile(srcTplPath)
	if err != nil {
		return err
	}

	destFile, err := os.Create(destPath)
	defer destFile.Close()
	if err != nil {
		return err
	}

	var buf bytes.Buffer
	for k, v := range data {
		if buf.Len() == 0 {
			buf.Write(bytes.Replace(tplData, []byte(k), []byte(v), -1))
		} else {
			res := buf.Bytes()
			buf.Reset()
			buf.Write(bytes.Replace(res, []byte(k), []byte(v), -1))
		}
	}

	if _, err = destFile.Write(buf.Bytes()); err != nil {
		return err
	}
	destFile.Sync()

	return nil
}

func (i Install) Install() error {

	var (
		err        error
		name       string = i.Name
		appPath    string = filepath.Join(i.RootPath, name)
		cliDepName string = "app-decorators-cli-deps"
		cliDeps    string = cliDepName + "@" + i.Version
	)
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt, syscall.SIGTERM)

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

	println("Run: start...")

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
	println("Run: init package.json...")
	NpmRun(
		exec.Command("npm", "init", "-y"),
		sigs,
		appPath,
		i.Debug,
	)

	println("Run: install cli dependencies...")
	NpmRun(
		exec.Command("npm", "install", cliDeps),
		sigs,
		appPath,
		i.Debug,
	)

	// remove package.json
	if err = i.Cleanup(appPath); err != nil {
		return err
	}

	//prepare dependency package.json
	println("Run: prepare dependency package.json...")
	err = i.PrepareDepsPkg(appPath, cliDepName, name)
	if err != nil {
		log.Fatalln(err)
	}

	// Install prepared dependencies
	println("Run: install dependencies...")
	NpmRun(
		exec.Command("npm", "install"),
		sigs,
		appPath,
		i.Debug,
	)

	// Copy core files
	fmt.Println("Run: create core files...")
	if err = i.CopyCoreFiles(osx.Os{}, "package.json"); err != nil {
		return err
	}

	// create index.html
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
	close(sigs)

	return nil
}
