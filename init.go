package main

import (
	"fmt"
	"github.com/serkansipahi/app-decorators-cli/osx"
)

func Init() (int, error) {

	fmt.Println("Run: initialize...")

	var (
		_   string
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
