package main

import (
	"path/filepath"
)

type BuildOptionsConfig struct {
	Src      string
	Dist     string
	Exclude  string
	Format   string // write struct that allows only cjs, umd, amd, etc
	NoMangle bool
	Minify   bool
}

func BuildOptions(opts BuildOptionsConfig) []string {

	var (
		srcPath  string   = filepath.Join(opts.Src)
		distPath string   = filepath.Join(opts.Dist)
		commands []string = []string{}
	)

	// set defaults
	if opts.Format == "" {
		opts.Format = "default"
	}

	// build options
	if opts.Format == "default" {
		commands = append(commands, "bundle")
	} else {
		commands = append(commands, "build")
	}

	commands = append(commands, srcPath)
	if opts.Exclude != "" {
		commands = append(commands, "-", opts.Exclude)
	}
	commands = append(commands, distPath)

	if opts.Minify {
		commands = append(commands, "--minify")
	}

	if opts.NoMangle {
		// @iusse
		// when --no-mangle doesnt set it work not as expected!
		// it convert e.g component name from "com-collapsible" to e.g "com-b"
		commands = append(commands, "--no-mangle")
	}
	if opts.Format != "default" {
		commands = append(commands, "--format", opts.Format)
	}

	return commands
}
