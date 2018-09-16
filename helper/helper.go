// Package helper provides improved env functions
package helper

import (
	"errors"
	"fmt"
	"os"
)

// GetEnv function provides a safe lookup for environment variables
func GetEnv(key string) (string, error) {
	if len(key) == 0 {
		return "", errors.New("Env variable name must be provided to getEnv")
	}
	val, ok := os.LookupEnv(key)
	if !ok {
		return "", fmt.Errorf("Could not find environment variable: %s", key)
	}

	return val, nil
}
