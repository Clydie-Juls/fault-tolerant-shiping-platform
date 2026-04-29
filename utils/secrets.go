package utils

import (
	"log"
	"os"
	"strings"
)

func ReadSecret(file string) string {
	data, err := os.ReadFile(file)
	if err != nil {
		log.Println("unable to read or find secret file")
		return ""
	}

	log.Println(string(data))
	return strings.TrimSpace(string(data))
}

func GetSecretString(content string, key string, fallback string) string {
	values := strings.Split(content, "=")
	for i := 0; i < len(values); i += 2 {
		if strings.EqualFold(values[i], key) {
			return values[i+1]
		}
	}

	return fallback
}
