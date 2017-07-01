package main

import (
	"os"
	"os/exec"
	"path/filepath"
)

func build(src string, dist string, format string, minify bool, noMangle bool, debug bool) *exec.Cmd {

	var (
		jspm     string = filepath.Join("node_modules", ".bin", "jspm")
		commands []string
		cmd      *exec.Cmd
	)

	commands = BuildOptions(BuildOptionsConfig{
		Src:      src,
		Dist:     dist,
		Format:   format,
		NoMangle: noMangle,
		Exclude:  "app-decorators",
		Minify:   minify,
	})

	cmd = exec.Command(jspm, commands...)
	if debug {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}

	return cmd
}
