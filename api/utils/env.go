package utils

import (
	"os"
)

func GetEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func VerifyEnv(keys []string) bool {
	for _, key := range keys {
		if _, ok := os.LookupEnv(key); !ok {
			return false
		}
	}
	return true
}
