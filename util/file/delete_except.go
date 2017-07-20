package file

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func DeleteExcept(walk string, exceptFile string, ext string) error {

	if exceptFile == "" {
		log.Fatalln("Please pass tje 'except file'")
	}

	var err error = filepath.Walk(walk, func(fullPath string, f os.FileInfo, err error) error {
		if err != nil {
			log.Fatalln(err)
		}

		if f.IsDir() {
			return nil
		}

		if fullPath == exceptFile+"."+ext {
			return nil
		}

		if filepath.Ext(fullPath) != "."+ext {
			return nil
		}

		err = os.Remove(fullPath)
		if err != nil {
			log.Fatalln(err)

		}
		fmt.Println("DELETE-EXCEPT: ", fullPath)

		return nil
	})

	if err != nil {
		log.Fatalln(err)
	}

	return nil
}
