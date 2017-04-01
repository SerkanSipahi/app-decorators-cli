package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

func Delete(module string, rootPath string, cliName string) (int, error) {

	// module json exists
	moduleJsonFilePath := filepath.Clean(
		filepath.Join(rootPath, module, cliName+".json"),
	)
	moduleFile, err := os.Open(moduleJsonFilePath)
	if err != nil {
		return -1, errors.New("Failed: module '" + module + "' is not part of " + cliName + ".json")
	}
	defer moduleFile.Close()

	// remove module
	moduleDirectory := filepath.Clean(
		filepath.Join(rootPath, module),
	)
	if err := os.RemoveAll(moduleDirectory); err != nil {
		return -1, err
	}

	fmt.Println("Run: removed " + module)

	return 1, nil
}
