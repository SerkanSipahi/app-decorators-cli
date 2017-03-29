package main

import (
	"fmt"
	"log"
	"os"
	"github.com/urfave/cli"
)

type Appdec struct {
	Name     string `json:"name"`
	Version  string `json:"version"`
}

const (
	CLI_NAME     = "appdec"
	AUTHOR_NAME  = "Serkan Sipahi"
	AUTHOR_EMAIL = "serkan.sipahi@yahoo.de"
	APP_VERSION  = "0.8.204"
	COPYRIGHT    = "(c) 2017"
)

func main() {

	_, err := Init()
	if err != nil {
		log.Fatalln("Failed while initializing...", err)
	}

	app := cli.NewApp()
	app.Name      = CLI_NAME
	app.Version   = APP_VERSION
	app.Copyright = COPYRIGHT + "" + CLI_NAME
	app.Authors   = []cli.Author {
		cli.Author {
			Name:  AUTHOR_NAME,
			Email: AUTHOR_EMAIL,
		},
	}

	/**
	 * Setting up allowed commands
	 */
	app.Commands = []cli.Command {
		{
			Name     : "init",
			Aliases  : []string{"i"},
			Usage    : "init usage",
			UsageText: "init usage text",
			Flags    : []cli.Flag {
				cli.StringFlag {
					Name: "name",
					Value: "",
					Usage: "set name of the app",
				},
				cli.BoolFlag {
					Name: "debug",
					Usage: "will show debug messages",
				},
			},
			Action: func(c *cli.Context) error {

				rootPath, err := os.Getwd()
				if err != nil {
					log.Fatalln("Failed whilte get root path")
				}

				name   := c.String("name")
				debug  := c.Bool("debug")
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

				fmt.Println("Run: done!")

				return nil
			},

		},
		{
			Name     : "server",
			Aliases  : []string{"s"},
			Usage    : "server usage",
			UsageText: "server usage text",
			Flags    : []cli.Flag {
				cli.BoolFlag {
					Name: "dev",
					Usage: "will show debug messages",
				},
				cli.BoolFlag {
					Name: "live",
					Usage: "force Cmd",
				},
			},
			Action : func(c *cli.Context) error {
				fmt.Println("completed task: ", c.Args().First())
				return nil
			},
		},
		{
			Name     : "bundle",
			Aliases  : []string{"b"},
			Usage    : "bundle usage",
			UsageText: "bundle usage text",
			Flags    : []cli.Flag {
				cli.BoolFlag {
					Name: "dev",
					Usage: "will show debug messages",
				},
				cli.BoolFlag {
					Name: "production",
					Usage: "production Cmd",
				},
			},
			Action : func(c *cli.Context) error {
				fmt.Println("completed task: ", c.Args().First())
				return nil
			},
		},
		{
			Name     : "delete",
			Aliases  : []string{"d"},
			Usage    : "delete module",
			UsageText: "delete usage text",
			Flags    : []cli.Flag {
				cli.BoolFlag {
					Name: "force",
					Usage: "force Cmd",
				},
			},
			Action : func(c *cli.Context) error {
				fmt.Println("completed task: ", c.Args().First())
				return nil
			},
		},
	}

	app.Run(os.Args)
}
