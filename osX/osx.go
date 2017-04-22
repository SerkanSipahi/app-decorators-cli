package osX

import (
	"errors"
	"os/exec"
)

// Errors
var ErrorBinNotExists = errors.New("Please make sure you have git or npm installed")

type Whicher interface {
	Which(string) (string, error)
}

type Which struct{}

// Receivers
func (w Which) Which(bin string) (string, error) {

	path, err := w.lookAtPath(bin)
	if err != nil {
		err = ErrorBinNotExists
	}

	return path, err
}

func (w Which) lookAtPath(bin string) (string, error) {
	return exec.LookPath(bin)
}
