package main

import (
	"fmt"
	"github.com/serkansipahi/app-decorators-cli/helper"
	"github.com/serkansipahi/app-decorators-cli/install"
	utilOs "github.com/serkansipahi/app-decorators-cli/util/os"
	"github.com/urfave/cli"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
)

const (
	CLI_NAME     = "appdec"
	AUTHOR_NAME  = "Serkan Sipahi"
	AUTHOR_EMAIL = "serkan.sipahi@yahoo.de"
	APP_VERSION  = "0.8.221"
	COPYRIGHT    = "(c) 2017"
)

func main() {

	_, err := Init(utilOs.Os{})
	if err != nil {
		log.Fatalln("Failed while initializing...", err)
	}

	rootPath, err := os.Getwd()
	if err != nil {
		log.Fatalln("Failed while getting root path")
	}
	app := cli.NewApp()
	app.Name = CLI_NAME
	app.Version = APP_VERSION
	app.Copyright = COPYRIGHT + " " + CLI_NAME
	app.Authors = []cli.Author{
		cli.Author{
			Name:  AUTHOR_NAME,
			Email: AUTHOR_EMAIL,
		},
	}

	/**
	 * Setting up allowed commands
	 */
	app.Commands = []cli.Command{
		{
			Name:      "init",
			Aliases:   []string{"i"},
			Usage:     "init usage",
			UsageText: "init usage text",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "name",
					Value: "",
					Usage: "set name of the app",
				},
				cli.BoolFlag{
					Name:  "debug",
					Usage: "will show debug messages",
				},
				cli.IntFlag{
					Name:  "timeout",
					Value: 60000,
					Usage: "set timeout",
				},
			},
			Action: func(c *cli.Context) error {

				// assign passed args
				name := strings.ToLower(c.String("name"))
				debug := c.Bool("debug")

				// create installer
				installer := install.New(
					name,
					rootPath,
					APP_VERSION,
					CLI_NAME,
					debug,
				)

				// start installer
				if err := installer.Run(); err != nil {
					log.Fatalln("Failed while installing...", err)
				}

				return nil
			},
		},
		{
			Name:      "delete",
			Aliases:   []string{"d"},
			Usage:     "delete module",
			UsageText: "delete usage text",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "name",
					Value: "",
					Usage: "set name of the app",
				},
			},
			Action: func(c *cli.Context) error {

				name := strings.ToLower(c.String("name"))
				if name == "" {
					log.Fatalln("Failed: please pass module-name with --name=mymodule")
				}

				err := Delete(rootPath, name)
				if err != nil {
					log.Fatalln("Failed while deleting...", err)
				}

				return nil
			},
		},
		{
			Name:      "server",
			Aliases:   []string{"s"},
			Usage:     "server usage",
			UsageText: "server usage text",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "name",
					Value: "",
					Usage: "set name of the app",
				},
				cli.StringFlag{
					Name:  "browser",
					Value: "",
					Usage: "will start any defined browser",
				},
				cli.BoolFlag{
					Name:  "compile",
					Usage: "will compile files before start server",
				},
			},
			Action: func(c *cli.Context) error {

				name := strings.ToLower(c.String("name"))
				compile := c.Bool("compile")

				if name == "" {
					log.Fatalln("Failed: please pass module-name with --name=mymodule")
				}

				appPath := filepath.Join(rootPath, name)
				_, module := path.Split(appPath)

				// modules exists (appdec.json)
				if err := helper.ModuleExists(appPath); err != nil {
					log.Fatalln("Module: " + module + " does not exists!")
				}

				if err := Server(name, compile); err != nil {
					log.Fatalln("Failed while Server...", err)
				}

				return nil
			},
		},
		{
			Name:      "bundle",
			Aliases:   []string{"b"},
			Usage:     "bundle usage",
			UsageText: "bundle usage text",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "name",
					Value: "",
					Usage: "set name of the app",
				},
				cli.StringFlag{
					Name:  "compile-for",
					Value: "all",
					Usage: "",
				},
			},
			Action: func(c *cli.Context) error {

				// idea: bundle per route

				name := strings.ToLower(c.String("name"))

				if name == "" {
					log.Fatalln("Failed: please pass module-name with --name=mymodule")
				}

				appPath := filepath.Join(rootPath, name)
				_, module := path.Split(appPath)

				// modules exists (appdec.json)
				if err := helper.ModuleExists(appPath); err != nil {
					log.Fatalln("Module: " + module + " does not exists!")
				}

				err := bundle(name)
				if err != nil {
					log.Fatalln("Failed while bundling...", err)
				}
				return nil
			},
		},
		{
			Name:      "list",
			Aliases:   []string{"l"},
			Usage:     "list usage",
			UsageText: "list usage text",
			Action: func(c *cli.Context) error {
				fmt.Println("list all modules")
				return nil
			},
		},
		{
			Name:      "install",
			Aliases:   []string{"l"},
			Usage:     "install usage",
			UsageText: "install usage text",
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "name",
					Usage: "name of module",
				},
			},
			Action: func(c *cli.Context) error {

				// When installing an existing app-dec or vendor module
				// it will store the name of module and the type(existing or vendor)
				// in a file. This is important if we make a bundle/codesplitting
				// for current developed module.

				// Some other ideas:
				//
				fmt.Println("list existing app-dec or vendor module")
				return nil
			},
		},
	}

	app.Run(os.Args)
}
