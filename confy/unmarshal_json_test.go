package confy_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/xmlking/toolkit/confy"
)

func TestUnmatchedKeyInJsonConfigFile(t *testing.T) {
	type configStruct struct {
		Name string
	}
	type configFile struct {
		Name string
		Test string
	}
	config := configFile{Name: "test", Test: "ATest"}

	file, err := os.CreateTemp("/tmp", "confy")
	if err != nil {
		t.Fatal("Could not create temp file")
	}

	if err := json.NewEncoder(file).Encode(config); err == nil {

		var result configStruct

		dir, filename := filepath.Split(file.Name())
		// Do not return error when there are unmatched keys but ErrorOnUnmatchedKeys is false
		confy.DefaultConfy = confy.NewConfy(confy.WithFS(os.DirFS(dir)))
		err = confy.Load(&result, filename)
		assert.NoError(t, err)

		// Return an error when there are unmatched keys and ErrorOnUnmatchedKeys is true
		confy.DefaultConfy = confy.NewConfy(confy.WithErrorOnUnmatchedKeys(), confy.WithFS(os.DirFS(dir)))
		if err := confy.Load(&result, filename); err == nil || !strings.Contains(err.Error(), "json: unknown field") {
			t.Errorf("Should get unknown field error when loading configuration with extra keys. Instead got error: %v", err)
		}

	} else {
		t.Errorf("failed to marshal config")
	}

	// Add .json to the file name and test
	err = os.Rename(file.Name(), file.Name()+".json")
	if err != nil {
		t.Errorf("Could not add suffix to file")
	}
	filename := file.Name() + ".json"

	var result configStruct

	dir, fName := filepath.Split(filename)
	// Do not return error when there are unmatched keys but ErrorOnUnmatchedKeys is false
	confy.DefaultConfy = confy.NewConfy(confy.WithFS(os.DirFS(dir)))
	err = confy.Load(&result, fName)
	assert.NoError(t, err)

	// Return an error when there are unmatched keys and ErrorOnUnmatchedKeys is true
	confy.DefaultConfy = confy.NewConfy(confy.WithErrorOnUnmatchedKeys(), confy.WithFS(os.DirFS(dir)))
	if err := confy.Load(&result, fName); err == nil || !strings.Contains(err.Error(), "json: unknown field") {
		t.Errorf("Should get unknown field error when loading configuration with extra keys. Instead got error: %v", err)
	}

	t.Cleanup(func() {
		t.Log("cleanup...")
		if err := file.Close(); err != nil {
			t.Error(err)
		}
		if err := os.Remove(filename); err != nil {
			t.Error(err)
		}
	})

}
