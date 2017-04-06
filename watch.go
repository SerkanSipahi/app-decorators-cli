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

func visit(filename string, f os.FileInfo, err error, paths *[]string) error {

	if err != nil {
		return err
	}

	if matched, _ := regexp.MatchString("node_modules", filename); matched {
		return nil
	}

	if !f.IsDir() {
		return nil
	}

	*paths = append(*paths, filename)

	return nil
}

func Watch(dir string, callback func(string)) {

	root := filepath.Clean(dir)
	paths := []string{}
	err := filepath.Walk(root, func(filename string, f os.FileInfo, err error) error {
		visit(filename, f, err, &paths)
		return nil
	})

	fmt.Println(paths)

	if err != nil {
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
