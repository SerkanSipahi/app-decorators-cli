package main

import (
	"fmt"
	"github.com/serkansipahi/app-decorators-cli/helper"
	"github.com/serkansipahi/app-decorators-cli/install"
	"github.com/serkansipahi/app-decorators-cli/options"
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
	app.Usage = "command line tool"
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
			Name:    "create",
			Aliases: []string{"i"},
			Usage:   "create new component",
			Flags: []cli.Flag{
				options.Name,
				options.Debug,
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
			Name:    "recreate",
			Aliases: []string{"i"},
			Usage:   "reinit an existing component (without deleting ./src files)",
			Flags: []cli.Flag{
				options.Name,
				options.Debug,
			},
			Action: func(c *cli.Context) error {

				return nil
			},
		},
		{
			Name:    "delete",
			Aliases: []string{"d"},
			Usage:   "delete existing component",
			Flags: []cli.Flag{
				options.Name,
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
			Name:    "run",
			Aliases: []string{"s"},
			Usage:   "starting server",
			Flags: []cli.Flag{
				options.Name,
				options.Debug,
				options.Browser,
				options.Watch,
				options.Format,
				options.Server,
				options.Dev,
				options.Production,
				options.Timeout,
				options.SourceMaps,
				options.Minify,
			},
			Action: func(c *cli.Context) error {

				var (
					name       = strings.ToLower(c.String("name"))
					timeout    = c.Int("timeout")
					port       = c.Int("port")
					debug      = c.Bool("debug")
					browser    = c.String("browser")
					watch      = c.Bool("watch")
					format     = c.String("format")
					server     = c.String("server")
					dev        = c.String("dev")
					production = c.String("production")
					SourceMaps = c.Bool("source-maps")
					Minify     = c.Bool("minify")
				)

				fmt.Println("options: ", name, timeout, port, debug, browser, watch, format, server, dev, production, SourceMaps, Minify)

				if name == "" {
					log.Fatalln("Failed: please pass component name e.g. --name=component")
				}

				appPath := filepath.Join(rootPath, name)
				_, module := path.Split(appPath)

				// check if component exists (appdec.json)
				if err := helper.ModuleExists(appPath); err != nil {
					log.Fatalln("Component: " + module + " does not exists!")
				}

				return nil
			},
		},

		{
			Name:    "server",
			Aliases: []string{"s"},
			Usage:   "starting server",
			Flags: []cli.Flag{
				options.Name,
				options.Debug,
				options.Browser,
				options.Watch,
				options.Format,
				options.Dev,
				options.Production,
			},
			Action: func(c *cli.Context) error {

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

				if err := Server(name, true); err != nil {
					log.Fatalln("Failed while Server...", err)
				}

				return nil
			},
		},
		{
			Name:    "build",
			Aliases: []string{"b"},
			Usage:   "build component",
			Flags: []cli.Flag{
				options.Name,
				options.Debug,
				options.Browser,
				options.Server,
				options.Watch,
				options.Format,
			},
			Action: func(c *cli.Context) error {

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
			Name:    "test",
			Aliases: []string{"l"},
			Usage:   "list usage",
			Flags: []cli.Flag{
				options.Name,
			},
			Action: func(c *cli.Context) error {
				fmt.Println("test component")
				return nil
			},
		},
		{
			Name:    "publish",
			Aliases: []string{"l"},
			Usage:   "publish component on npm",
			Flags: []cli.Flag{
				options.Name,
			},
			Action: func(c *cli.Context) error {
				// use lerna (internal)
				fmt.Println("publish component")
				return nil
			},
		},
		{
			Name:    "list",
			Aliases: []string{"l"},
			Usage:   "list all available modules of app-decorators",
			Action: func(c *cli.Context) error {
				fmt.Println("list all modules")
				return nil
			},
		},
		{
			Name:    "install",
			Aliases: []string{"l"},
			Usage:   "install usage",
			Flags: []cli.Flag{
				options.Name,
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
