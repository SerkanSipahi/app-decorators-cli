package helper

import (
	"os"
	"path/filepath"
)

func ModuleExists(appPath string) error {

	moduleJsonFilePath := filepath.Clean(
		filepath.Join(appPath, "appdec.json"),
	)
	moduleFile, err := os.Open(moduleJsonFilePath)
	defer moduleFile.Close()

	return err
}
