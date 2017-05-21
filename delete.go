package main

import (
	"fmt"
	"github.com/serkansipahi/app-decorators-cli/helper"
	"log"
	"os"
	"path"
	"path/filepath"
)

func Delete(rootPath string, name string) error {

	appPath := filepath.Join(rootPath, name)
	_, module := path.Split(appPath)

	if err := helper.ModuleExists(appPath); err != nil {
		log.Fatalln("Module: " + module + " does not exists!")
	}

	if err := os.RemoveAll(appPath); err != nil {
		return err
	}

	fmt.Println("Run: removed " + module)

	return nil
}
