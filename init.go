package main

import "github.com/serkansipahi/app-decorators-cli/util/os"

func Init(e os.Whicher) (int, error) {

	var (
		_   string
		err error
	)

	_, err = e.Which("git")
	if err != nil {
		return -1, err
	}
	_, err = e.Which("npm")
	if err != nil {
		return -1, err
	}

	return 1, nil
}
