package confy_test

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"testing/fstest"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"

	"github.com/xmlking/toolkit/confy"
)

//go:embed fixtures/goodbye.yml
var goodbye []byte

var sysValue int
var testFS = fstest.MapFS{
	"hello.json": {
		Data:    []byte(`{"APPName":"confy","Hosts":["http://example.org","http://jinzhu.me"],"DB":{"Name":"confy","User":"confy","Password":"confy","Port":3306,"SSL":true},"Contacts":[{"Name":"sumo","Email":"sumo@gmail.com"},{"Name":"sumo2","Email":"sumo2@gmail.com"}],"Description":"This is an anonymous embedded struct whose environment variables should NOT include 'ANONYMOUS'"}`),
		Mode:    0456,
		ModTime: time.Now(),
		Sys:     &sysValue,
	},
	"sub/goodbye.yaml": {
		Data:    goodbye,
		Mode:    0456,
		ModTime: time.Now(),
		Sys:     &sysValue,
	},
}

type Anonymous struct {
	Description string
}

type Database struct {
	Name     string
	User     string `yaml:",omitempty" default:"root"`
	Password string `required:"true" env:"DBPassword"`
	Port     uint   `default:"3306" yaml:",omitempty" json:",omitempty"`
	SSL      bool   `default:"true" yaml:",omitempty" json:",omitempty"`
}

type Contact struct {
	Name  string `default:"sumo" yaml:",omitempty"  json:",omitempty"`
	Email string `required:"true"`
}

type testConfig struct {
	APPName   string   `default:"confy" yaml:",omitempty" json:",omitempty"`
	Hosts     []string `validate:"omitempty,dive,url" default:"[\"https://abc.org\"]"`
	DB        *Database
	Contacts  []Contact
	Anonymous `anonymous:"true" default:"-"`
	private   string
}

func generateDefaultConfig(t *testing.T) testConfig {
	t.Helper()
	return testConfig{
		APPName: "confy",
		Hosts:   []string{"http://example.org", "http://jinzhu.me"},
		DB: &Database{
			Name:     "confy",
			User:     "confy",
			Password: "confy",
			Port:     3306,
			SSL:      true,
		},
		Contacts: []Contact{
			{
				Name:  "sumo",
				Email: "sumo@gmail.com",
			},
			{
				Name:  "sumo2",
				Email: "sumo2@gmail.com",
			},
		},
		Anonymous: Anonymous{
			Description: "This is an anonymous embedded struct whose environment variables should NOT include 'ANONYMOUS'",
		},
	}
}

func setup() {
	confy.DefaultConfy = confy.NewConfy()
	fmt.Println("Setup completed")
}

func teardown() {
	// Do something here.
	confy.DefaultConfy = nil
	fmt.Println("Teardown completed")
}

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	teardown()
	os.Exit(code)
}

func TestLoadNormaltestConfig(t *testing.T) {

	config := generateDefaultConfig(t)
	if bytes, err := json.Marshal(config); err == nil {
		if file, err := os.CreateTemp("/tmp", "confy"); err == nil {
			t.Cleanup(func() {
				t.Log("cleanup...")
				if err := file.Close(); err != nil {
					t.Error(err)
				}
				if err := os.Remove(file.Name()); err != nil {
					t.Error(err)
				}
			})

			file.Write(bytes)

			var result testConfig
			dir, filename := filepath.Split(file.Name())
			confy.DefaultConfy = confy.NewConfy(confy.WithFS(os.DirFS(dir)))
			err = confy.Load(&result, filename)
			assert.NoError(t, err)
			assert.Exactly(t, result, config, "result should equal to original configuration")
		}

	} else {
		t.Errorf("failed to marshal config")
	}

}

