package env

import (
	"os"
	"strconv"
)

/*
* GetString retrieves a String from an environment variable.
* if the environment variable is not set or cannot be parsed as an String,
* it returns the provided fallback value.
 */
func GetString(key, fallback string) string {
	val, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}

	return val
}

/*
* GetInt retrieves an integer value from an environment variable.
* If the environment variable is not set or cannot be parsed as an integer,
* it returns the provided fallback value.
 */
func GetInt(key string, fallback int) int {
	val, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}

	valAsInt, err := strconv.Atoi(val)
	if err != nil {
		return fallback
	}

	return valAsInt
}
