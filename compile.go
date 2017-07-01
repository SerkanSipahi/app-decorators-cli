package main

import (
	"fmt"
	"github.com/serkansipahi/app-decorators-cli/util/file"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

var (
	writeCount          int
	callCallbackOnCount int
	skip                bool
)

type CustomWrite struct {
	callback func()
	w        io.Writer
}

func (cw CustomWrite) Write(p []byte) (n int, err error) {

	writeCount++

	n, err = cw.w.Write(p)
	if err != nil {
		log.Fatalln(err)
	}

	if writeCount == callCallbackOnCount || skip {
		cw.callback()
		skip = true
	}
	return n, err
}

func compile(src string, dist string, watch bool, callback func()) *exec.Cmd {

	var (
		err      error
		srcPath  string   = filepath.Join(src)
		libPath  string   = filepath.Join(dist)
		babel    string   = filepath.Join("node_modules", ".bin", "babel")
		commands []string = []string{srcPath, "--out-dir", libPath, "--ignore", "node_modules"}
		count    int      = 0
		cmd      *exec.Cmd
	)

	count = file.Count(src, file.CountOptions{
		Ignore:   "node_modules",
		FileType: "js",
	})

	var cw CustomWrite = CustomWrite{
		callback: callback,
		w:        os.Stdout,
	}
	callCallbackOnCount = count

	// remove compiled files
	err = os.RemoveAll(libPath)
	if err != nil {
		panic(fmt.Sprint("Something gone wrong while removing compiled files: "+libPath, err))
	}

	if watch {
		commands = append(commands, "--watch")
	}

	cmd = exec.Command(babel, commands...)
	cmd.Stdout = cw
	cmd.Stderr = cw

	return cmd
}
