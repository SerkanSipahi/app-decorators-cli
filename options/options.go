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
	Port = cli.StringFlag{
		Name:  "port",
		Value: "3000",
		Usage: "set port",
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
	Server = cli.BoolFlag{
		Name:  "server",
		Usage: "start server",
	}
	Production = cli.BoolFlag{
		Name:  "production",
		Usage: "set production environment",
	}
	Minify = cli.BoolFlag{
		Name:  "minify",
		Usage: "will minify code",
	}
	CompileFor = cli.StringFlag{
		// babel-preset-env autoprefixer
		Name:  "compile-for",
		Usage: "compile for",
	}
)