// CONFY_DEBUG_MODE=true CONFY_VERBOSE_MODE=true go test -v -run TestDefaultValue -count=1 ./confy/...
func TestDefaultValue(t *testing.T) {
	config := generateDefaultConfig(t)
	config.APPName = ""
	config.DB.Port = 0
	config.DB.SSL = false
	config.Contacts[0].Name = ""

	if bytes, err := json.Marshal(config); err == nil {
		if file, err := os.CreateTemp("/tmp", "confy"); err == nil {
			defer file.Close()
			defer os.Remove(file.Name())
			file.Write(bytes)

			var result testConfig
			dir, filename := filepath.Split(file.Name())
			confy.DefaultConfy = confy.NewConfy(confy.WithFS(os.DirFS(dir)))
			err = confy.Load(&result, filename)
			assert.NoError(t, err)
			assert.Exactly(t, result, generateDefaultConfig(t), "result should equal to original configuration")
		}
	} else {
		t.Errorf("failed to marshal config")
	}
}

func TestMissingRequiredValue(t *testing.T) {
	config := generateDefaultConfig(t)
	config.DB.Password = ""

	if bytes, err := json.Marshal(config); err == nil {
		if file, err := os.CreateTemp("/tmp", "confy"); err == nil {
			defer file.Close()
			defer os.Remove(file.Name())
			file.Write(bytes)

			var result testConfig
			if err := confy.Load(&result, file.Name()); err == nil {
				t.Errorf("Should got error when load configuration missing db password")
			}
		}
	} else {
		t.Errorf("failed to marshal config")
	}
}

func TestUnmatchedKeyInYamltestConfigFile(t *testing.T) {
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

	defer os.Remove(file.Name())
	defer file.Close()

	filename := file.Name()

	if data, err := yaml.Marshal(config); err == nil {
		file.WriteString(string(data))

		var result configStruct

		dir, fName := filepath.Split(file.Name())
		// Do not return error when there are unmatched keys but ErrorOnUnmatchedKeys is false
		confy.DefaultConfy = confy.NewConfy(confy.WithFS(os.DirFS(dir)))

		if err := confy.Load(&result, fName); err != nil {
			t.Errorf("Should NOT get error when loading configuration with extra keys. Error: %v", err)
		}

		// Return an error when there are unmatched keys and ErrorOnUnmatchedKeys is true
		confy.DefaultConfy = confy.NewConfy(confy.WithErrorOnUnmatchedKeys(), confy.WithFS(os.DirFS(dir)))
		if err := confy.Load(&result, fName); err == nil {
			t.Errorf("Should get error when loading configuration with extra keys")

			// The error should be of type *yaml.TypeError
		} else if _, ok := err.(*yaml.TypeError); !ok {
			// || !strings.Contains(err.Error(), "not found in struct") {
			t.Errorf("Error should be of type yaml.TypeError. Instead error is %v", err)
		}

	} else {
		t.Errorf("failed to marshal config")
	}

	// Add .yaml to the file name and test again
	err = os.Rename(filename, filename+".yaml")
	if err != nil {
		t.Errorf("Could not add suffix to file")
	}
	filename = filename + ".yaml"
	defer os.Remove(filename)

	var result configStruct

	dir, fName := filepath.Split(filename)
	confy.DefaultConfy = confy.NewConfy(confy.WithFS(os.DirFS(dir)))

	// Do not return error when there are unmatched keys but ErrorOnUnmatchedKeys is false
	if err := confy.NewConfy().Load(&result, fName); err != nil {
		t.Errorf("Should NOT get error when loading configuration with extra keys. Error: %v", err)
	}

	// Return an error when there are unmatched keys and ErrorOnUnmatchedKeys is true
	confy.DefaultConfy = confy.NewConfy(confy.WithErrorOnUnmatchedKeys(), confy.WithFS(os.DirFS(dir)))
	if err := confy.Load(&result, fName); err == nil {
		t.Errorf("Should get error when loading configuration with extra keys")

		// The error should be of type *yaml.TypeError
	} else if _, ok := err.(*yaml.TypeError); !ok {
		// || !strings.Contains(err.Error(), "not found in struct") {
		t.Errorf("Error should be of type yaml.TypeError. Instead error is %v", err)
	}
}

