package main

import (
	"errors"
	"fmt"
	"github.com/serkansipahi/app-decorators-cli/util/watch"
)

func Server(name string, dev bool, production bool) (int8, error) {

	if name == "" {
		return -1, errors.New("Please pass module name")
	}

	watcher := watch.New("node_modules")
	watcher.Watch("./collapsible", func(file string) {
		fmt.Println(file)
	})

	return 1, nil
}
