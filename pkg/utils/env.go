package utils

import (
	"os"
)

func GetEnvVarValue(key string, allowEmpty bool) string {
	value := os.Getenv(key)
	if value == "" && !allowEmpty {
		panic("'" + key + "' environment variable not set")
	}
	return value
}
