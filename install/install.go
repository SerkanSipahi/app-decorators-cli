package install

import (
	"errors"
	"fmt"
	"github.com/serkansipahi/app-decorators-cli/util/exec"
	"github.com/serkansipahi/app-decorators-cli/util/json"
	"github.com/serkansipahi/app-decorators-cli/util/os"
	"io/ioutil"
	"log"
	"path/filepath"
)

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

type Config struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type Install struct {
	Name     string
	AppPath  string
	RootPath string
	CliName  string
	Version  string
	Debug    bool
}

func (i Install) Run() (int, error) {

	return i.Install(
		os.Os{},
		json.New(),
		exec.New(false, i.Debug),
	)
}

func (i Install) Install(os os.Os, json json.Stringifyer, exec exec.Execers) (int, error) {

	var (
		err     error
		name    string = i.Name
		appPath string = i.RootPath + "/" + name
		// npm packages
		appDecPkg   string = "app-decorators@" + i.Version
		babelCliPkg string = "babel-cli@6.24.1"
	)

	if name == "" {
		return -1, errors.New("Failed: Please set module name e.g. 'appdec init --name=mymodule'")
	}

	if name == "commands" || name == "osx" {
		return -1, errors.New("Failed: cant install module as '" + name + "' because its reserved")
	}

	/**
	 * Return when  "appPath" exists
	 */
	if err = os.Chdir(appPath); err == nil {
		err = errors.New(fmt.Sprintf("\n"+
			"Run: '%s' already created\n"+
			"You can delete it with 'appdec delete --name=%s\n"+
			"", name, name))

		return -1, err
	}

	/**
	 * Create "appPath" if not exists
	 */
	if err = os.Mkdir(appPath, 0755); err != nil {
		return -1, err
	}
	if err = os.Chdir(appPath); err != nil {
		return -1, err
	}

	err = exec.Run([]string{
		"npm init -y",
		"npm install " + appDecPkg,
		"npm install " + babelCliPkg,
	})

	if err != nil {
		log.Fatal(err)
	}

	/**
	 * Cleanup
	 */
	pkgJsonPath := filepath.Join(appPath, "package.json")
	if err = os.Remove(pkgJsonPath); err != nil {
		return -1, err
	}

	/**
	 * Create app specific json file
	 */
	config := Config{
		Name:    i.Name,
		Version: i.Version,
	}
	jsonData, err := json.Stringify(config)
	if err != nil {
		return -1, err
	}

	fmt.Println("Run: create " + i.CliName + ".json...")
	appDecJsonPath := filepath.Join(appPath, i.CliName+".json")
	if err = ioutil.WriteFile(appDecJsonPath, jsonData, 0755); err != nil {
		return -1, err
	}

	/**
	 * Copy core files
	 */
	fmt.Println("Run: create core files...")

	appDecoratorPath := filepath.Clean(
		filepath.Join(i.RootPath, i.Name, "node_modules", "app-decorators"),
	)
	files, _ := os.ReadFiles(appDecoratorPath)
	for _, file := range files {

		src := filepath.Join(appDecoratorPath, file.Name())
		dist := filepath.Join(i.RootPath, i.Name, file.Name())

		err := os.CopyFile(src, dist)
		if err != nil {
			return -1, err
		}
	}

	fmt.Println("Run: done!")

	return 1, nil
}
