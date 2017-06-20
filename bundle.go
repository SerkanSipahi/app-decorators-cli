package main

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

func bundle(appPath string) error {

	var (
		err     error
		libPath string
		cmd     *exec.Cmd
	)

	// remove lib (compiled files) folder
	libPath = filepath.Join(appPath, "lib")
	err = os.RemoveAll(libPath)
	if err != nil {
		log.Fatalln("Something gone wrong while removing: "+libPath, err)
	}

	// compile file
	babel := filepath.Join(appPath, "node_modules", ".bin", "babel")
	srcPath := filepath.Join(appPath, "src")
	libPath = filepath.Join(appPath, "lib")
	babelCmd := exec.Command(
		babel, srcPath, "--out-dir", libPath, "--ignore", "node_modules",
	)
	babelCmd.Stdout = os.Stdout
	babelCmd.Stderr = os.Stderr
	if err = babelCmd.Run(); err != nil {
		return err
	}

	// change directory to module dir e.g. collapsible
	if err = os.Chdir(appPath); err != nil {
		log.Fatalln("Cant change to: "+appPath, err)
	}

	// bundle
	// node_modules/.bin/jspm bundle 'src/**/* - app-decorators - [src/**/*]' deps.js --config jspm.config.json --minify
	cmd = exec.Command("node_modules/.bin/jspm", "bundle",
		"src/index", "-", "app-decorators",
		"lib/index.js",
		"--minify",
		"--no-mangle",
		"--config", "jspm.config.json",
		"--skip-source-maps",
		// FIXME: if --no-mangle the build doesnt work well! it convert e.g component name from collapsible to e.g b
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		log.Fatalln("Cant bundle component: "+appPath, err)
	}

	return nil
}
