package env

import "flag"

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
