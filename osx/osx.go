package osx

import (
	"os/exec"
	"os"
	"errors"
)

func Which(bin string) (string, error) {

	path, err := exec.LookPath(bin)
	if err != nil {
		err = errors.New("Please make sure you have '" + bin + "' installed")
	}

	return path, err
}

func Cmd(name string, arg string, option string, debug ...bool) exec.Cmd {

	cmd := *exec.Command(name, arg, option)
	if debug[0] == true {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}

	return cmd
}