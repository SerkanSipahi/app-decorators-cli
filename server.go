package main

import (
	"errors"
	"fmt"
	"github.com/serkansipahi/app-decorators-cli/helper"
	"github.com/serkansipahi/app-decorators-cli/util/watch"
	"path"
)

func Server(appPath string, dev bool, production bool, excludeDir string) error {

	_, module := path.Split(appPath)

	if err := helper.ModuleExists(appPath); err != nil {
		return errors.New("Module: " + module + " does not exists!")
	}

	watcher := watch.New(excludeDir)
	watcher.Watch(module, func(file string) {
		fmt.Println(file)
	})

	return nil
}
