package main

import (
	"log"
	"os"
)

func bundle(appPath string) error {

	var err error

	// change to component
	if err = os.Chdir(appPath); err != nil {
		log.Fatalln("Cant change to: "+appPath, err)
	}

	// compile files
	err = compile("src", "lib", false, true)
	if err != nil {
		panic(err)
	}

	// bundle
	err = build("src/index.js", "lib/index.js", "default", true, true)
	if err != nil {
		panic(err)
	}

	return nil
}