func TestYamlDefaultValue(t *testing.T) {
	config := generateDefaultConfig(t)
	config.APPName = ""
	config.DB.Port = 0
	config.DB.SSL = false

	if bytes, err := yaml.Marshal(config); err == nil {
		if file, err := os.CreateTemp("/tmp", "confy.*.yaml"); err == nil {
			defer file.Close()
			defer os.Remove(file.Name())
			file.Write(bytes)

			var result testConfig
			dir, filename := filepath.Split(file.Name())
			confy.DefaultConfy = confy.NewConfy(confy.WithFS(os.DirFS(dir)))
			err = confy.Load(&result, filename)
			assert.NoError(t, err)

			assert.Exactly(t, result, generateDefaultConfig(t), "result should equal to original configuration")
		}
	} else {
		t.Errorf("failed to marshal config")
	}
}

func TestLoadConfigurationByEnvironment(t *testing.T) {
	config := generateDefaultConfig(t)
	config2 := struct {
		APPName string
	}{
		APPName: "config2",
	}

	if file, err := os.CreateTemp("/tmp", "confy"); err == nil {
		defer file.Close()
		defer os.Remove(file.Name())
		configBytes, _ := yaml.Marshal(config)
		config2Bytes, _ := yaml.Marshal(config2)
		os.WriteFile(file.Name()+".yaml", configBytes, 0644)
		defer os.Remove(file.Name() + ".yaml")
		os.WriteFile(file.Name()+".production.yaml", config2Bytes, 0644)
		defer os.Remove(file.Name() + ".production.yaml")

		var result testConfig
		os.Setenv("CONFY_ENV", "production")
		defer os.Setenv("CONFY_ENV", "")

		dir, filename := filepath.Split(file.Name())
		confy.DefaultConfy = confy.NewConfy(confy.WithFS(os.DirFS(dir)))
		err = confy.Load(&result, filename)

		if err := confy.Load(&result, filename+".yaml"); err != nil {
			t.Errorf("No error should happen when load configurations, but got %v", err)
		}

		defaultConfig := generateDefaultConfig(t)
		defaultConfig.APPName = "config2"
		assert.Exactly(t, result, defaultConfig, "result should equal to original configuration")
	}

	t.Cleanup(func() {
		t.Log("Cleanup...")
	})
}

func TestLoadtestConfigurationByEnvironmentSetBytestConfig(t *testing.T) {
	config := generateDefaultConfig(t)
	config2 := struct {
		APPName string
	}{
		APPName: "production_config2",
	}

	if file, err := os.CreateTemp("/tmp", "confy"); err == nil {
		defer file.Close()
		defer os.Remove(file.Name())
		configBytes, _ := yaml.Marshal(config)
		config2Bytes, _ := yaml.Marshal(config2)
		os.WriteFile(file.Name()+".yaml", configBytes, 0644)
		defer os.Remove(file.Name() + ".yaml")
		os.WriteFile(file.Name()+".production.yaml", config2Bytes, 0644)
		defer os.Remove(file.Name() + ".production.yaml")

		var result testConfig

		dir, filename := filepath.Split(file.Name())
		confy.DefaultConfy = confy.NewConfy(confy.WithEnvironment("production"), confy.WithFS(os.DirFS(dir)))
		err = confy.Load(&result, filename+".yaml")
		assert.NoError(t, err)

		defaultConfig := generateDefaultConfig(t)
		defaultConfig.APPName = "production_config2"
		assert.Exactly(t, result, defaultConfig, "result should be load configurations by environment correctly")

		if confy.GetEnvironment() != "production" {
			t.Errorf("confy's environment should be production")
		}
	}
}

