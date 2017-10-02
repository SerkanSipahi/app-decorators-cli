package main

import (
	"github.com/serkansipahi/app-decorators-cli/helper"
	"github.com/serkansipahi/app-decorators-cli/install"
	"github.com/serkansipahi/app-decorators-cli/options"
	utilOs "github.com/serkansipahi/app-decorators-cli/util/os"
	"github.com/urfave/cli"
	"log"
	"os"
	"regexp"
	"strings"
)

const (
	CLI_NAME     = "appdec"
	AUTHOR_NAME  = "Serkan Sipahi"
	AUTHOR_EMAIL = "serkan.sipahi@yahoo.de"
	APP_VERSION  = "0.8.249"
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
			Aliases: []string{"c"},
			Usage:   "create new component",
			Flags: []cli.Flag{
				options.Name,
				options.Debug,
			},
			Action: func(c *cli.Context) error {

				// assign passed args
				name := strings.ToLower(c.String("name"))
				name = strings.Trim(name, " ")
				debug := c.Bool("debug")

				regex := regexp.MustCompile(`[^a-zA-Z]+`)
				if ok := regex.MatchString(name); ok {
					log.Fatalln("Only letters between a and z allowed!")
				}

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
		/*
			{
				Name:    "recreate",
				Aliases: []string{"rc"},
				Usage:   "reinit an existing component (without deleting ./src files)",
				Flags: []cli.Flag{
					options.Name,
					options.Debug,
				},
				Action: func(c *cli.Context) error {

					return nil
				},
			},
		*/
		{
			Name:    "run",
			Aliases: []string{"r"},
			Usage:   "starting workflow, see: run --help",
			Flags: []cli.Flag{
				options.Name,
				options.Watch,
				options.Server,
				options.Production,
				options.Minify,
				options.Debug,
				options.Format,
				options.Port,
				options.Browser,
				options.NoMangle,
			},
			Action: func(c *cli.Context) error {

				var name string = strings.ToLower(c.String("name"))
				if name == "" {
					log.Fatalln("\nFailed: please pass component name e.g. --name=component")
				}

				// component has appdec.json
				if err := helper.ModuleExists(name); err != nil {
					log.Fatalln("\nComponent: " + name + " does not exists!")
				}

				var err error = Run(RunConfig{
					Name:       name,
					Watch:      c.Bool("watch"),
					Server:     c.Bool("server"),
					Production: c.Bool("production"),
					Minify:     c.Bool("minify"),
					NoMangle:   c.Bool("no-mangle"),
					Debug:      c.Bool("debug"),
					Ch:         make(chan string, 1),
					KillSigs:   make(chan os.Signal, 1),
					Format:     c.String("format"),
					Port:       c.String("port"),
				})

				if err != nil {
					log.Fatalln(err)
				}

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
		/*
			{
				Name:    "test",
				Aliases: []string{"t"},
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
				Aliases: []string{"p"},
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
				Aliases: []string{"i"},
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
		*/
	}

	app.Run(os.Args)
}
