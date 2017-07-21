package main

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"sync"
	"syscall"
)

var (
	err error
	cmd *exec.Cmd
)

func webserver(port string, lock *sync.Mutex) error {

	if cmd != nil {
		cmd.Process.Kill()
	}

	cmd = exec.Command("node", "server.js", port)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	lock.Lock()
	if err = cmd.Start(); err != nil {
		return err
	}
	lock.Unlock()

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

	lock.Lock()
	if err = cmd.Process.Kill(); err != nil {
		return err
	}
	lock.Unlock()

	fmt.Println("Stop server!")

	return nil

}
