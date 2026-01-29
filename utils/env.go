package utils

import (
	"os"
	"strconv"
	"strings"
)

func GetEnvInt(key string, fallback int) int {
	val := GetEnvString(key, "")
	intVal, err := strconv.Atoi(val)
	FailOnError(err, "unable to parse string into int")

	return intVal
}

func GetEnvFloat(key string, fallback float64) float64 {
	val := GetEnvString(key, "")
	floatVal, err := strconv.ParseFloat(strings.TrimSpace(val), 64)
	FailOnError(err, "unable to parse string into int")

	return floatVal
}

func GetEnvString(key string, fallback string) string {
	val, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}

	return val
}