func TestOverwritetestConfigurationWithEnvironmentWithDefaultPrefix(t *testing.T) {
	config := generateDefaultConfig(t)

	if bytes, err := json.Marshal(config); err == nil {
		if file, err := os.CreateTemp("/tmp", "confy"); err == nil {
			file.Write(bytes)
			var result testConfig
			os.Setenv("CONFY_APP_NAME", "config2")
			os.Setenv("CONFY_HOSTS", "- http://example2.org\n- http://jinzhu2.me")
			os.Setenv("CONFY_APP_NAME", "config2")
			os.Setenv("CONFY_DB_NAME", "db_name")

			t.Cleanup(func() {
				t.Log("cleanup...")
				if err := file.Close(); err != nil {
					t.Error(err)
				}
				if err := os.Remove(file.Name()); err != nil {
					t.Error(err)
				}
				os.Setenv("CONFY_APP_NAME", "")
				os.Setenv("CONFY_HOSTS", "")
				os.Setenv("CONFY_APP_NAME", "")
				os.Setenv("CONFY_DB_NAME", "")
			})

			dir, filename := filepath.Split(file.Name())
			confy.DefaultConfy = confy.NewConfy(confy.WithFS(os.DirFS(dir)))
			err = confy.Load(&result, filename)
			assert.NoError(t, err)

			t.Log(result.Contacts[0])
			t.Log(result.Contacts[1])

			defaultConfig := generateDefaultConfig(t)
			defaultConfig.APPName = "config2"
			defaultConfig.Hosts = []string{"http://example2.org", "http://jinzhu2.me"}
			defaultConfig.DB.Name = "db_name"
			assert.Exactly(t, result, defaultConfig, "result should equal to original configuration")
		}
	} else {
		t.Errorf("failed to marshal config")
	}
}

// go test -v -run TestENV -count=1 ./confy/...
func TestENV(t *testing.T) {
	if confy.GetEnvironment() != "test" {
		t.Skipf("skipping test. Env should be test when running `go test`, instead env is %v", confy.GetEnvironment())
	}

	os.Setenv("CONFY_ENV", "production")
	defer os.Setenv("CONFY_ENV", "")
	confy.DefaultConfy = confy.NewConfy()
	if confy.GetEnvironment() != "production" {
		t.Errorf("Env should be production when set it with CONFY_ENV")
	}
}

func TestOverwritetestConfigurationWithEnvironment(t *testing.T) {
	config := generateDefaultConfig(t)

	if bytes, err := json.Marshal(config); err == nil {
		if file, err := os.CreateTemp("/tmp", "confy"); err == nil {
			defer file.Close()
			defer os.Remove(file.Name())
			file.Write(bytes)
			var result testConfig
			os.Setenv("CONFY_ENV_PREFIX", "app")
			os.Setenv("APP_APP_NAME", "config2")
			os.Setenv("APP_DB_NAME", "db_name")
			defer os.Setenv("CONFY_ENV_PREFIX", "")
			defer os.Setenv("APP_APP_NAME", "")
			defer os.Setenv("APP_DB_NAME", "")

			dir, filename := filepath.Split(file.Name())
			confy.DefaultConfy = confy.NewConfy(confy.WithFS(os.DirFS(dir)))
			err = confy.Load(&result, filename)
			assert.NoError(t, err)

			defaultConfig := generateDefaultConfig(t)
			defaultConfig.APPName = "config2"
			defaultConfig.DB.Name = "db_name"
			assert.Exactly(t, result, defaultConfig, "result should equal to original configuration")
		}
	} else {
		t.Errorf("failed to marshal config")
	}
}

func TestOverwritetestConfigurationWithEnvironmentThatSetBytestConfig(t *testing.T) {
	config := generateDefaultConfig(t)

	if bytes, err := json.Marshal(config); err == nil {
		if file, err := os.CreateTemp("/tmp", "confy"); err == nil {
			defer file.Close()
			defer os.Remove(file.Name())
			file.Write(bytes)
			os.Setenv("APP1_APP_NAME", "config2")
			os.Setenv("APP1_DB_NAME", "db_name")
			defer os.Setenv("APP1_APP_NAME", "")
			defer os.Setenv("APP1_DB_NAME", "")

			var result testConfig
			dir, filename := filepath.Split(file.Name())
			confy.DefaultConfy = confy.NewConfy(confy.WithEnvironmentVariablePrefix("APP1"), confy.WithFS(os.DirFS(dir)))
			err = confy.Load(&result, filename)
			assert.NoError(t, err)

			defaultConfig := generateDefaultConfig(t)
			defaultConfig.APPName = "config2"
			defaultConfig.DB.Name = "db_name"
			assert.Exactly(t, result, defaultConfig, "result should equal to original configuration")
		}
	} else {
		t.Errorf("failed to marshal config")
	}
}

