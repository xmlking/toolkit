package configurator

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"reflect"
	"strings"

	"github.com/creasty/defaults"
	"github.com/markbates/pkger"
	"github.com/markbates/pkger/pkging"
	// "github.com/kelseyhightower/envconfig"
	"github.com/stoewer/go-strcase"
	"gopkg.in/yaml.v2"
)

func getConfigFileWithEnv(file, env string, usePkger bool) (envFile string, err error) {
	var extname = path.Ext(file)

	if extname == "" {
		envFile = fmt.Sprintf("%v.%v", file, env)
	} else {
		envFile = fmt.Sprintf("%v.%v%v", strings.TrimSuffix(file, extname), env, extname)
	}

	if fileInfo, er := stat(envFile, usePkger); er == nil && fileInfo.Mode().IsRegular() {
		return
	}
	return "", fmt.Errorf("failed to find file %v", file)
}

func (c *configurator) getConfigFiles(files ...string) (filesFound []string) {

	if c.config.debug || c.config.verbose {
		fmt.Printf("Current environment: '%v'\n", c.config.environment)
	}

	for _, file := range files {
		found := false

		// check for config file
		if fileInfo, err := stat(file, c.config.usePkger); err == nil && fileInfo.Mode().IsRegular() {
			found = true
			filesFound = append(filesFound, file)
		}

		// check for config file with env
		if file, err := getConfigFileWithEnv(file, c.config.environment, c.config.usePkger); err == nil {
			found = true
			filesFound = append(filesFound, file)
		}

		// still not found? check for example config file
		if !found {
			if example, err := getConfigFileWithEnv(file, "example", c.config.usePkger); err == nil {
				if !c.config.silent {
					fmt.Printf("Failed to find config: %v, using example file: %v\n", file, example)
				}
				filesFound = append(filesFound, example)
			} else if !c.config.silent {
				fmt.Printf("Failed to find config: %v\n", file)
			}
		}
	}
	return
}

func (c *configurator) processFile(config interface{}, file string) (err error) {
	var data []byte
	if c.config.usePkger {
		var fh pkging.File
		if fh, err = pkger.Open(file); err != nil {
			return err
		}
		defer fh.Close()
		if data, err = ioutil.ReadAll(fh); err != nil {
			return err
		}
	} else {
		if data, err = ioutil.ReadFile(file); err != nil {
			return err
		}
	}

	switch {
	case strings.HasSuffix(file, ".yaml") || strings.HasSuffix(file, ".yml"):
		if c.config.errorOnUnmatchedKeys {
			return yaml.UnmarshalStrict(data, config)
		}
		return yaml.Unmarshal(data, config)
	case strings.HasSuffix(file, ".json"):
		return unmarshalJSON(data, config, c.config.errorOnUnmatchedKeys)
	default:

		if err := unmarshalJSON(data, config, c.config.errorOnUnmatchedKeys); err == nil {
			return nil
		} else if strings.Contains(err.Error(), "json: unknown field") {
			return err
		}

		var yamlError error
		if c.config.errorOnUnmatchedKeys {
			yamlError = yaml.UnmarshalStrict(data, config)
		} else {
			yamlError = yaml.Unmarshal(data, config)
		}

		if yamlError == nil {
			return nil
		} else if yErr, ok := yamlError.(*yaml.TypeError); ok {
			return yErr
		}

		return errors.New("failed to decode config")
	}
}

// unmarshalJSON unmarshals the given data into the config interface.
// If the errorOnUnmatchedKeys boolean is true, an error will be returned if there
// are keys in the data that do not match fields in the config interface.
func unmarshalJSON(data []byte, config interface{}, errorOnUnmatchedKeys bool) error {
	reader := strings.NewReader(string(data))
	decoder := json.NewDecoder(reader)

	if errorOnUnmatchedKeys {
		decoder.DisallowUnknownFields()
	}

	err := decoder.Decode(config)
	if err != nil && err != io.EOF {
		return err
	}
	return nil
}

func getPrefixForStruct(prefixes []string, fieldStruct *reflect.StructField) []string {
	if fieldStruct.Anonymous && fieldStruct.Tag.Get("anonymous") == "true" {
		return prefixes
	}
	return append(prefixes, fieldStruct.Name)
}

