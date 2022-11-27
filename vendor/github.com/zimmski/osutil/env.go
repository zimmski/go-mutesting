package osutil

import (
	"os"
)

// EnvOrDefault returns the environment variable with the given key, or the default value if the key is not defined.
func EnvOrDefault(key string, defaultValue string) string {
	if v, ok := os.LookupEnv(key); ok {
		return v
	}

	return defaultValue
}