func TestResetPrefixToBlank(t *testing.T) {
	config := generateDefaultConfig(t)

	if bytes, err := json.Marshal(config); err == nil {
		if file, err := os.CreateTemp("/tmp", "confy"); err == nil {
			defer file.Close()
			defer os.Remove(file.Name())
			file.Write(bytes)
			var result testConfig
			os.Setenv("CONFY_ENV_PREFIX", "-")
			os.Setenv("APP_NAME", "config2")
			os.Setenv("DB_NAME", "db_name")
			defer os.Setenv("CONFY_ENV_PREFIX", "")
			defer os.Setenv("APP_NAME", "")
			defer os.Setenv("DB_NAME", "")

			dir, filename := filepath.Split(file.Name())
			confy.DefaultConfy = confy.NewConfy(confy.WithFS(os.DirFS(dir)))
			err = confy.Load(&result, filename)
			assert.NoError(t, err)

			defaultConfig := generateDefaultConfig(t)
			defaultConfig.APPName = "config2"
			defaultConfig.DB.Name = "db_name"
			assert.Exactly(t, result, defaultConfig, "result should equal to original configuration")
		}
	} else {
		t.Errorf("failed to marshal config")
	}
}

func TestResetPrefixToBlank2(t *testing.T) {
	config := generateDefaultConfig(t)

	if bytes, err := json.Marshal(config); err == nil {
		if file, err := os.CreateTemp("/tmp", "confy"); err == nil {
			defer file.Close()
			defer os.Remove(file.Name())
			file.Write(bytes)
			var result testConfig
			os.Setenv("CONFY_ENV_PREFIX", "-")
			os.Setenv("APP_NAME", "config2")
			os.Setenv("DB_NAME", "db_name")
			defer os.Setenv("CONFY_ENV_PREFIX", "")
			defer os.Setenv("APPName", "")
			defer os.Setenv("DB_NAME", "")

			dir, filename := filepath.Split(file.Name())
			confy.DefaultConfy = confy.NewConfy(confy.WithFS(os.DirFS(dir)))
			err = confy.Load(&result, filename)
			assert.NoError(t, err)

			defaultConfig := generateDefaultConfig(t)
			defaultConfig.APPName = "config2"
			defaultConfig.DB.Name = "db_name"
			assert.Exactly(t, result, defaultConfig, "result should equal to original configuration")
		}
	} else {
		t.Errorf("failed to marshal config")
	}
}

func TestReadFromEnvironmentWithSpecifiedEnvName(t *testing.T) {
	config := generateDefaultConfig(t)

	if bytes, err := json.Marshal(config); err == nil {
		if file, err := os.CreateTemp("/tmp", "confy"); err == nil {
			defer file.Close()
			defer os.Remove(file.Name())
			file.Write(bytes)
			var result testConfig
			os.Setenv("DBPassword", "db_password")
			defer os.Setenv("DBPassword", "")

			dir, filename := filepath.Split(file.Name())
			confy.DefaultConfy = confy.NewConfy(confy.WithFS(os.DirFS(dir)))
			err = confy.Load(&result, filename)
			assert.NoError(t, err)

			defaultConfig := generateDefaultConfig(t)
			defaultConfig.DB.Password = "db_password"
			assert.Exactly(t, result, defaultConfig, "result should equal to original configuration")
		}
	} else {
		t.Errorf("failed to marshal config")
	}
}

