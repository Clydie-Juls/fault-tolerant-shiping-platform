package utils

import (
	"os"
	"strconv"
	"strings"
)

func GetEnvInt(key string, fallback int) int {
	val := GetEnvString(key, "")
	intVal, err := strconv.Atoi(val)
	if err != nil {
		return fallback
	}

	return intVal
}

func GetEnvFloat(key string, fallback float64) float64 {
	val := GetEnvString(key, "")
	floatVal, err := strconv.ParseFloat(strings.TrimSpace(val), 64)
	if err != nil {
		return fallback
	}

	return floatVal
}

func GetEnvString(key string, fallback string) string {
	val, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}

	return val
}
