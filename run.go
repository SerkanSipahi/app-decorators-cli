package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strings"
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
	NoMangle   bool
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

	dist := "lib"
	if c.Production {
		dist = "tmp"
	}

	// compile files
	go compile("src", dist, c.Watch, c.Ch)

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
				tmpPath := filepath.Join(dist, "index.js")
				libPath := filepath.Join("lib", "index.js")
				c.CmdBuild = build(tmpPath, libPath, c.Format, c.Minify, c.NoMangle, c.Debug)
				err = c.CmdBuild.Run()
				if err != nil {
					log.Fatalln(err)
				}

				read, err := ioutil.ReadFile(libPath)
				newContents := strings.Replace(string(read), tmpPath, libPath, -1)
				err = ioutil.WriteFile(libPath, []byte(newContents), 0)
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
