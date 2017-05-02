package watch

// resources/inspiration:
// http://stackoverflow.com/questions/6608873/file-system-scanning-in-golang#6612243
// https://github.com/noypi/filemon/blob/master/example_test.go

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

func New(excludeDir ...string) *Watch {
	return &Watch{
		ExcludeDir: excludeDir[0],
	}
}

type Watch struct {
	ExcludeDir string
}

func (w Watch) visit(filename string, f os.FileInfo, params ...string) (error, string) {

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

func (w Watch) watchDir(path string, callback func(string)) {

	// create a new watcher
	fmw := filemon.NewWatcher(func(ev *filemon.WatchEvent) {
		file := fmt.Sprint(ev.Fpath)
		// ignore __ (intellij), .swp and ~ files
		if _, err := os.Stat(file); err != nil {
			return
		}

		// it should not call two times successively for same file change
		// create 0, modify 1, delete 2, rename 3, attrib 5 see ev.Type
		evT := ev.Type
		if evT == 1 {
			return
		}

		callback(file)
	})

	fmw.Watch(path)
}

func (w Watch) Watch(dir string, callback func(string)) {

	root := filepath.Clean(dir)
	paths := []string{}

	err := filepath.Walk(root, func(name string, f os.FileInfo, err error) error {

		if err != nil {
			panic(err)
		}

		if err, name = w.visit(name, f, w.ExcludeDir); err == nil {
			paths = append(paths, name)
		}

		return nil
	})

	if err != nil {
		log.Fatal(err)
	}

	// watch determined paths
	for _, path := range paths {
		go w.watchDir(path, callback)
	}

	// wait for kill signal
	onKill := make(chan os.Signal, 1)
	signal.Notify(onKill, os.Interrupt, os.Kill)
	<-onKill // wait for event

	fmt.Fprintln(os.Stderr, "\nKill triggered. Exiting...")
}