func TestAnonymousStruct(t *testing.T) {
	config := generateDefaultConfig(t)

	if bytes, err := json.Marshal(config); err == nil {
		if file, err := os.CreateTemp("/tmp", "confy"); err == nil {
			defer file.Close()
			defer os.Remove(file.Name())
			file.Write(bytes)
			var result testConfig
			os.Setenv("CONFY_DESCRIPTION", "environment description")
			defer os.Setenv("CONFY_DESCRIPTION", "")

			dir, filename := filepath.Split(file.Name())
			confy.DefaultConfy = confy.NewConfy(confy.WithFS(os.DirFS(dir)))
			err = confy.Load(&result, filename)
			assert.NoError(t, err)

			defaultConfig := generateDefaultConfig(t)
			defaultConfig.Anonymous.Description = "environment description"
			assert.Exactly(t, result, defaultConfig, "result should equal to original configuration")
		}
	} else {
		t.Errorf("failed to marshal config")
	}
}

type slicetestConfig struct {
	Test1 int
	Test2 []struct {
		Test2Ele1 int
		Test2Ele2 int
	}
}

func TestSliceFromEnv(t *testing.T) {
	var tc = slicetestConfig{
		Test1: 1,
		Test2: []struct {
			Test2Ele1 int
			Test2Ele2 int
		}{
			{
				Test2Ele1: 1,
				Test2Ele2: 2,
			},
			{
				Test2Ele1: 3,
				Test2Ele2: 4,
			},
		},
	}

	var result slicetestConfig
	os.Setenv("CONFY_TEST1", "1")
	os.Setenv("CONFY_TEST2_0_TEST2ELE1", "1")
	os.Setenv("CONFY_TEST2_0_TEST2ELE2", "2")

	os.Setenv("CONFY_TEST2_1_TEST2ELE1", "3")
	os.Setenv("CONFY_TEST2_1_TEST2ELE2", "4")
	err := confy.Load(&result)
	if err != nil {
		t.Fatalf("load from env err:%v", err)
	}

	assert.Exactly(t, result, tc, "result should equal to original configuration")
}

func TestConfigFromEnv(t *testing.T) {
	type config struct {
		LineBreakString string `required:"true"`
		Count           int64
		Slient          bool
		HomeAddress     struct {
			StreetName string
			City       string
		}
	}

	cfg := &config{}

	os.Setenv("CONFY_ENV_PREFIX", "CONFY")
	os.Setenv("CONFY_LINE_BREAK_STRING", "Line one\nLine two\nLine three\nAnd more lines")
	os.Setenv("CONFY_SLIENT", "1")
	os.Setenv("CONFY_COUNT", "10")
	os.Setenv("CONFY_HOME_ADDRESS_STREET_NAME", "abc")
	confy.Load(cfg)

	t.Log(cfg)

	if os.Getenv("CONFY_LINE_BREAK_STRING") != cfg.LineBreakString {
		t.Error("Failed to load value has line break from env")
	}

	if !cfg.Slient {
		t.Error("Failed to load bool from env")
	}

	if cfg.Count != 10 {
		t.Error("Failed to load number from env")
	}

	if os.Getenv("CONFY_HOME_ADDRESS_STREET_NAME") != cfg.HomeAddress.StreetName {
		t.Error("Failed to load StreetName from env")
	}
}

func TestValidation(t *testing.T) {
	type config struct {
		Name     string `validate:"-"`
		Title    string `validate:"alphanum,required"`
		AuthorIP string `validate:"ipv4"`
		Email    string `validate:"email"`
		Email2   string `validate:"email"`
		Endpoint string `validate:"required"`
		Count    int64  `validate:"required"`
		Slient   bool   `validate:"required"`
	}

	cfg := &config{Email: "a@b.com", Email2: "", AuthorIP: "1.1"}
	err := confy.Load(cfg)
	fmt.Printf("%+v\n", cfg)
	if err != nil {
		errs := err.(validator.ValidationErrors)
		for index, err := range errs {
			fmt.Printf("\t%d.  %s\n", index, err)
		}
		// t.Error("Error validating")
	}
}

