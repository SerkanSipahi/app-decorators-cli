package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

type customWriter struct {
	callback func()
	w        io.Writer
}

func (cw *customWriter) Write(p []byte) (n int, err error) {
	n, err = cw.w.Write(p)
	if err != nil {
		log.Fatalln(err)
	}
	cw.callback()
	return n, err
}

func compile(src string, dist string, watch bool, callback func()) *exec.Cmd {

	var (
		err      error
		srcPath  string   = filepath.Join(src)
		libPath  string   = filepath.Join(dist)
		commands []string = []string{srcPath, "--out-dir", libPath, "--ignore", "node_modules"}
		babel    string   = filepath.Join("node_modules", ".bin", "babel")
		cmd      *exec.Cmd
		cw       = &customWriter{callback: callback, w: os.Stdout}
	)

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
