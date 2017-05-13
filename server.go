package main

import (
	"github.com/serkansipahi/app-decorators-cli/util/exec"
	"path/filepath"
)

func Server(appPath string, dev bool, production bool, excludeDir string) error {

	babel := filepath.Join(appPath, "node_modules", ".bin", "babel")
	srcPath := filepath.Join(appPath, "src")
	webRoot := filepath.Join(appPath, "webroot", "lib")

	commander := exec.New(false, true, true)
	err := commander.Run([]string{
		babel + " " + srcPath + " --out-dir " + webRoot + " --watch",
	})

	if err != nil {
		return err
	}

	return nil
}
