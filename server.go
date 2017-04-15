package main

import (
	"errors"
	"fmt"
)

func Server(name string) (int8, error) {

	if name == "" {
		return -1, errors.New("Please pass module name")
	}

	Watch("./collapsible", func(file string) {
		fmt.Println(file)
	})

	return 1, nil
}
