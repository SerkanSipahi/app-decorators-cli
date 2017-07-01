package file

import (
	"log"
	"os"
	"path/filepath"
	"regexp"
)

type CountOptions struct {
	Ignore   string
	FileType string
}

// Count count recursive the count of files in passed directory
// as second parameter it can passed a ignore dir
func Count(dir string, option CountOptions) int {

	var countFile int = 0

	var err error = filepath.Walk(dir, func(fullPath string, f os.FileInfo, err error) error {
		if err != nil {
			log.Fatalln(err)
		}

		if f.IsDir() {
			return nil
		}
		if option.Ignore != "" {
			if matched, _ := regexp.MatchString(option.Ignore, fullPath); matched == true {
				return nil
			}
		}
		if option.FileType != "" {
			if filepath.Ext(fullPath) != "."+option.FileType {
				return nil
			}
		}

		countFile++

		return nil

	})

	if err != nil {
		log.Fatalln(err)
	}

	return countFile
}
