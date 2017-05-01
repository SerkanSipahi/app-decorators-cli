package os

import (
	"errors"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
)

// interfaces
type Whicher interface {
	Which(string) (string, error)
}

type ReadFiler interface {
	ReadFiles(src string) ([]os.FileInfo, error)
}

type CopyFiler interface {
	CopyFile(src, dst string) (err error)
}

// struct
type Os struct{}

func (o Os) Which(bin string) (string, error) {

	path, err := o.lookAtPath(bin)
	if err != nil {
		err = errors.New("Please make sure you have " + bin + " installed")
	}

	return path, err
}

func (o Os) lookAtPath(bin string) (string, error) {
	return exec.LookPath(bin)
}

// CopyFile/CopyDir https://gist.github.com/m4ng0squ4sh/92462b38df26839a3ca324697c8cba04
// CopyFile copies the contents of the file named src to the file named
// by dst. The file will be created if it does not already exist. If the
// destination file exists, all it's contents will be replaced by the contents
// of the source file. The file mode will be copied from the source and
// the copied data is synced/flushed to stable storage.
func (o Os) CopyFile(src, dst string) (err error) {

	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	distFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer func() {
		if e := distFile.Close(); e != nil {
			err = e
		}
	}()

	_, err = io.Copy(distFile, srcFile)
	if err != nil {
		return err
	}

	err = distFile.Sync()
	if err != nil {
		return err
	}

	// Get Fileinfo
	fileInfo, err := os.Stat(src)
	if err != nil {
		return err
	}
	// Sync Filemode from src and src
	err = os.Chmod(dst, fileInfo.Mode())
	if err != nil {
		return err
	}

	return nil
}

// Read all files from given src directory. If given directory not
// exists or is invalid it return an error
func (o Os) ReadFiles(src string) ([]os.FileInfo, error) {

	list, err := ioutil.ReadDir(filepath.Clean(src))
	if err != nil {
		return list, err
	}

	var files []os.FileInfo
	for _, entry := range list {
		// Skip directories
		if entry.IsDir() {
			continue
		}
		// Skip symlinks
		if entry.Mode()&os.ModeSymlink != 0 {
			continue
		}

		files = append(files, entry)
	}

	return files, nil
}
