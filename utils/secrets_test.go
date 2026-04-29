package utils

import (
	"os"
	"path/filepath"
	"testing"
)

func TestReadSecrets(t *testing.T) {
	tempDir := t.TempDir()
	path := filepath.Join(tempDir, "rabbitmq_pass")

	if err := os.WriteFile(path, []byte("user=111\n"), 0600); err != nil {
		t.Fatal(err)
	}

	content := ReadSecret(path)
	if content != "user=111" {
		t.Errorf("content not correct. result: %s, expected: %s\n", content, "user111")
	}
}

func TestGetSecretString(t *testing.T) {
	val := GetSecretString("AMQP_PASS=1234", "AMQP_PASS", "opps")
	if val != "1234" {
		t.Errorf("value not correct. result %s, expected :%s\n", val, "1234")
	}
}
