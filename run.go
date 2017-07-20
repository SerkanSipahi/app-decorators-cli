package main

import (
	"fmt"
	"github.com/serkansipahi/app-decorators-cli/util/file"
	"log"
	"os"
	"os/exec"
	"os/signal"
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
	if err = os.Chdir(c.Name); err != nil {
		log.Fatalln("\nCant change to: "+c.Name, err)
	}

	if !c.Watch && !c.Server && !c.Production && !c.Minify {
		log.Fatalln("Please use any option flag: ./appdec --help")
	}

	// compile files
	if c.Watch {
		go compile("src", "lib", c.Watch, c.Ch)
	} else {
		go func(ch chan<- string) { ch <- "chan: [no watch]" }(c.Ch)
	}

	go func(ch chan string) {
		for {
			var chMsg string = <-ch
			fmt.Println("DEBUG: go1", chMsg)
			if !c.Production {
				ch <- "chan: [no production]"
				return
			}

			if c.CmdBuild != nil {
				c.CmdBuild.Process.Kill()
			}

			c.CmdBuild = build("lib/index.js", "lib/index.js", c.Format, c.Minify, true, true)
			err = c.CmdBuild.Run()
			if err != nil {
				log.Fatalln(err)
			}

			ch <- "chan: [build done]"
		}
	}(c.Ch)

	go func(ch <-chan string) {
		for {
			var chMsg string = <-ch
			fmt.Println("DEBUG: go2", chMsg)
		}
	}(c.Ch)

	if c.Server {
		go webserver("3000")
	}

	<-c.KillSigs

	if c.Production {
		file.DeleteExcept("./lib", "lib/index", "js")
	}

	fmt.Println("Stop appdec!")

	return nil
}
