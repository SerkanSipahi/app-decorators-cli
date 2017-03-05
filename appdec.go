package main

import (
	"fmt"
	"os"
	"github.com/urfave/cli"
	"os/exec"
	"errors"
	"log"
	"encoding/json"
	"io/ioutil"
)

const (
	cliName     = "appdec"
	authorName  = "Serkan Sipahi"
	authorEmail = "serkan.sipahi@yahoo.de"
	appVersion  = "0.8.204"
	npmPackage  = "app-decorators"
)

type Appdec struct {
	Version string `json:"version"`
	Name    string `json:"name"`
}

func which(bin string) (string, error) {

	path, err := exec.LookPath(bin)
	if err != nil {
		err = errors.New("Please make sure you have '" + bin + "' installed")
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

func jsonStringify(value interface{}) ([]byte, error){

	data, err := json.MarshalIndent(value, "", "\t")

	if err != nil {
		return data, err
	}

	return data, nil
}

func initialize() (int, error) {

	fmt.Println("Run: initialize...")

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
		err                 error
		appPath             string
		currentAbsolutePath string
		cliPackagePath      string
	)

	currentAbsolutePath, err = os.Getwd()
	appPath = currentAbsolutePath + "/" + name
	cliPackagePath = currentAbsolutePath + "/" + cliName + ".json"

	/**
	 * Der soll die appdec.json prüfen
	 */
	if err = os.Chdir(appPath); err == nil && force == false {
		err = errors.New(fmt.Sprintf("\n"+
			"Run: '%s' already created\n" +
			"Please use --force for new init\n" +
			"Note: this command will delete all your files in %s\n"+
		"", name, name))
		return -1, err
	}

	if err = os.RemoveAll(appPath); err != nil {
		return 1, err
	}

	/**
	 * 1.) Wenn appdec.json existiert
	 * 2.) Name auslesen
	 * 3.) Name{ordner}, und appdec.json löschen
	 */

	if err = os.Chdir(appPath); err != nil {
		if err = os.Mkdir(appPath, 0755); err != nil {
			return -1, err
		}
		if err = os.Chdir(appPath); err != nil {
			return -1, err
		}
	}

	fmt.Println("Run: install...")
	cmdNpmInit := command("npm", "init", "-y", debugCommand)
	if err = cmdNpmInit.Run(); err != nil {
		return -1, err
	}

	cmdNpmInstall := command("npm", "install", npmPackage, debugCommand)
	if err = cmdNpmInstall.Run(); err != nil {
		return -1, err
	}

	jsonData, err := jsonStringify(Appdec{appVersion, name})
	if err != nil {
		return -1, err
	}

	fmt.Println("Run: create appdec.json...")
	if err = ioutil.WriteFile(cliPackagePath, jsonData, 0755); err != nil {
		return -1, err
	}

	if err = os.Remove(appPath + "/package.json"); err != nil {
		return -1, err
	}

	return 1, nil
}

func main() {

	_, err := initialize()
	if err != nil {
		log.Fatalln("Failed while initializing...", err)
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
					Name: "production",
					Usage: "production command",
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
