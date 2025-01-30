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
	// Create a new string pointer to store the variable value.
	v := new(string)

	// Append a new environment variable definition to `envs`.
	envs = append(envs, envVar{
		v,            // Pointer to the string variable.
		name,         // The name of the environment variable.
		"string",     // The data type (for documentation/help purposes).
		required,     // Whether the variable is required.
		defaultValue, // The default value if the variable is not set.
		help,         // Help text describing the variable.

		// Function to parse and set the string value from input.
		func(a interface{}, b string) error {
			*a.(*string) = b // Directly assign the string value.
			return nil
		},

		// Function to set the default value if the environment variable is not set.
		func(a interface{}, b interface{}) {
			*a.(*string) = b.(string) // Assign the default value, ensuring it's a string.
		},

		new(string), // Pointer to store the raw string representation of the environment variable.
	})

	// Return the pointer to the string variable so it can be accessed elsewhere.
	return v
}

// Int defines an integer environment variable with the specified parameters.
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
	// Create a new integer pointer to store the variable value.
	v := new(int)

	// Append a new environment variable definition to `envs`.
	envs = append(envs, envVar{
		v,            // Pointer to the integer variable.
		name,         // The name of the environment variable.
		"integer",    // The data type (for documentation/help purposes).
		required,     // Whether the variable is required.
		defaultValue, // The default value if the variable is not set.
		help,         // Help text describing the variable.

		// Function to parse and set the integer value from a string.
		func(a interface{}, b string) error {
			v, err := strconv.ParseInt(b, 10, 64) // Convert string to int64.
			if err != nil {
				a = nil // If parsing fails, set `a` to nil.
				return err
			}

			*a.(*int) = int(v) // Store the parsed value as an int.
			return nil
		},

		// Function to set the default value if the environment variable is not set.
		func(a interface{}, b interface{}) {
			if val, ok := b.(int); ok { // Ensure `b` is an int before assignment.
				*a.(*int) = val
			}
		},

		new(string), // Pointer to store the raw string representation of the environment variable.
	})

	// Return the pointer to the integer variable so it can be accessed elsewhere.
	return v
}

// Float64 defines a float64 environment variable, adds it to the list of expected environment variables (`envs`),
// and returns a pointer to its value.
//
// Parameters:
//   - name: Environment variable name.
//   - required: Whether the variable is mandatory.
//   - defaultValue: Default value if not set.
//   - help: Description for documentation.
//
// Example:
//
//	timeout := env.Float64("TIMEOUT", false, 30.0, "Request timeout in seconds")
func Float64(name string, required bool, defaultValue float64, help string) *float64 {
	// Create a new float64 pointer to store the variable value.
	v := new(float64)

	// Append a new environment variable definition to `envs`.
	envs = append(envs, envVar{
		v,            // Pointer to the float64 variable.
		name,         // The name of the environment variable.
		"float",      // The data type (for documentation/help purposes).
		required,     // Whether the variable is required.
		defaultValue, // The default value if the variable is not set.
		help,         // Help text describing the variable.

		// Function to parse and set the float64 value from a string.
		func(i interface{}, s string) error {
			v, err := strconv.ParseFloat(s, 64) // Convert string to float64.
			if err != nil {
				i = nil // If parsing fails, set `i` to nil.
				return err
			}

			*i.(*float64) = float64(v) // Store the parsed value.
			return nil
		},

		// Function to set the default value if the environment variable is not set.
		func(i1, i2 interface{}) {
			*i1.(*float64) = i2.(float64) // Assign default float64 value.
		},

		new(string), // Pointer to store the raw string representation of the environment variable.
	})

	// Return the pointer to the float64 variable so it can be accessed elsewhere.
	return v
}

// Bool defines a boolean environment variable, adds it to the list of expected environment variables (`envs`),
// and returns a pointer to its value.
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
	// Create a new boolean pointer to store the variable value.
	v := new(bool)

	// Append a new environment variable definition to `envs`.
	envs = append(envs, envVar{
		v,            // Pointer to the boolean variable.
		name,         // The name of the environment variable.
		"boolean",    // The data type (for documentation/help purposes).
		required,     // Whether the variable is required.
		defaultValue, // The default value if the variable is not set.
		help,         // Help text describing the variable.

		// Function to parse and set the boolean value from a string.
		func(i interface{}, s string) error {
			v, err := strconv.ParseBool(s) // Convert string to boolean.
			if err != nil {
				i = nil // If parsing fails, set `i` to nil.
				return err
			}

			*i.(*bool) = bool(v) // Store the parsed value.
			return nil
		},

		// Function to set the default value if the environment variable is not set.
		func(i1, i2 interface{}) {
			*i1.(*bool) = i2.(bool) // Assign default boolean value.
		},

		new(string), // Pointer to store the raw string representation of the environment variable.
	})

	// Return the pointer to the boolean variable so it can be accessed elsewhere.
	return v
}

// Parse processes command-line flags and environment variables.
func Parse() error {
	// Parse the main flags package to enable the --help function.
	flag.Parse()

	// If the --help flag is provided, print the help message and exit.
	if *help {
		fmt.Println("Config values are set using environment variables. For more info please see the following list.")
		fmt.Println("")
		fmt.Println(Help())

		// Exit the program after displaying help.
		os.Exit(0)
	}

	// Collect errors encountered while processing environment variables.
	errors := make([]string, 0)

	// Iterate through all expected environment variables.
	for _, e := range envs {
		err := processEnvVar(e)
		if err != nil {
			// Append an error message if the environment variable is invalid or missing.
			errors = append(errors, fmt.Sprintf("expected: %s type: %s got: %s", e.name, e.varType, *e.envValue))
		}
	}

	// If there were any errors, format and return them as a single error message.
	if len(errors) > 0 {
		errString := strings.Join(errors, "\n")
		return fmt.Errorf("%s", errString)
	}

	// Return nil if all environment variables were processed successfully.
	return nil
}

// processEnvVar retrieves and validates a single environment variable.
func processEnvVar(e envVar) error {
	// Get the environment variable value from the system.
	*e.envValue = os.Getenv(e.name)

	// If the variable is empty and it's not required, set its default value.
	if *e.envValue == "" && !e.required {
		e.setDefault(e.value, e.defaultValue)
		return nil
	}

	// If the variable is empty but required, return an error.
	if *e.envValue == "" && e.required {
		return fmt.Errorf("%s should be provided", e.name)
	}

	// Try setting the value using a method that processes it.
	err := e.setValue(e.value, *e.envValue)
	if err != nil {
		return err
	}

	// Return nil if everything is successful.
	return nil
}

// Help generates and returns a help message listing all environment variables.
func Help() string {
	// Initialize the help message with a title.
	h := make([]string, 1)
	h[0] = "Environment variables:"

	// Iterate through all environment variables to generate their descriptions.
	for _, e := range envs {
		def := fmt.Sprintf("'%v'", e.defaultValue)
		if def == "''" {
			def = "no default"
		}

		// Append the variable name and default value to the help message.
		h = append(h, "  "+e.name+" default: "+def)
		h = append(h, "       ") // Add a blank line for better readability.
	}

	// Join all the help message parts into a single string and return.
	return strings.Join(h, "\n")
}
