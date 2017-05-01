package main

import (
	"errors"
	"fmt"
	"github.com/serkansipahi/app-decorators-cli/util/watcher"
)

func Server(name string) (int8, error) {

	if name == "" {
		return -1, errors.New("Please pass module name")
	}

	watcher.Watch("./collapsible", func(file string) {
		fmt.Println(file)
	})

	return 1, nil
}
