package bootstrap

import (
	"os"
	"fmt"
	"github.com/serkansipahi/app-decorators-cli/osx"
	"encoding/json"
	"io/ioutil"
	"errors"
	"github.com/serkansipahi/app-decorators-cli/config"
)

func Initialize() (int, error) {

	fmt.Println("Run: initialize...")

	var (
		_ string
		err error
	)

	_, err = osx.Which("git")
	if err != nil {
		return -1, err
	}
	_, err = osx.Which("npm")
	if err != nil {
		return -1, err
	}

	return 1, nil
}

func Install(appdecConfig config.Appdec, rootPath string, debugCommand bool, force bool) (int, error) {

	var (
		err            error
		name           string = appdecConfig.Name
		appPath        string = rootPath + "/" + name
		cliPackagePath string = rootPath + "/" + config.AppName + ".json"
	)

	/**
	 * Der soll die appdec.json prüfen
	 */
	if err = os.Chdir(appPath); err == nil && force == false {
		err = errors.New(fmt.Sprintf("\n"+
			"Run: '%s' already created\n" +
			"Please use --force for new init\n" +
			"Note: this Cmd will delete all your files in %s\n"+
			"", name, name))
		return -1, err
	}

	/**
	 * 1.) Wenn appdec.json existiert
	 * 2.) Name auslesen
	 * 3.) Name{ordner}, und appdec.json löschen
	 */
	if err = os.RemoveAll(appPath); err != nil {
		return 1, err
	}

	if err = os.Chdir(appPath); err != nil {
		if err = os.Mkdir(appPath, 0755); err != nil {
			return -1, err
		}
		if err = os.Chdir(appPath); err != nil {
			return -1, err
		}
	}

	fmt.Println("Run: install...")
	cmdNpmInit := osx.Cmd("npm", "init", "-y", debugCommand)
	if err = cmdNpmInit.Run(); err != nil {
		return -1, err
	}

	cmdNpmInstall := osx.Cmd("npm", "install", config.NpmPackage, debugCommand)
	if err = cmdNpmInstall.Run(); err != nil {
		return -1, err
	}

	// json stringify
	jsonData, err := json.MarshalIndent(appdecConfig, "", "\t")
	if err != nil {
		return -1, err
	}

	fmt.Println("Run: create appdec.json...")
	if err = ioutil.WriteFile(cliPackagePath, jsonData, 0755); err != nil {
		return -1, err
	}

	if err = os.Remove(appPath + "/package.json"); err != nil {
		return -1, err
	}

	return 1, nil
}