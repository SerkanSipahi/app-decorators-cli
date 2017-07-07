package main

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
)

var (
	err error
	cmd *exec.Cmd
)

func webserver(port string) error {

	if cmd != nil {
		cmd.Process.Kill()
	}

	cmd = exec.Command("node", "server.js", port)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err = cmd.Start(); err != nil {
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

	if err = cmd.Process.Kill(); err != nil {
		return err
	}

	fmt.Println("Stop server!")
	os.Exit(1)

	return nil

}
