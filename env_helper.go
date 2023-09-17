package main

import (
	"os"
	"strconv"
)

func getEnvStr(key string, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	if value, ok := os.LookupEnv(key); ok {
		ans, err := strconv.Atoi(value)
		if err != nil {
			return 0
		} else {
			return ans
		}
	}
	return fallback
}
