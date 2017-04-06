package main

// ### resources ###
// http://stackoverflow.com/questions/6608873/file-system-scanning-in-golang#6612243
// https://github.com/noypi/filemon/blob/master/example_test.go

import (
	"fmt"
	"github.com/labstack/gommon/log"
	"github.com/noypi/filemon"
	"os"
	"path/filepath"
	"regexp"
)

func visit(path string, f os.FileInfo, err error) error {

	if err != nil {
		return err
	}

	if matched, _ := regexp.MatchString("node_modules", path); matched {
		return nil
	}

	if !f.IsDir() {
		return nil
	}

	fmt.Printf("Visited: %s\n", path)

	return nil
}

func Watch(dir string, callback func(string)) {

	root := filepath.Clean(dir)
	if err := filepath.Walk(root, visit); err {
		log.Fatal(err)
	}

	// create a new watcher
	w := filemon.NewWatcher(func(ev *filemon.WatchEvent) {

		file := fmt.Sprint(ev.Fpath)
		if matched, _ := regexp.MatchString("(__|\\.swp|~)$", file); matched {
			return
		}
		callback(file)
	})

	// watch the current path
	w.Watch(dir)

	// wait for a ctrl+c
	w.WaitForKill()
}
