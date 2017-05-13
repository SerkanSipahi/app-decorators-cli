package helper

import (
	"errors"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"path/filepath"
)

func ModuleExists(appPath string) error {

	// Module package.json property exists
	pkgJsonPath := filepath.Join(appPath, "package.json")
	pkgDataJson, _ := ioutil.ReadFile(pkgJsonPath)

	pkgDataStruct := string(pkgDataJson)
	name := gjson.Get(pkgDataStruct, "appdec.name")
	version := gjson.Get(pkgDataStruct, "appdec.version")

	if name.Exists() && version.Exists() {
		return nil
	}

	return errors.New("Module does not exists!")
}
