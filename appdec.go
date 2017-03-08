package main

import (
	"fmt"
	"log"
	"os"
	"github.com/urfave/cli"
	"github.com/serkansipahi/app-decorators-cli/bootstrap"
	"github.com/serkansipahi/app-decorators-cli/config"
)

func main() {

	_, err := bootstrap.Initialize()
	if err != nil {
		log.Fatalln("Failed while initializing...", err)
	}

	app := cli.NewApp()
	app.Name      = config.AppName
	app.Version   = config.AppVersion
	app.Copyright = "(c) 2017 " + config.AppName
	app.Authors   = []cli.Author {
		cli.Author {
			Name:  config.AuthorName,
			Email: config.AuthorEmail,
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
					Value: "app",
					Usage: "set name of the app",
				},
				cli.BoolFlag {
					Name: "debug",
					Usage: "will show debug messages",
				},
				cli.BoolFlag {
					Name: "force",
					Usage: "force Cmd",
				},
			},
			Action: func(c *cli.Context) error {

				rootPath, err := os.Getwd()
				if err != nil {
					log.Fatalln("Failed whilte get root path")
				}

				name  := c.String("name")
				debug := c.Bool("debug")
				force := c.Bool("force")

				_, err = bootstrap.Install(
					config.Appdec {
						name,
						config.AppVersion,
					},
					rootPath,
					debug,
					force,
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
