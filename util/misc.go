package util

import "os"

// Environ ...
func Environ(name string, fallback string) string {
	result := os.Getenv(name)

	if result == "" {
		result = fallback
	}

	return result
}
