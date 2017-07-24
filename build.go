package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type BuildWrite struct {
	w io.Writer
}

func (cw BuildWrite) Write(p []byte) (n int, err error) {

	n, err = cw.w.Write(p)
	if err != nil {
		log.Fatalln(err)
	}
	return n, err
}

func build(src, dist, format string, minify, noMangle, debug bool) *exec.Cmd {

	var (
		jspm     string = filepath.Join("node_modules", ".bin", "jspm")
		commands []string
		cmd      *exec.Cmd
	)

	commands = BuildOptions(BuildOptionsConfig{
		Src:            src,
		Dist:           dist,
		Format:         format,
		AllowedFormats: "default|static|cjs|amd|umd",
		NoMangle:       noMangle,
		Exclude:        "app-decorators",
		Minify:         minify,
	})

	var bw BuildWrite = BuildWrite{w: os.Stdout}

	fmt.Println("COMMAND: ", jspm, strings.Join(commands, " "))
	if debug {
		//fmt.Println("COMMAND: ", jspm, strings.Join(commands, " "))
	}

	cmd = exec.Command(jspm, commands...)
	//if debug {
	cmd.Stdout = bw
	cmd.Stderr = bw
	//}

	return cmd
}
