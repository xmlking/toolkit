package ioutil

import (
	"github.com/markbates/pkger"
	"path/filepath"
	"testing"
)

func TestCreateDirectory(t *testing.T) {
	var tempDir = "/config/testing"
	if err := CreateDirectory(tempDir); err != nil {
		t.Error(err)
	}
	if err := pkger.RemoveAll(tempDir); err != nil {
		t.Error(err)
	}
}

// NOTE: Be Careful. you may accidentally delete importent files if you miss-configure `tempFile`
func TestWriteFile(t *testing.T) {
	var tempFile = "/config/testing/sumo.txt"
	if err := WriteFile(tempFile, []byte("ABC"), 077); err != nil {
		t.Error(err)
	}
	if err := pkger.RemoveAll(filepath.Dir(tempFile)); err != nil {
		t.Error(err)
	}
}
