package main

import (
	"log"
	"os"
	"os/exec"
)

func bundle(appPath string) error {

	var (
		err error
		cmd *exec.Cmd
	)

	// change to component
	if err = os.Chdir(appPath); err != nil {
		log.Fatalln("Cant change to: "+appPath, err)
	}

	// compile files
	err = compile("src", "lib", true, true)
	if err != nil {
		panic(err)
	}

	// bundle
	// node_modules/.bin/jspm bundle 'src/**/* - app-decorators - [src/**/*]' deps.js --config jspm.config.json --minify
	cmd = exec.Command("node_modules/.bin/jspm", "bundle",
		"src/index", "-", "app-decorators",
		"lib/index.js",
		"--minify",
		// FIXME: if --no-mangle the build doesnt work well! it convert e.g component name from collapsible to e.g b
		"--no-mangle",
		"--skip-source-maps",
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		log.Fatalln("Cant bundle component: "+appPath, err)
	}

	return nil
}
