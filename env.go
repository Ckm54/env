package env

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type envVar struct {
	value        interface{}
	name         string
	varType      string
	required     bool
	defaultValue interface{}
	help         string
	setValue     func(interface{}, string) error
	setDefault   func(interface{}, interface{})
	envValue     *string
}

var envs []envVar
var help = flag.Bool("help", false, "--help to show help")

func init() {
	envs = make([]envVar, 0)
}

// String registers a new string environment variable with the specified parameters.
// It appends the variable configuration to the global envs slice and returns a pointer
// to the string value that will be populated when environment variables are parsed.
//
// Parameters:
//   - name: The environment variable name to look for
//   - required: Whether this environment variable must be set
//   - defaultValue: Value to use if environment variable is not set
//   - help: Description of this environment variable for documentation
//
// Example usage:
//
//	port := env.String("PORT", false, "8080", "HTTP server port")
//
// The returned pointer will be populated with the environment variable value
// after calling env.Parse()
func String(name string, required bool, defaultValue, help string) *string {
	v := new(string)

	envs = append(envs, envVar{
		v,
		name,
		"string",
		required,
		defaultValue,
		help,
		func(a interface{}, b string) error {
			*a.(*string) = b
			return nil
		},
		func(a interface{}, b interface{}) {
			*a.(*string) = b.(string)
		},
		new(string),
	})

	return v
}

// Int registers a new integer environment variable with the specified parameters.
// It appends the variable configuration to the global envs slice and returns a pointer
// to the integer value that will be populated when environment variables are parsed.
//
// The function handles string to integer conversion internally and supports values
// up to 64-bit integers. If the environment variable contains an invalid integer,
// Parse() will return an error.
//
// Parameters:
//   - name: The environment variable name to look for
//   - required: Whether this environment variable must be set
//   - defaultValue: Value to use if environment variable is not set
//   - help: Description of this environment variable for documentation
//
// Example usage:
//
//	port := env.Int("PORT", false, 8080, "HTTP server port")
//
// The returned pointer will be populated with the environment variable value
// after calling env.Parse()
func Int(name string, required bool, defaultValue int, help string) *int {
	v := new(int)

	envs = append(envs, envVar{
		v,
		name,
		"integer",
		required,
		defaultValue,
		help,
		func(a interface{}, b string) error {
			v, err := strconv.ParseInt(b, 10, 64)
			if err != nil {
				a = nil
				return err
			}

			*a.(*int) = int(v)

			return nil
		},
		func(a interface{}, b interface{}) {
			if val, ok := b.(int); ok {
				*a.(*int) = val
			}
		},
		new(string),
	})

	return v
}

// Float64 registers a float64 environment variable with the given name, default value, and metadata.
//
// Parameters:
//   - name: Environment variable name.
//   - required: Whether the variable is mandatory.
//   - defaultValue: Default value if not set.
//   - help: Description for documentation.
//
// Example:
//   timeout := env.Float64("TIMEOUT", false, 30.0, "Request timeout in seconds")

func Float64(name string, required bool, defaultValue float64, help string) *float64 {
	v := new(float64)

	envs = append(envs, envVar{
		v,
		name,
		"float",
		required,
		defaultValue,
		help,
		func(i interface{}, s string) error {
			v, err := strconv.ParseFloat(s, 64)
			if err != nil {
				i = nil
				return err
			}

			*i.(*float64) = float64(v)
			return nil
		},
		func(i1, i2 interface{}) {
			*i1.(*float64) = i2.(float64)
		},
		new(string),
	})

	return v
}

// Bool registers a boolean environment variable with the given name, default value, and metadata.
// It appends the variable configuration to the global envs slice and returns a pointer to the boolean value.
//
// Parameters:
//   - name: Environment variable name.
//   - required: Whether the variable is mandatory.
//   - defaultValue: Default value if the variable is not set.
//   - help: Description of the variable for documentation.
//
// Example:
//
//	debugMode := env.Bool("DEBUG_MODE", false, false, "Enable debug mode")
func Bool(name string, required bool, defaultValue bool, help string) *bool {
	v := new(bool)

	envs = append(envs, envVar{
		v,
		name,
		"boolean",
		required,
		defaultValue,
		help,
		func(i interface{}, s string) error {
			v, err := strconv.ParseBool(s)
			if err != nil {
				i = nil
				return err
			}

			*i.(*bool) = bool(v)

			return nil
		},
		func(i1, i2 interface{}) {
			*i1.(*bool) = i2.(bool)
		},
		new(string),
	})

	return v
}

func Parse() error {
	// parse the maub flags package to enable the --help function
	flag.Parse()
	if *help == true {
		fmt.Println("Config values are set using environment variables. For more info please see the following list.")
		fmt.Println("")
		fmt.Println(Help())

		os.Exit(0)
	}

	errors := make([]string, 0)

	for _, e := range envs {
		err := processEnvVar(e)
		if err != nil {
			errors = append(errors, fmt.Sprintf("expected: %s type: %s got: %s", e.name, e.varType, *e.envValue))
		}
	}

	if len(errors) > 0 {
		errString := strings.Join(errors, "\n")
		return fmt.Errorf(errString)
	}

	return nil
}

func processEnvVar(e envVar) error {
	*e.envValue = os.Getenv(e.name)
	if *e.envValue == "" && !e.required {
		e.setDefault(e.value, e.defaultValue)
		return nil
	}

	if *e.envValue == "" && e.required {
		return fmt.Errorf("%s should be provided", e.name)
	}

	err := e.setValue(e.value, *e.envValue)
	if err != nil {
		return err
	}

	return nil
}

// Help is help
func Help() string {
	h := make([]string, 1)
	h[0] = "Environment variables:"

	for _, e := range envs {
		def := fmt.Sprintf("'%v'", e.defaultValue)
		if def == "''" {
			def = "no default"
		}

		h = append(h, "  "+e.name+" default: "+def)
		h = append(h, "       ")
	}

	return strings.Join(h, "\n")
}
