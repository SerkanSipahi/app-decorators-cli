package osx

import (
	"os/exec"
	"os"
	"errors"
	"encoding/json"
	"io/ioutil"
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

func FromStructToJson(config interface{}) ([]byte, error) {
	return json.MarshalIndent(config, "", "\t")
}

func FromJsonFileToStruct(file string, s interface{}) interface{} {

	raw, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}

	json.Unmarshal(raw, &s)
	return s
}