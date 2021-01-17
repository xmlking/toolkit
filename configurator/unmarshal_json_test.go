package configurator_test

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/xmlking/toolkit/configurator"
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

	file, err := ioutil.TempFile("/tmp", "configurator")
	if err != nil {
		t.Fatal("Could not create temp file")
	}

	filename := file.Name()

	if err := json.NewEncoder(file).Encode(config); err == nil {

		var result configStruct

		// Do not return error when there are unmatched keys but ErrorOnUnmatchedKeys is false
		if err := configurator.NewConfigurator().Load(&result, filename); err != nil {
			t.Errorf("Should NOT get error when loading configuration with extra keys. Error: %v", err)
		}

		// Return an error when there are unmatched keys and ErrorOnUnmatchedKeys is true
		if err := configurator.NewConfigurator(configurator.WithErrorOnUnmatchedKeys()).Load(&result, filename); err == nil || !strings.Contains(err.Error(), "json: unknown field") {

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
	filename = file.Name() + ".json"

	var result configStruct

	// Do not return error when there are unmatched keys but ErrorOnUnmatchedKeys is false
	if err := configurator.NewConfigurator().Load(&result, filename); err != nil {
		t.Errorf("Should NOT get error when loading configuration with extra keys. Error: %v", err)
	}

	// Return an error when there are unmatched keys and ErrorOnUnmatchedKeys is true
	if err := configurator.NewConfigurator(configurator.WithErrorOnUnmatchedKeys()).Load(&result, filename); err == nil || !strings.Contains(err.Error(), "json: unknown field") {

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
