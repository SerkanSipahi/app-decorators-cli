package main

import (
	"fmt"
	"os"
	"github.com/urfave/cli"
	"os/exec"
	"errors"
	"log"
)

const (
	cliName     string = "appdec"
	authorName  string = "Serkan Sipahi"
	authorEmail string = "serkan.sipahi@yahoo.de"
	appVersion  string = "0.8.204"
	npmPackage  string = "app-decorators"
)

func which(binary string) (string, error) {

	path, err := exec.LookPath(binary)
	if err != nil {
		err = errors.New("Please make sure you have '" + binary + "' installed")
	}

	return path, err
}

func command(name string, arg string, option string, debug ...bool) exec.Cmd {

	cmd := *exec.Command(name, arg, option)
	if debug[0] == true {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}
	return cmd
}

func cd(directory string) error {
	return os.Chdir(directory)
}

func mkdir(directory string) error {
	return os.Mkdir(directory, 0700)
}

func pwd() (string, error) {
	return os.Getwd()
}

func initialize() (int, error) {

	var (
		_ string
		err error
	)

	_, err = which("git")
	if err != nil {
		return -1, err
	}
	_, err = which("npm")
	if err != nil {
		return -1, err
	}

	return 1, nil
}

func install(name string, debugCommand bool, force bool) (int, error) {

	var (
		err error
		currentPath string
		appPath string
	)

	currentPath, err = pwd()
	appPath = currentPath+  "/" + name

	if err = cd(appPath); err == nil && force == false {
		fmt.Println("Process: " + name + " already created! Please use --force for new init")
		return 1, nil
	}

	fmt.Println("Process: create required directory...")
	if err = cd(appPath); err != nil {
		if err = mkdir(appPath); err != nil {
			return -1, err
		}
	}

	fmt.Println("Process: change to dir...")
	if err = cd(appPath); err != nil {
		return -1, err
	}

	fmt.Println("Process: init package.json...")
	cmdNpmInit := command("npm", "init", "-y", debugCommand)
	if err := cmdNpmInit.Run(); err != nil {
		return -1, err
	}

	fmt.Println("Process: install required libs...")
	cmdNpmInstall := command("npm", "install", npmPackage, debugCommand)
	if err := cmdNpmInstall.Run(); err != nil {
		return -1, err
	}

	/**
	 * Am Ende eine appdec.json erstellen die alle infos enthälts,
	 * so das beim nächsten mal wenn du server gestartet etc, die app
	 * weis wo sie zugreifen soll
	 */

	return 1, nil
}

func main() {

	_, err := initialize()
	if err != nil {
		log.Fatal("Failed while initializing...", err)
	}

	app := cli.NewApp()
	app.Name = cliName
	app.Version = appVersion
	app.Copyright = "(c) 2017 " + cliName
	app.Authors = []cli.Author {
		cli.Author {
			Name:  authorName,
			Email: authorEmail,
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
					Usage: "force command",
				},
			},
			Action : func(c *cli.Context) error {

				_, err = install(
					c.String("name"),
					c.Bool("debug"),
					c.Bool("force"),
				)
				if err != nil {
					log.Fatal("Failed while installing...", err)
				}

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
					Usage: "force command",
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
					Name: "live",
					Usage: "force command",
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
