package main

import (
	"log"
	"path/filepath"
	"regexp"
)

type BuildOptionsConfig struct {
	Src            string
	Dist           string
	Exclude        string
	Format         string
	AllowedFormats string
	Alias          string
	NoMangle       bool
	Minify         bool
}

func BuildOptions(opts BuildOptionsConfig) []string {

	var (
		srcPath  string   = filepath.Join(opts.Src)
		distPath string   = filepath.Join(opts.Dist)
		commands []string = []string{}
	)

	if opts.Format != "default" {
		if matched, _ := regexp.MatchString(opts.AllowedFormats, opts.Format); matched == false {
			log.Fatalln("Allowed formats: " + opts.AllowedFormats + " - but you passed: " + opts.Format)
		}
	}

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
		commands = append(commands, "--no-mangle")
	}
	if opts.Format != "default" {
		commands = append(commands, "--format", opts.Format)
	}

	return commands
}
