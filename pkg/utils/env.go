package utils

import (
	"os"
	"strings"
)

func GetEnvVarValue(key string, allowEmpty bool) string {
	value := os.Getenv(key)
	if value == "" && !allowEmpty {
		panic("'" + key + "' environment variable not set")
	}
	// Trim quote from both ends to avoid issues with string values
	// (sometime they are included in quotes, sometimes not)
	value = strings.Trim(value, "\"")
	value = strings.Trim(value, "'")
	return value
}
