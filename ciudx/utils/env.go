package utils

import (
	"os"
)

// Getenv returns env var is set else the default value
func Getenv(key, defaultVal string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return defaultVal
}
