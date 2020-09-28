package ioutil

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/markbates/pkger"
	"github.com/markbates/pkger/pkging"
)

func ReadFile(filename string) ([]byte, error) {
	f, err := pkger.Open(filename)
	if err != nil {
		return nil, err
	}

	b, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func CreateDirectory(dir string) (err error) {
	if _, err = pkger.Stat(dir); os.IsNotExist(err) {
		err = pkger.MkdirAll(dir, 0755)
	}
	return
}

func WriteFile(fileName string, data []byte, perm os.FileMode) (err error) {
	if err = CreateDirectory(filepath.Dir(fileName)); err != nil {
		return
	}
	var f pkging.File
	if f, err = pkger.Create(fileName); err != nil {
		return
	}
	err = ioutil.WriteFile(filepath.Join(f.Info().Dir+fileName), data, perm)
	if err != nil {
		return
	}

	return
}
