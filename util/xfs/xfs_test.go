package xfs_test

import (
	"embed"
	"io/fs"
	"os"
	"strings"
	"testing"
	"testing/fstest"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/xmlking/toolkit/util/xfs"
)

//go:embed fixtures
var fixtures embed.FS

var testFsys = fstest.MapFS{
	"hello.txt": {
		Data:    []byte("hello, world"),
		Mode:    0456,
		ModTime: time.Now(),
		Sys:     &sysValue,
	},
	"sub/goodbye.txt": {
		Data:    []byte("goodbye, world"),
		Mode:    0456,
		ModTime: time.Now(),
		Sys:     &sysValue,
	},
}

var sysValue int

type readFileOnly struct{ fs.ReadFileFS }

func (readFileOnly) Open(name string) (fs.File, error) { return nil, fs.ErrNotExist }

type openOnly struct{ fs.FS }

func TestReadFile(t *testing.T) {
	// Test that ReadFile uses the method when present.
	data, err := fs.ReadFile(readFileOnly{testFsys}, "hello.txt")
	if string(data) != "hello, world" || err != nil {
		t.Fatalf(`ReadFile(readFileOnly, "hello.txt") = %q, %v, want %q, nil`, data, err, "hello, world")
	}

	// Test that ReadFile uses Open when the method is not present.
	data, err = fs.ReadFile(openOnly{testFsys}, "hello.txt")
	if string(data) != "hello, world" || err != nil {
		t.Fatalf(`ReadFile(openOnly, "hello.txt") = %q, %v, want %q, nil`, data, err, "hello, world")
	}

	// Test that ReadFile on Sub of . works (sub_test checks non-trivial subs).
	sub, err := fs.Sub(testFsys, "sub")
	if err != nil {
		t.Fatal(err)
	}
	data, err = fs.ReadFile(sub, "goodbye.txt")
	if string(data) != "goodbye, world" || err != nil {
		t.Fatalf(`ReadFile(sub(.), "goodbye.txt") = %q, %v, want %q, nil`, data, err, "goodbye, world")
	}
}

func TestXFS_embed(t *testing.T) {
	expected := []string{"fixtures/hello.txt", "fixtures"}
	if err := fstest.TestFS(fixtures, expected...); err != nil {
		t.Fatal(err)
	}
}

func TestXFS_relative_path(t *testing.T) {
	efx := xfs.FS(fixtures)

	b, err := fs.ReadFile(efx, "fixtures/hello.txt")
	assert.NoError(t, err)
	assert.Equal(t, "hello, world", strings.TrimSpace(string(b)))
}

func TestXFS_relative_path_from_project_root(t *testing.T) {
	efx := xfs.FS(fixtures)

	b, err := fs.ReadFile(efx, "util/xfs/fixtures/hello.txt")
	assert.NoError(t, err)
	assert.Equal(t, "hello, world", strings.TrimSpace(string(b)))
}

func TestXFS_absolute_path(t *testing.T) {
	f, err := os.CreateTemp("", "sample")
	assert.NoError(t, err)
	err = os.WriteFile(f.Name(), []byte("hello, world"), 0666)
	t.Log(f.Name())

	efx := xfs.FS(fixtures)
	b, err := fs.ReadFile(efx, f.Name())
	assert.NoError(t, err)
	assert.Equal(t, "hello, world", strings.TrimSpace(string(b)))

	_ = f.Close()
	err = os.Remove(f.Name())
	assert.NoError(t, err)
}
