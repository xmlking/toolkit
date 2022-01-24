package confy

import (
	"context"
	"fmt"
	"os"
	"reflect"
	"regexp"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog/log"
)

var (
	// Default Confy
	DefaultConfy Confy
)

type Confy interface {
	GetEnvironment() string
	Load(config interface{}, files ...string) error
}

type confy struct {
	config
	validate *validator.Validate
}

func (c *confy) init() {

	if env := os.Getenv("CONFY_ENV"); env != "" {
		c.config.environment = env
	}

	if envPrefix := os.Getenv("CONFY_ENV_PREFIX"); envPrefix != "" {
		c.config.environmentVariablePrefix = envPrefix
	}

	if debugMode, _ := strconv.ParseBool(os.Getenv("CONFY_DEBUG_MODE")); debugMode {
		c.config.debug = debugMode
	}

	if verboseMode, _ := strconv.ParseBool(os.Getenv("CONFY_VERBOSE_MODE")); verboseMode {
		c.config.verbose = verboseMode
	}

	if silentMode, _ := strconv.ParseBool(os.Getenv("CONFY_SILENT_MODE")); silentMode {
		c.config.silent = silentMode
	}

	c.validate = validator.New()
	// c.validate.SetTagName("valid")
}

// NewConfy creates a new configurator configured with the given options.
func NewConfy(opts ...Option) Confy {
	// Set default config
	environment := "development"
	if testRegexp.MatchString(os.Args[0]) {
		environment = "test"
	}
	cfg := config{environment: environment, environmentVariablePrefix: "CONFY", context: context.Background(), fs: os.DirFS(".")}
	cfg.options(opts...)
	confy := &confy{config: cfg}
	confy.init()
	return confy
}

var testRegexp = regexp.MustCompile("_test|(\\.test$)")

// GetEnvironment return runtime environment
func (c *confy) GetEnvironment() string {
	return c.config.environment
}

// Load will unmarshal configurations to struct from files that you provide
func (c *confy) Load(config interface{}, files ...string) (err error) {
	defaultValue := reflect.Indirect(reflect.ValueOf(config))
	if !defaultValue.CanAddr() {
		return fmt.Errorf("config %v should be addressable", config)
	}
	err = c.load(config, files...)
	return
}

// GetEnvironment return environment
func GetEnvironment() string {
	if DefaultConfy == nil {
		log.Fatal().Msg("config not initialized yet...")
	}
	return DefaultConfy.GetEnvironment()
}

// Load will unmarshal configurations to struct from files that you provide
func Load(config interface{}, files ...string) error {
	return DefaultConfy.Load(config, files...)
}
