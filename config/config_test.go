package config

import (
	"bufio"
	"os"
	"testing"
)

func createEnvFile(envVars map[string]string) error {
	file, err := os.Create(".env.testonly")
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for key, value := range envVars {
		os.Setenv(key, value)
		_, err := writer.WriteString(key + "=" + value + "\n")
		if err != nil {
			return err
		}
	}

	writer.Flush()

	return nil
}

func TestNewAppConfig(t *testing.T) {
	// Call func to generate .env file
	err := createEnvFile(map[string]string{
		"PG_USERNAME": "testuser",
		"PG_PASS":     "testpassword",
		"PG_DB":       "testdb",
		"PG_PORT":     "5432",
		"PG_HOST":     "localhost",
		"REDIS_HOST":  "localhost",
		"REDIS_PORT":  "6379",
		"REDIS_USN":   "testuser",
		"REDIS_PASS":  "testpassword",
		"KEY":         "testkey",
		"PORT":        "8080",
	})
	if err != nil {
		t.Error("Failed to create .env file:", err)
		return
	}

	// Test reading from environment variables
	config, err := NewAppConfig(".env.testonly")
	if err != nil {
		t.Errorf("Error creating config: %v", err)
		return
	}
	if config.DbUsername != "testuser" {
		t.Errorf("Expected DbUsername to be 'testuser', got '%s'", config.DbUsername)
	}
	if config.DbPassword != "testpassword" {
		t.Errorf("Expected DbPassword to be 'testpassword', got '%s'", config.DbPassword)
	}
	if config.DbName != "testdb" {
		t.Errorf("Expected DbName to be 'testdb', got '%s'", config.DbName)
	}
	if config.DbPort != 5432 {
		t.Errorf("Expected DbPort to be 5432, got '%d'", config.DbPort)
	}
	if config.DbHost != "localhost" {
		t.Errorf("Expected DbHost to be 'localhost', got '%s'", config.DbHost)
	}
	if config.RedisHost != "localhost" {
		t.Errorf("Expected RedisHost to be 'localhost', got '%s'", config.RedisHost)
	}
	if config.RedisPort != "6379" {
		t.Errorf("Expected RedisPort to be '6379', got '%s'", config.RedisPort)
	}
	if config.RedisUsn != "testuser" {
		t.Errorf("Expected RedisUsn to be 'testuser', got '%s'", config.RedisUsn)
	}
	if config.RedisPass != "testpassword" {
		t.Errorf("Expected RedisPass to be 'testpassword', got '%s'", config.RedisPass)
	}
	if config.EncryptKey != "testkey" {
		t.Errorf("Expected EncryptKey to be 'testkey', got '%s'", config.EncryptKey)
	}
	if config.Port != 8080 {
		t.Errorf("Expected Port to be 8080, got '%d'", config.Port)
	}

	if _, err := os.Stat(".env.testonly"); err == nil {
		defer os.Remove(".env.testonly")
	}
}
