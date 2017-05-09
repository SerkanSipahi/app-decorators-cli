package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

func Delete(rootPath string, module string, cliName string) error {

	modulePath := filepath.Join(rootPath, module)
	moduleFile, err := os.Open(modulePath)
	if err != nil {
		return errors.New("Failed: module '" + module + "' does not exists")
	}
	defer moduleFile.Close()

	// module json exists
	moduleJsonFilePath := filepath.Clean(
		filepath.Join(rootPath, module, cliName+".json"),
	)
	moduleJsonFile, err := os.Open(moduleJsonFilePath)
	if err != nil {
		return errors.New("Failed: module '" + module + "' is not part of " + cliName + ".json")
	}
	defer moduleJsonFile.Close()

	// remove module
	moduleDirectory := filepath.Clean(
		filepath.Join(rootPath, module),
	)
	if err := os.RemoveAll(moduleDirectory); err != nil {
		return err
	}

	fmt.Println("Run: removed " + module)

	return nil
}
