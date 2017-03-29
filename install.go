package main

import (
	"os"
	"fmt"
	"errors"
	"github.com/serkansipahi/app-decorators-cli/osx"
	"io/ioutil"
	"encoding/json"
	"io"
)

func Install(appdecConfig Appdec, rootPath string, cliName string, debugCommand bool) (int, error) {

	var (
		err     error
		name    string = appdecConfig.Name
		appPath string = rootPath + "/" + name
	)

	if name == "" {
		return -1, errors.New("Failed: Please set module name e.g. 'appdec init --name=mymodule'")
	}

	/**
	 * Return when  "appPath" exists
	 */
	if err = os.Chdir(appPath); err == nil {
		err = errors.New(fmt.Sprintf("\n"+
			"Run: '%s' already created\n" +
			"You can delete it with 'appdec delete --name=%s\n'" +
			"", name, name))

		return -1, err
	}

	/**
	 * Create "appPath" if not exists
	 */
	if err = os.Mkdir(appPath, 0755); err != nil {
		fmt.Println("222")
		return -1, err
	}
	if err = os.Chdir(appPath); err != nil {
		fmt.Println("333")
		return -1, err
	}

	/**
	 * Init npm package.json
	 */
	fmt.Println("Run: install...")
	cmdNpmInit := osx.Cmd("npm", "init", "-y", debugCommand)
	if err = cmdNpmInit.Run(); err != nil {
		return -1, err
	}

	/*
	 * Install app-decorators via npm
	 */
	cmdNpmInstall := osx.Cmd("npm", "install", "app-decorators@" + appdecConfig.Version, debugCommand)
	if err = cmdNpmInstall.Run(); err != nil {
		return -1, err
	}

	/**
	 * Create app specific json file
	 */
	jsonData, err := json.MarshalIndent(appdecConfig, "", "\t")
	if err != nil {
		return -1, err
	}

	fmt.Println("Run: create "+ cliName + ".json...")
	if err = ioutil.WriteFile(appPath + "/" + cliName + ".json", jsonData, 0755); err != nil {
		return -1, err
	}

	/**
	 * Cleanup
	 */
	if err = os.Remove(appPath + "/package.json"); err != nil {
		return -1, err
	}

	fmt.Println("Run: done!")

	return 1, nil
}