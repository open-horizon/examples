//
// envutil -- utilities for getting configuration from environment variables.
//
// Written by Glen Darling (glendarling@us.ibm.com), Oct. 2016
//
package envutil

import (
	"os"
	"strconv"
	"logutil"
)

// Provides global access to the env vars that were set. This struct is populated early in main()
type CfgType struct {
	GPS_PORT int
	USE_GPS bool
	LAT float64
	LON float64
	LOCATION_ACCURACY_KM float64
}
var Cfg CfgType

/* not used...
// Get a string value from the first prefixed environment variable found
// that matches, or "" if none is found.
func GetFirstMatchingVariable(prefixes []string, env_var string) (result string) {
	result = ""
	for _, prefix := range prefixes {
		full_name := fmt.Sprintf("%s_%s", prefix, env_var)
		result = os.Getenv(full_name)
		if "" != result { return }
	}
	return
}
*/

// Get a string value from an environment variable or return the passed default
// Emit a warning if the variable is not set and 'warn' is true
func GetString(env_var string, default_value string, warn bool) string {
	env_str := os.Getenv(env_var)
	if "" == env_str {
		if warn {
			logutil.Logf("WARNING: Env variable %s should have been set, but it was not.", env_var)
		}
		return default_value
	}
	return env_str
}

// Get a boolean value from an environment variable or return the passed default
// Emit a warning if the variable is not set and 'warn' is true
// Emit a warning if the value cannot be parsed as a bool
func GetBool(env_var string, default_value bool, warn bool) bool {
	env_str := os.Getenv(env_var)
	if "" == env_str {
		if warn {
			logutil.Logf("WARNING: Env variable %s should have been set, but it was not.", env_var)
		}
		return default_value
	}
	// Convert the string to bool. Can be 0, 1, true, false
	var env_bool bool
	env_int, bool_err := strconv.ParseInt(env_str, 10, 64)
	if bool_err == nil {
		env_bool = (0 != env_int)
	} else {
		env_bool, bool_err = strconv.ParseBool(env_str)
	}
	if bool_err != nil {
		logutil.Logf("WARNING: Unable to parse env variable %s as a boolean.  Found: '%q'.", env_var, env_str)
		return default_value
	}
	return env_bool
}

// Get an int value from an environment variable or return the passed default
// Emit a warning if the variable is not set and 'warn' is true
// Emit a warning if the value cannot be parsed as an int
func GetInt(env_var string, default_value int, warn bool) int {
	env_str := os.Getenv(env_var)
	if "" == env_str {
		if warn {
			logutil.Logf("WARNING: Env variable %s should have been set, but it was not.", env_var)
		}
		return default_value
	}
	env_int, int_err := strconv.ParseInt(env_str, 10, 64)
	if int_err != nil {
		logutil.Logf("WARNING: Unable to parse env variable %s as an integer.  Found: '%q'.", env_var, env_str)
		return default_value
	}
	return int(env_int)
}

// Get a float value from an environment variable or return the passed default
// Emit a warning if the variable is not set and 'warn' is true
// Emit a warning if the value cannot be parsed as a float64
func GetFloat(env_var string, default_value float64, warn bool) float64 {
	env_str := os.Getenv(env_var)
	if "" == env_str {
		if warn {
			logutil.Logf("WARNING: Env variable %s should have been set, but it was not.", env_var)
		}
		return default_value
	}
	env_float, float_err := strconv.ParseFloat(env_str, 64)
	if float_err != nil {
		logutil.Logf("WARNING: Unable to parse env variable %s as a float.  Found: '%q'.", env_var, env_str)
		return default_value
	}
	return env_float
}
