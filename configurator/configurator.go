package configurator

import (
	"context"
	"fmt"
	"os"
	"reflect"
	"regexp"
	"strconv"

	"github.com/go-playground/validator/v10"
)

var (
	// Default Configurator
	DefaultConfigurator Configurator
)

type Configurator interface {
	GetEnvironment() string
	Load(config interface{}, files ...string) error
}

type configurator struct {
	config
	validate *validator.Validate
}

func (c *configurator) init() {

	if env := os.Getenv("CONFIG_ENV"); env != "" {
		c.config.environment = env
	}

	if envPrefix := os.Getenv("CONFIG_ENV_PREFIX"); envPrefix != "" {
		c.config.environmentVariablePrefix = envPrefix
	}

	if debugMode, _ := strconv.ParseBool(os.Getenv("CONFIG_DEBUG_MODE")); debugMode {
		c.config.debug = debugMode
	}

	if verboseMode, _ := strconv.ParseBool(os.Getenv("CONFIG_VERBOSE_MODE")); verboseMode {
		c.config.verbose = verboseMode
	}

	if silentMode, _ := strconv.ParseBool(os.Getenv("CONFIG_SILENT_MODE")); silentMode {
		c.config.silent = silentMode
	}

	if usePkger, _ := strconv.ParseBool(os.Getenv("CONFIG_USE_PKGER")); usePkger {
		c.config.usePkger = usePkger
	}

	c.validate = validator.New()
	// c.validate.SetTagName("valid")
}

// NewConfigurator creates a new configurator configured with the given options.
func NewConfigurator(opts ...Option) Configurator {
	// Set default config
	environment := "development"
	if testRegexp.MatchString(os.Args[0]) {
		environment = "test"
	}
	cfg := config{environment: environment, environmentVariablePrefix: "CONFIG", context: context.Background()}
	cfg.options(opts...)
	configor := &configurator{config: cfg}
	configor.init()
	return configor
}

var testRegexp = regexp.MustCompile("_test|(\\.test$)")

// GetEnvironment return runtime environment
func (c *configurator) GetEnvironment() string {
	return c.config.environment
}

// Load will unmarshal configurations to struct from files that you provide
func (c *configurator) Load(config interface{}, files ...string) (err error) {
	defaultValue := reflect.Indirect(reflect.ValueOf(config))
	if !defaultValue.CanAddr() {
		return fmt.Errorf("config %v should be addressable", config)
	}
	err = c.load(config, files...)
	return
}

// ENV return environment
func GetEnvironment() string {
	return DefaultConfigurator.GetEnvironment()
}

// Load will unmarshal configurations to struct from files that you provide
func Load(config interface{}, files ...string) error {
	return DefaultConfigurator.Load(config, files...)
}
