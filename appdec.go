package main

import (
	"fmt"
	"github.com/urfave/cli"
	"log"
	"os"
)

type Appdec struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

const (
	CLI_NAME     = "appdec"
	AUTHOR_NAME  = "Serkan Sipahi"
	AUTHOR_EMAIL = "serkan.sipahi@yahoo.de"
	APP_VERSION  = "0.8.206"
	COPYRIGHT    = "(c) 2017"
)

func main() {

	_, err := Init()
	if err != nil {
		log.Fatalln("Failed while initializing...", err)
	}

	rootPath, err := os.Getwd()
	if err != nil {
		log.Fatalln("Failed whilte get root path")
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

				name := c.String("name")
				debug := c.Bool("debug")
				appdec := Appdec{
					name,
					APP_VERSION,
				}

				_, err = Install(
					appdec,
					rootPath,
					CLI_NAME,
					debug,
				)
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
			},
			Action: func(c *cli.Context) error {
				fmt.Println("completed task: ", c.Args().First())
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
