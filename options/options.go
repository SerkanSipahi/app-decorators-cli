package options

import (
	"github.com/urfave/cli"
)

var (
	Name = cli.StringFlag{
		Name:  "name",
		Value: "",
		Usage: "set name of component",
	}
	Timeout = cli.IntFlag{
		Name:  "timeout",
		Value: 60000,
		Usage: "set timeout",
	}
	Port = cli.IntFlag{
		Name:  "port",
		Value: 3000,
		Usage: "set timeout",
	}
	Debug = cli.BoolFlag{
		Name:  "debug",
		Usage: "will show debug messages",
	}
	Browser = cli.StringFlag{
		Name:  "browser",
		Value: "chrome",
		Usage: "will start any defined browser",
	}
	Watch = cli.BoolFlag{
		Name:  "watch",
		Usage: "compile files on any change",
	}
	Format = cli.StringFlag{
		Name:  "format",
		Value: "default",
		Usage: "define component format (amd|cjs|umd|esm)",
	}
	Server = cli.StringFlag{
		Name:  "server",
		Usage: "start server",
		Value: "dev",
	}
	Dev = cli.StringFlag{
		Name:  "dev",
		Usage: "set dev environment",
		Value: "all",
	}
	Production = cli.StringFlag{
		Name:  "production",
		Usage: "set production environment",
		Value: "all",
	}
	SourceMaps = cli.BoolFlag{
		Name:  "source-maps",
		Usage: "will generate source-maps",
	}
	Minify = cli.BoolFlag{
		Name:  "minify",
		Usage: "will minify code",
	}
)
