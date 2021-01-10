package configurator_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/go-playground/validator/v10"
	"gopkg.in/yaml.v2"

	"github.com/xmlking/toolkit/configurator"
)

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
	APPName   string   `default:"configurator" yaml:",omitempty" json:",omitempty"`
	Hosts     []string `validate:"omitempty,dive,url" default:"[\"https://abc.org\"]"`
	DB        *Database
	Contacts  []Contact
	Anonymous `anonymous:"true" default:"-"`
	private   string
}

func generateDefaultConfig() testConfig {
	return testConfig{
		APPName: "configurator",
		Hosts:   []string{"http://example.org", "http://jinzhu.me"},
		DB: &Database{
			Name:     "configurator",
			User:     "configurator",
			Password: "configurator",
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
	configurator.DefaultConfigurator = configurator.NewConfigurator()
	fmt.Println("Setup completed")
}

func teardown() {
	// Do something here.
	configurator.DefaultConfigurator = nil
	fmt.Println("Teardown completed")
}

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	teardown()
	os.Exit(code)
}

func TestLoadNormaltestConfig(t *testing.T) {

	config := generateDefaultConfig()
	if bytes, err := json.Marshal(config); err == nil {
		if file, err := ioutil.TempFile("/tmp", "configurator"); err == nil {
			defer file.Close()
			defer os.Remove(file.Name())
			file.Write(bytes)

			var result testConfig
			configurator.Load(&result, file.Name())
			if !reflect.DeepEqual(result, config) {
				t.Errorf("result should equal to original configuration")
			}
		}
	} else {
		t.Errorf("failed to marshal config")
	}
}

// CONFIG_DEBUG_MODE=true CONFIG_VERBOSE_MODE=true go test -v   -run TestDefaultValue -count=1
func TestDefaultValue(t *testing.T) {
	config := generateDefaultConfig()
	config.APPName = ""
	config.DB.Port = 0
	config.DB.SSL = false
	config.Contacts[0].Name = ""

	if bytes, err := json.Marshal(config); err == nil {
		if file, err := ioutil.TempFile("/tmp", "configurator"); err == nil {
			defer file.Close()
			defer os.Remove(file.Name())
			file.Write(bytes)

			var result testConfig
			if err := configurator.Load(&result, file.Name()); err != nil {
				t.Error(err)
			}

			if !reflect.DeepEqual(result, generateDefaultConfig()) {
				t.Errorf("\nExpected: %+v, \nGot: %+v", generateDefaultConfig(), result)
			}
		}
	} else {
		t.Errorf("failed to marshal config")
	}
}

func TestMissingRequiredValue(t *testing.T) {
	config := generateDefaultConfig()
	config.DB.Password = ""

	if bytes, err := json.Marshal(config); err == nil {
		if file, err := ioutil.TempFile("/tmp", "configurator"); err == nil {
			defer file.Close()
			defer os.Remove(file.Name())
			file.Write(bytes)

			var result testConfig
			if err := configurator.Load(&result, file.Name()); err == nil {
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

	file, err := ioutil.TempFile("/tmp", "configurator")
	if err != nil {
		t.Fatal("Could not create temp file")
	}

	defer os.Remove(file.Name())
	defer file.Close()

	filename := file.Name()

	if data, err := yaml.Marshal(config); err == nil {
		file.WriteString(string(data))

		var result configStruct

		// Do not return error when there are unmatched keys but ErrorOnUnmatchedKeys is false
		if err := configurator.NewConfigurator().Load(&result, filename); err != nil {
			t.Errorf("Should NOT get error when loading configuration with extra keys. Error: %v", err)
		}

		// Return an error when there are unmatched keys and ErrorOnUnmatchedKeys is true
		if err := configurator.NewConfigurator(configurator.WithErrorOnUnmatchedKeys()).Load(&result, filename); err == nil {
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

	// Do not return error when there are unmatched keys but ErrorOnUnmatchedKeys is false
	if err := configurator.NewConfigurator().Load(&result, filename); err != nil {
		t.Errorf("Should NOT get error when loading configuration with extra keys. Error: %v", err)
	}

	// Return an error when there are unmatched keys and ErrorOnUnmatchedKeys is true
	if err := configurator.NewConfigurator(configurator.WithErrorOnUnmatchedKeys()).Load(&result, filename); err == nil {
		t.Errorf("Should get error when loading configuration with extra keys")

		// The error should be of type *yaml.TypeError
	} else if _, ok := err.(*yaml.TypeError); !ok {
		// || !strings.Contains(err.Error(), "not found in struct") {
		t.Errorf("Error should be of type yaml.TypeError. Instead error is %v", err)
	}
}

func TestYamlDefaultValue(t *testing.T) {
	config := generateDefaultConfig()
	config.APPName = ""
	config.DB.Port = 0
	config.DB.SSL = false

	if bytes, err := yaml.Marshal(config); err == nil {
		if file, err := ioutil.TempFile("/tmp", "configurator.*.yaml"); err == nil {
			defer file.Close()
			defer os.Remove(file.Name())
			file.Write(bytes)

			var result testConfig
			configurator.Load(&result, file.Name())

			if !reflect.DeepEqual(result, generateDefaultConfig()) {
				t.Errorf("result should be set default value correctly")
			}
		}
	} else {
		t.Errorf("failed to marshal config")
	}
}

func TestLoadConfigurationByEnvironment(t *testing.T) {
	config := generateDefaultConfig()
	config2 := struct {
		APPName string
	}{
		APPName: "config2",
	}

	if file, err := ioutil.TempFile("/tmp", "configurator"); err == nil {
		defer file.Close()
		defer os.Remove(file.Name())
		configBytes, _ := yaml.Marshal(config)
		config2Bytes, _ := yaml.Marshal(config2)
		ioutil.WriteFile(file.Name()+".yaml", configBytes, 0644)
		defer os.Remove(file.Name() + ".yaml")
		ioutil.WriteFile(file.Name()+".production.yaml", config2Bytes, 0644)
		defer os.Remove(file.Name() + ".production.yaml")

		var result testConfig
		os.Setenv("CONFIG_ENV", "production")
		defer os.Setenv("CONFIG_ENV", "")
		configurator.DefaultConfigurator = configurator.NewConfigurator()
		if err := configurator.Load(&result, file.Name()+".yaml"); err != nil {
			t.Errorf("No error should happen when load configurations, but got %v", err)
		}

		var defaultConfig = generateDefaultConfig()
		defaultConfig.APPName = "config2"
		if !reflect.DeepEqual(result, defaultConfig) {
			t.Errorf("result should load configurations by environment correctly")
		}
	}

	t.Cleanup(func() {
		t.Log("Cleanup...")
	})
}

func TestLoadtestConfigurationByEnvironmentSetBytestConfig(t *testing.T) {
	config := generateDefaultConfig()
	config2 := struct {
		APPName string
	}{
		APPName: "production_config2",
	}

	if file, err := ioutil.TempFile("/tmp", "configurator"); err == nil {
		defer file.Close()
		defer os.Remove(file.Name())
		configBytes, _ := yaml.Marshal(config)
		config2Bytes, _ := yaml.Marshal(config2)
		ioutil.WriteFile(file.Name()+".yaml", configBytes, 0644)
		defer os.Remove(file.Name() + ".yaml")
		ioutil.WriteFile(file.Name()+".production.yaml", config2Bytes, 0644)
		defer os.Remove(file.Name() + ".production.yaml")

		var result testConfig
		configurator := configurator.NewConfigurator(configurator.WithEnvironment("production"))
		if configurator.Load(&result, file.Name()+".yaml"); err != nil {
			t.Errorf("No error should happen when load configurations, but got %v", err)
		}

		var defaultConfig = generateDefaultConfig()
		defaultConfig.APPName = "production_config2"
		if !reflect.DeepEqual(result, defaultConfig) {
			t.Errorf("result should be load configurations by environment correctly")
		}

		if configurator.GetEnvironment() != "production" {
			t.Errorf("configurator's environment should be production")
		}
	}
}

func TestOverwritetestConfigurationWithEnvironmentWithDefaultPrefix(t *testing.T) {
	config := generateDefaultConfig()

	if bytes, err := json.Marshal(config); err == nil {
		if file, err := ioutil.TempFile("/tmp", "configurator"); err == nil {
			defer file.Close()
			defer os.Remove(file.Name())
			file.Write(bytes)
			var result testConfig
			os.Setenv("CONFIG_APP_NAME", "config2")
			os.Setenv("CONFIG_HOSTS", "- http://example.org\n- http://jinzhu.me")
			os.Setenv("CONFIG_DB_NAME", "db_name")
			defer os.Setenv("CONFIG_APP_NAME", "")
			defer os.Setenv("CONFIG_HOSTS", "")
			defer os.Setenv("CONFIG_DB_NAME", "")
			configurator.Load(&result, file.Name())

			var defaultConfig = generateDefaultConfig()
			defaultConfig.APPName = "config2"
			defaultConfig.Hosts = []string{"http://example.org", "http://jinzhu.me"}
			defaultConfig.DB.Name = "db_name"
			if !reflect.DeepEqual(result, defaultConfig) {
				t.Errorf("result should equal to original configuration")
			}
		}
	}
}

func TestOverwritetestConfigurationWithEnvironment(t *testing.T) {
	config := generateDefaultConfig()

	if bytes, err := json.Marshal(config); err == nil {
		if file, err := ioutil.TempFile("/tmp", "configurator"); err == nil {
			defer file.Close()
			defer os.Remove(file.Name())
			file.Write(bytes)
			var result testConfig
			os.Setenv("CONFIG_ENV_PREFIX", "app")
			os.Setenv("APP_APP_NAME", "config2")
			os.Setenv("APP_DB_NAME", "db_name")
			defer os.Setenv("CONFIG_ENV_PREFIX", "")
			defer os.Setenv("APP_APP_NAME", "")
			defer os.Setenv("APP_DB_NAME", "")
			configurator.DefaultConfigurator = configurator.NewConfigurator()
			configurator.Load(&result, file.Name())

			var defaultConfig = generateDefaultConfig()
			defaultConfig.APPName = "config2"
			defaultConfig.DB.Name = "db_name"
			if !reflect.DeepEqual(result, defaultConfig) {
				t.Errorf("result should equal to original configuration")
			}
		}
	}
}

func TestOverwritetestConfigurationWithEnvironmentThatSetBytestConfig(t *testing.T) {
	config := generateDefaultConfig()

	if bytes, err := json.Marshal(config); err == nil {
		if file, err := ioutil.TempFile("/tmp", "configurator"); err == nil {
			defer file.Close()
			defer os.Remove(file.Name())
			file.Write(bytes)
			os.Setenv("APP1_APP_NAME", "config2")
			os.Setenv("APP1_DB_NAME", "db_name")
			defer os.Setenv("APP1_APP_NAME", "")
			defer os.Setenv("APP1_DB_NAME", "")

			var result testConfig
			configurator := configurator.NewConfigurator(configurator.WithEnvironmentVariablePrefix("APP1"))
			if err := configurator.Load(&result, file.Name()); err != nil {
				t.Error(err)
			}

			var defaultConfig = generateDefaultConfig()
			defaultConfig.APPName = "config2"
			defaultConfig.DB.Name = "db_name"
			if !reflect.DeepEqual(result, defaultConfig) {
				t.Errorf("result should equal to original configuration")
			}
		}
	}
}

func TestResetPrefixToBlank(t *testing.T) {
	config := generateDefaultConfig()

	if bytes, err := json.Marshal(config); err == nil {
		if file, err := ioutil.TempFile("/tmp", "configurator"); err == nil {
			defer file.Close()
			defer os.Remove(file.Name())
			file.Write(bytes)
			var result testConfig
			os.Setenv("CONFIG_ENV_PREFIX", "-")
			os.Setenv("APP_NAME", "config2")
			os.Setenv("DB_NAME", "db_name")
			defer os.Setenv("CONFIG_ENV_PREFIX", "")
			defer os.Setenv("APP_NAME", "")
			defer os.Setenv("DB_NAME", "")
			configurator.DefaultConfigurator = configurator.NewConfigurator()

			configurator.Load(&result, file.Name())

			var defaultConfig = generateDefaultConfig()
			defaultConfig.APPName = "config2"
			defaultConfig.DB.Name = "db_name"
			if !reflect.DeepEqual(result, defaultConfig) {
				t.Errorf("result should equal to original configuration")
			}
		}
	}
}

func TestResetPrefixToBlank2(t *testing.T) {
	config := generateDefaultConfig()

	if bytes, err := json.Marshal(config); err == nil {
		if file, err := ioutil.TempFile("/tmp", "configurator"); err == nil {
			defer file.Close()
			defer os.Remove(file.Name())
			file.Write(bytes)
			var result testConfig
			os.Setenv("CONFIG_ENV_PREFIX", "-")
			os.Setenv("APP_NAME", "config2")
			os.Setenv("DB_NAME", "db_name")
			defer os.Setenv("CONFIG_ENV_PREFIX", "")
			defer os.Setenv("APPName", "")
			defer os.Setenv("DB_NAME", "")

			configurator.DefaultConfigurator = configurator.NewConfigurator()
			configurator.Load(&result, file.Name())

			var defaultConfig = generateDefaultConfig()
			defaultConfig.APPName = "config2"
			defaultConfig.DB.Name = "db_name"
			if !reflect.DeepEqual(result, defaultConfig) {
				t.Errorf("result should equal to original configuration")
			}
		}
	}
}

func TestReadFromEnvironmentWithSpecifiedEnvName(t *testing.T) {
	config := generateDefaultConfig()

	if bytes, err := json.Marshal(config); err == nil {
		if file, err := ioutil.TempFile("/tmp", "configurator"); err == nil {
			defer file.Close()
			defer os.Remove(file.Name())
			file.Write(bytes)
			var result testConfig
			os.Setenv("DBPassword", "db_password")
			defer os.Setenv("DBPassword", "")

			configurator.DefaultConfigurator = configurator.NewConfigurator()
			configurator.Load(&result, file.Name())

			var defaultConfig = generateDefaultConfig()
			defaultConfig.DB.Password = "db_password"
			if !reflect.DeepEqual(result, defaultConfig) {
				t.Errorf("result should equal to original configuration")
			}
		}
	}
}

func TestAnonymousStruct(t *testing.T) {
	config := generateDefaultConfig()

	if bytes, err := json.Marshal(config); err == nil {
		if file, err := ioutil.TempFile("/tmp", "configurator"); err == nil {
			defer file.Close()
			defer os.Remove(file.Name())
			file.Write(bytes)
			var result testConfig
			os.Setenv("CONFIG_DESCRIPTION", "environment description")
			defer os.Setenv("CONFIG_DESCRIPTION", "")

			configurator.DefaultConfigurator = configurator.NewConfigurator()
			configurator.Load(&result, file.Name())

			var defaultConfig = generateDefaultConfig()
			defaultConfig.Anonymous.Description = "environment description"
			if !reflect.DeepEqual(result, defaultConfig) {
				t.Errorf("result should equal to original configuration")
			}
		}
	}
}

func TestENV(t *testing.T) {
	if configurator.GetEnvironment() != "test" {
		t.Errorf("Env should be test when running `go test`, instead env is %v", configurator.GetEnvironment())
	}

	os.Setenv("CONFIG_ENV", "production")
	defer os.Setenv("CONFIG_ENV", "")
	configurator.DefaultConfigurator = configurator.NewConfigurator()
	if configurator.GetEnvironment() != "production" {
		t.Errorf("Env should be production when set it with CONFIG_ENV")
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
	os.Setenv("CONFIG_TEST1", "1")
	os.Setenv("CONFIG_TEST2_0_TEST2ELE1", "1")
	os.Setenv("CONFIG_TEST2_0_TEST2ELE2", "2")

	os.Setenv("CONFIG_TEST2_1_TEST2ELE1", "3")
	os.Setenv("CONFIG_TEST2_1_TEST2ELE2", "4")
	err := configurator.Load(&result)
	if err != nil {
		t.Fatalf("load from env err:%v", err)
	}

	if !reflect.DeepEqual(result, tc) {
		t.Fatalf("unexpected result:%+v", result)
	}
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

	os.Setenv("CONFIG_ENV_PREFIX", "CONFIG")
	os.Setenv("CONFIG_LINE_BREAK_STRING", "Line one\nLine two\nLine three\nAnd more lines")
	os.Setenv("CONFIG_SLIENT", "1")
	os.Setenv("CONFIG_COUNT", "10")
	os.Setenv("CONFIG_HOME_ADDRESS_STREET_NAME", "abc")
	configurator.Load(cfg)

	t.Log(cfg)

	if os.Getenv("CONFIG_LINE_BREAK_STRING") != cfg.LineBreakString {
		t.Error("Failed to load value has line break from env")
	}

	if !cfg.Slient {
		t.Error("Failed to load bool from env")
	}

	if cfg.Count != 10 {
		t.Error("Failed to load number from env")
	}

	if os.Getenv("CONFIG_HOME_ADDRESS_STREET_NAME") != cfg.HomeAddress.StreetName {
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
	err := configurator.Load(cfg)
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
		"home": Address{"", "ABC456D89"},
		"work": Address{"", "12345"},
		"wor2": Address{"", "12345"},
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
		err := configurator.Load(test.param)
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
	config := generateDefaultConfig()
	if bytes, err := json.Marshal(config); err == nil {
		if file, err := ioutil.TempFile(".", "temp_configurator"); err == nil {
			defer file.Close()
			defer os.Remove(file.Name())
			file.Write(bytes)

			var result testConfig
			configurator.NewConfigurator(configurator.WithPkger()).Load(&result, "/"+file.Name())
			if !reflect.DeepEqual(result, config) {
				t.Errorf("result should equal to original configuration")
			}
		}
	} else {
		t.Errorf("failed to marshal config")
	}
}
