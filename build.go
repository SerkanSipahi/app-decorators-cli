package main

import (
	"os"
	"os/exec"
	"path/filepath"
)

func build(src string, dist string, format string, minify bool, noMangle bool) error {

	var (
		err      error
		commands []string
		jspm     string = filepath.Join("node_modules", ".bin", "jspm")
		jspmCmd  *exec.Cmd
	)

	commands = BuildOptions(BuildOptionsConfig{
		Src:      src,
		Dist:     dist,
		Format:   format,
		NoMangle: noMangle,
		Exclude:  "app-decorators",
		Minify:   minify,
	})

	jspmCmd = exec.Command(jspm, commands...)
	jspmCmd.Stdout = os.Stdout
	jspmCmd.Stderr = os.Stderr

	err = jspmCmd.Run()
	if err != nil {
		return err
	}

	return nil
}
