package main

import (
	"fmt"
	"github.com/serkansipahi/app-decorators-cli/install"
	"github.com/serkansipahi/app-decorators-cli/osX"
	"github.com/urfave/cli"
	"log"
	"os"
)

// @todo/@fixme
// - check weather an new version app-decorator is available. When "yes" ask for upgrade

// build: go build *.go

func main() {

	//config := install.Config{"a", "1.0"}
	//installer := install.New(config, "b", "c")
	//installer.Run()

	_, err := Init(osX.Which{})
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
			},
			Action: func(c *cli.Context) error {

				// assign passed args
				name := c.String("name")
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
				_, err := installer.Run()
				if err != nil {
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

				name := c.String("name")
				_, err := Delete(name, rootPath, CLI_NAME)
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
				cli.BoolFlag{
					Name:  "dev",
					Usage: "will show debug messages",
				},
				cli.BoolFlag{
					Name:  "live",
					Usage: "force Cmd",
				},
				cli.StringFlag{
					Name:  "name",
					Value: "",
					Usage: "set name of the app",
				},
			},
			Action: func(c *cli.Context) error {

				name := c.String("name")
				if name == "" {
					log.Fatalln("Failed: please pass module-name with --name=mymodule")
				}

				if _, err := Server(name); err != nil {
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
				cli.BoolFlag{
					Name:  "dev",
					Usage: "will show debug messages",
				},
				cli.BoolFlag{
					Name:  "production",
					Usage: "production Cmd",
				},
			},
			Action: func(c *cli.Context) error {
				fmt.Println("completed task: ", c.Args().First())
				return nil
			},
		},
	}

	app.Run(os.Args)
}