func (c *configurator) processTags(config interface{}, prefixes ...string) error {
	configValue := reflect.Indirect(reflect.ValueOf(config))
	if configValue.Kind() != reflect.Struct {
		return errors.New("invalid config, should be struct")
	}

	configType := configValue.Type()
	for i := 0; i < configType.NumField(); i++ {
		var (
			envNames    []string
			fieldStruct = configType.Field(i)
			field       = configValue.Field(i)
			envName     = fieldStruct.Tag.Get("env") // read configuration from shell env
		)

		if !field.CanAddr() || !field.CanInterface() {
			continue
		}

		if envName == "" {
			envNames = append(envNames, strcase.UpperSnakeCase(strings.Join(append(prefixes, fieldStruct.Name), "_"))) // CONFIG_DB_NAME
			// envNames = append(envNames, strings.Join(append(prefixes, fieldStruct.Name), "_"))                  // Config_DB_Name
			// envNames = append(envNames, strings.ToUpper(strings.Join(append(prefixes, fieldStruct.Name), "_"))) // CONFIG_DB_NAME
		} else {
			envNames = []string{envName}
		}

		if c.config.verbose {
			fmt.Printf("Trying to load struct `%v`'s field `%v` from env %v\n", configType.Name(), fieldStruct.Name, strings.Join(envNames, ", "))
		}

		// Load From Shell ENV
		for _, env := range envNames {
			if value := os.Getenv(env); value != "" {
				if c.config.debug || c.config.verbose {
					fmt.Printf("Loading configuration for struct `%v`'s field `%v` from env %v...\n", configType.Name(), fieldStruct.Name, env)
				}

				switch reflect.Indirect(field).Kind() {
				case reflect.Bool:
					switch strings.ToLower(value) {
					case "", "0", "f", "false":
						field.Set(reflect.ValueOf(false))
					default:
						field.Set(reflect.ValueOf(true))
					}
				case reflect.String:
					field.Set(reflect.ValueOf(value))
				default:
					if err := yaml.Unmarshal([]byte(value), field.Addr().Interface()); err != nil {
						return err
					}
				}
				break
			}
		}

		if isBlank := reflect.DeepEqual(field.Interface(), reflect.Zero(field.Type()).Interface()); isBlank && fieldStruct.Tag.Get("required") == "true" {
			// return error if it is required but blank
			return errors.New(fieldStruct.Name + " is required, but blank")
		}

		for field.Kind() == reflect.Ptr {
			field = field.Elem()
		}

		if field.Kind() == reflect.Struct {
			if err := c.processTags(field.Addr().Interface(), getPrefixForStruct(prefixes, &fieldStruct)...); err != nil {
				return err
			}
		}

		if field.Kind() == reflect.Slice {
			if arrLen := field.Len(); arrLen > 0 {
				for i := 0; i < arrLen; i++ {
					if reflect.Indirect(field.Index(i)).Kind() == reflect.Struct {
						if err := c.processTags(field.Index(i).Addr().Interface(), append(getPrefixForStruct(prefixes, &fieldStruct), fmt.Sprint(i))...); err != nil {
							return err
						}
					}
				}
			} else {
				// load slice from env
				newVal := reflect.New(field.Type().Elem()).Elem()
				if newVal.Kind() == reflect.Struct {
					idx := 0
					for {
						newVal = reflect.New(field.Type().Elem()).Elem()
						if err := c.processTags(newVal.Addr().Interface(), append(getPrefixForStruct(prefixes, &fieldStruct), fmt.Sprint(idx))...); err != nil {
							return err
						} else if reflect.DeepEqual(newVal.Interface(), reflect.New(field.Type().Elem()).Elem().Interface()) {
							break
						} else {
							idx++
							field.Set(reflect.Append(field, newVal))
						}
					}
				}
			}
		}
	}
	return nil
}

func (c *configurator) load(config interface{}, files ...string) (err error) {
	defer func() {
		if c.config.debug || c.config.verbose {
			if err != nil {
				fmt.Printf("Failed to load configuration from %v, got %v\n", files, err)
			}

			fmt.Printf("Configuration:\n  %#v\n", config)
		}
	}()

	configFiles := c.getConfigFiles(files...)

	for _, file := range configFiles {
		if c.config.debug || c.config.verbose {
			fmt.Printf("Loading configurations from file '%v'...\n", file)
		}
		if err = c.processFile(config, file); err != nil {
			return err
		}
	}

	if c.config.verbose {
		fmt.Printf("Configuration after loading, and before setting Defaults :\n  %#+v\n", config)
	}

	// process defaults
	if err = defaults.Set(config); err != nil {
		return err
	}

	if c.config.verbose {
		fmt.Printf("Configuration after loading files and setting Defaults, before processing ENV:\n  %#v\n", config)
	}

	if c.config.environmentVariablePrefix == "-" { // ???
		err = c.processTags(config)
	} else {
		err = c.processTags(config, c.config.environmentVariablePrefix)
	}

	// validate config only if no parsing errors
	if err == nil {
		err = c.validate.Struct(config)
	}

	return err
}

func stat(name string, usePkger bool) (os.FileInfo, error) {
	if usePkger {
		return pkger.Stat(name)
	} else {
		return os.Stat(name)
	}
}
