package main

// resources/inspiration:
// http://stackoverflow.com/questions/6608873/file-system-scanning-in-golang#6612243
// https://github.com/noypi/filemon/blob/master/example_test.go

// @todo:
// - move watch.go to own repo: http://github.com/serkansipahi/watcher
// - check weather on runtime new directory with file will created
// - it should not call two times successively for same file change
// - allow to pass something like this ./collapsible -r --ignore=node_modules

import (
	"errors"
	"fmt"
	"github.com/noypi/filemon"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"regexp"
)

func getDir(filename string, f os.FileInfo, params ...string) (error, string) {

	var excludeDir string
	if len(params) == 1 {
		excludeDir = params[0]
	}

	if matched, _ := regexp.MatchString(excludeDir, filename); matched {
		return errors.New("matched excludeDir"), filename
	}

	if !f.IsDir() {
		return errors.New("no dir"), filename
	}

	return nil, filename
}

func watchPath(path string, callback func(string)) {

	// create a new watcher
	w := filemon.NewWatcher(func(ev *filemon.WatchEvent) {
		file := fmt.Sprint(ev.Fpath)
		// ignore __ (intellij), .swp and ~ files
		if matched, _ := regexp.MatchString("(__|\\.swp|~)$", file); matched {
			return
		}

		callback(file)
	})

	w.Watch(path)
}

func Watch(dir string, callback func(string)) {

	root := filepath.Clean(dir)
	paths := []string{}

	err := filepath.Walk(root, func(name string, f os.FileInfo, err error) error {

		if err != nil {
			panic(err)
		}

		if err, name = getDir(name, f, "node_modules"); err == nil {
			paths = append(paths, name)
		}

		return nil
	})

	if err != nil {
		log.Fatal(err)
	}

	// watch determined paths
	for _, path := range paths {
		go watchPath(path, callback)
	}

	// wait for kill signal
	onKill := make(chan os.Signal, 1)
	signal.Notify(onKill, os.Interrupt, os.Kill)
	<-onKill // wait for event

	fmt.Fprintln(os.Stderr, "\nKill triggered. Exiting...")
}
