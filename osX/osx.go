package osX

import (
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
)

// Errors
var ErrorBinNotExists = errors.New("Please make sure you have git or npm installed")

type Which struct{}

type Whicher interface {
	Which(string) (string, error)
}

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

// Commander
type Commander interface {
	Run(name string, arg string, option string, debug ...bool)
}

type Command struct {
	Name   string
	Arg    string
	Option string
	Debug  bool
}

func (c Command) Run() error {
	return c.run()
}

func (c Command) run() error {
	cmd := *exec.Command(c.Name, c.Arg, c.Option)
	if c.Debug == true {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}

	return cmd.Run()
}

// CopyFile/CopyDir https://gist.github.com/m4ng0squ4sh/92462b38df26839a3ca324697c8cba04

// CopyFile copies the contents of the file named src to the file named
// by dst. The file will be created if it does not already exist. If the
// destination file exists, all it's contents will be replaced by the contents
// of the source file. The file mode will be copied from the source and
// the copied data is synced/flushed to stable storage.
func CopyFile(src, dst string) (err error) {

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
func ReadFiles(src string) ([]os.FileInfo, error) {

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

func JsonStringify(config interface{}) ([]byte, error) {
	return json.MarshalIndent(config, "", "\t")
}

func JsonParse(file string, s interface{}) interface{} {

	raw, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}

	json.Unmarshal(raw, &s)
	return s
}
