package main

import (
	"github.com/serkansipahi/app-decorators-cli/util/exec"
	"github.com/serkansipahi/app-decorators-cli/util/watch"
	"log"
	"path"
	"path/filepath"
)

func Server(appPath string, dev bool, production bool, excludeDir string) error {

	babel := filepath.Join(appPath, "node_modules", ".bin", "babel")
	srcPath := filepath.Join(appPath, "src")
	libPath := filepath.Join(appPath, "lib")

	commander := exec.New(false, true, true)
	err := commander.Run([]string{
		babel + " " + srcPath + " --out-dir " + libPath,
	})

	if err != nil {
		return err
	}

	watcher := watch.New(excludeDir)
	watcher.Watch(filepath.Join(appPath, "src"), func(file string) {

		_, fileName := path.Split(file)
		fileExt := filepath.Ext(fileName)
		if fileExt != ".js" {
			return
		}
		srcPath := filepath.Join(appPath, "src", fileName)
		libPath := filepath.Join(appPath, "lib", fileName)
		commander := exec.New(false, true, true)
		err := commander.Run([]string{
			babel + " " + srcPath + " --out-file " + libPath,
		})

		if err != nil {
			log.Fatalln(err)
		}
	})

	return nil
}
