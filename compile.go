package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func compile(src string, dist string, watch bool, debug bool) error {

	var (
		err      error
		srcPath  string   = filepath.Join(src)
		libPath  string   = filepath.Join(dist)
		commands []string = []string{srcPath, "--out-dir", libPath, "--ignore", "node_modules"}
		babel    string   = filepath.Join("node_modules", ".bin", "babel")
		babelCmd *exec.Cmd
	)

	// remove compiled files
	err = os.RemoveAll(libPath)
	if err != nil {
		return errors.New(fmt.Sprint("Something gone wrong while removing compiled files: "+libPath, err))
	}

	if watch {
		commands = append(commands, "--watch")
	}

	babelCmd = exec.Command(babel, commands...)

	if debug {
		babelCmd.Stdout = os.Stdout
		babelCmd.Stderr = os.Stderr
	}

	if err = babelCmd.Run(); err != nil {
		return err
	}

	return nil
}