func TestValidationMore(t *testing.T) {

	// usage : https://github.com/go-playground/validator/blob/master/doc.go
	type Address struct {
		// Skip Field: Usage: -
		Street string `validate:"-"`
		Zip    string `json:"zip" validate:"numeric,required"`
	}

	type User struct {
		Name           string `validate:"required"`
		Email          string `validate:"required,email"`
		Password       string `validate:"required"`
		Age            int    `validate:"required,numeric,gte=0,lte=130"`
		FavouriteColor string `validate:"hexcolor|rgb|rgba"`
		Home           *Address
		AddArray       []Address `validate:"unique,gt=0,dive,required"`
		// Multidimensional nesting is also supported, each level you wish to dive will
		// require another dive tag. dive has some sub-tags, 'keys' & 'endkeys', please see
		// the Keys & EndKeys section just below.
		// gt=0 will be applied to the map itself
		// len=4 will be applied to the map keys
		// required will be applied to map values
		AddMap map[string]Address `validate:"required,unique=zip,gt=1,dive,keys,required,len=4,endkeys,required"`
		// omitempty
		// Allows conditional validation, for example if a field is not set with
		// a value (Determined by the "required" validator) then other validation
		// such as min or max won't run, but if a value is set validation will run.
		Sleep time.Duration `validate:"omitempty,gt=1h30m"`
	}

	addMap := map[string]Address{
		"home": {"", "ABC456D89"},
		"work": {"", "12345"},
		"wor2": {"", "12345"},
	}
	var tests = []struct {
		param    interface{}
		expected bool
	}{
		{&User{"John", "john@yahoo.com", "123G#678", 20, "#010", &Address{"", "ABC456D89"}, []Address{{"Street", "123456"}, {"Street", "54321"}}, addMap, time.Hour * 2}, false},
		{&User{"John", "john!yahoo.com", "12345678", 20, "#001", &Address{"Street", "ABC456D89"}, []Address{{"Street", "ABC456D89"}, {"Street", "123456"}}, addMap, 0}, false},
		{&User{"John", "", "12345", -1, "rgb(255,255,255)", &Address{"Street", "123456789"}, []Address{{"Street", "ABC456D89"}, {"Street", "123456"}}, addMap, 0}, false},
		{&User{"", "john@yahoo.com", "123G#678", 20, "#000", &Address{"Street", "95504"}, []Address{{"Street", "123456"}, {"Street", "A123456"}}, nil, time.Minute}, false},
	}
	for _, test := range tests {
		err := confy.Load(test.param)
		if err != nil {
			t.Logf("Error for: %#v", test.param)
			// this check is only needed when your code could produce
			// an invalid value for validation such as interface with nil
			// value most including myself do not usually have code like this.
			if _, ok := err.(*validator.InvalidValidationError); ok {
				t.Log(err)
			}
			for _, err := range err.(validator.ValidationErrors) {
				t.Logf("Error: %v, Value: %v", err, err.Value())
			}
			if test.expected {
				t.Errorf("Got Error: %s", err)
			}
		}
		t.Log("-----------------------")
	}
}

func TestUsePkger(t *testing.T) {
	config := generateDefaultConfig(t)
	if bytes, err := json.Marshal(config); err == nil {
		if file, err := os.CreateTemp("..", "temp_confy"); err == nil {
			t.Cleanup(func() {
				t.Log("cleanup...")
				if err := file.Close(); err != nil {
					t.Error(err)
				}
				if err := os.Remove(file.Name()); err != nil {
					t.Error(err)
				}
			})

			_, _ = file.Write(bytes)

			var result testConfig
			dir, filename := filepath.Split(file.Name())
			confy.DefaultConfy = confy.NewConfy(confy.WithFS(os.DirFS(dir)))
			err = confy.Load(&result, filename)
			assert.NoError(t, err)
			assert.Exactly(t, result, config, "result should equal to original configuration")
		}
	} else {
		t.Errorf("failed to marshal config")
	}
}
