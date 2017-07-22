package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"sync"
	"syscall"
)

type RunConfig struct {
	CmdBuild   *exec.Cmd
	Name       string
	Watch      bool
	Server     bool
	Production bool
	Minify     bool
	Debug      bool
	Ch         chan string
	KillSigs   chan os.Signal
	Format     string
	Port       string
}

func Run(c RunConfig) error {

	//@Todo:
	// kill server (express) before starting

	signal.Notify(c.KillSigs, os.Interrupt, syscall.SIGTERM)

	// change to component directory
	if err := os.Chdir(c.Name); err != nil {
		log.Fatalln("\nCant change to: "+c.Name, err)
	}

	if !c.Watch && !c.Server && !c.Production && !c.Minify {
		log.Fatalln("Please use any option flag: ./appdec --help")
	}

	dist := "lib"
	if c.Production {
		dist = "tmp"
	}

	// compile files
	if c.Watch {
		go compile("src", dist, c.Watch, c.Ch)
	} else {
		go func(ch chan<- string) {
			ch <- "chan: [no watch]"
		}(c.Ch)
	}

	lock := &sync.Mutex{}

	go func(ch <-chan string, lock *sync.Mutex) {
		for {
			select {
			case chMsg := <-ch:
				fmt.Println("DEBUG: go1", chMsg)
				if !c.Production {
					return
				}

				lock.Lock()
				c.CmdBuild = build(filepath.Join(dist, "index.js"), "lib/index.js", c.Format, c.Minify, true, true)
				err = c.CmdBuild.Run()
				if err != nil {
					log.Fatalln(err)
				}
				lock.Unlock()

				if !c.Watch {
					os.Exit(1)
				}
			default:
				// Channel full. Discarding value
			}
		}
	}(c.Ch, lock)

	if c.Server {
		go webserver("3000", lock)
	}

	<-c.KillSigs
	fmt.Println("Stop appdec!")

	return nil
}
