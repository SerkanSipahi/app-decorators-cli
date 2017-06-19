package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"syscall"
)

func Server(appPath string, compile bool) error {

	var err error
	babel := filepath.Join(appPath, "node_modules", ".bin", "babel")
	srcPath := filepath.Join(appPath, "src")
	webRoot := filepath.Join(appPath, "lib")

	// compile file
	babelCmd := exec.Command(
		babel, srcPath, "--out-dir", webRoot, "--source-maps", "--watch", "--ignore", "node_modules",
	)
	if compile {
		babelCmd.Stdout = os.Stdout
		babelCmd.Stderr = os.Stderr
		if err = babelCmd.Start(); err != nil {
			return err
		}
	}

	// start server
	if err = os.Chdir(appPath); err != nil {
		return errors.New("Cant change to: " + appPath)
	}
	serverCmd := exec.Command("node", "server.js")
	serverCmd.Stdout = os.Stdout
	serverCmd.Stderr = os.Stderr
	if err = serverCmd.Start(); err != nil {
		return err
	}

	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	signal.Notify(sigs, os.Interrupt, syscall.SIGTERM)

	go func() {
		sigs := <-sigs
		if sigs.String() != "interrupt" {
			return
		}
		done <- true
	}()
	<-done

	if compile {
		fmt.Println("\nStop compiler")
		if err = babelCmd.Process.Kill(); err != nil {
			return err
		}
	}
	fmt.Println("Stop server!")
	if err = serverCmd.Process.Kill(); err != nil {
		return err
	}

	return nil
}
