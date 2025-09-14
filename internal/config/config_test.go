package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetEnv(t *testing.T) {
	t.Run("environment_variable_exists", func(t *testing.T) {
		key := "TEST_ENV_VAR"
		expectedValue := "test_value"
		os.Setenv(key, expectedValue)
		defer os.Unsetenv(key) // Clean up after test

		result := getEnv(key, "fallback_value")

		assert.Equal(t, expectedValue, result, "Should return the environment variable value")
	})

	t.Run("environment_variable_does_not_exist", func(t *testing.T) {
		key := "NONEXISTENT_ENV_VAR"
		os.Unsetenv(key)

		fallbackValue := "fallback_value"

		result := getEnv(key, fallbackValue)

		assert.Equal(t, fallbackValue, result, "Should return the fallback value")
	})

	t.Run("environment_variable_empty", func(t *testing.T) {
		key := "EMPTY_ENV_VAR"
		os.Setenv(key, "")
		defer os.Unsetenv(key) // Clean up after test

		fallbackValue := "fallback_value"

		result := getEnv(key, fallbackValue)

		assert.Equal(t, "", result, "Should return the empty environment variable value")
	})
}

func TestLoadConfig(t *testing.T) {
	t.Run("default_values", func(t *testing.T) {
		os.Unsetenv("MONGODB_DATABASE_URL")
		os.Unsetenv("MONGODB_DATABASE")

		config := LoadConfig()

		assert.Equal(t, "mongodb://localhost:27017", config.MongoURI, "Should use default MongoDB URI")
		assert.Equal(t, "testdb", config.MongoDatabase, "Should use default MongoDB database name")
	})

	t.Run("custom_values_from_env", func(t *testing.T) {
		customURI := "mongodb://customhost:27017"
		customDB := "customdb"
		os.Setenv("MONGODB_DATABASE_URL", customURI)
		os.Setenv("MONGODB_DATABASE", customDB)
		defer func() {
			os.Unsetenv("MONGODB_DATABASE_URL")
			os.Unsetenv("MONGODB_DATABASE")
		}()

		config := LoadConfig()

		assert.Equal(t, customURI, config.MongoURI, "Should use custom MongoDB URI from environment")
		assert.Equal(t, customDB, config.MongoDatabase, "Should use custom MongoDB database name from environment")
	})

	t.Run("mixed_values", func(t *testing.T) {
		customURI := "mongodb://mixedhost:27017"
		os.Setenv("MONGODB_DATABASE_URL", customURI)
		os.Unsetenv("MONGODB_DATABASE")
		defer os.Unsetenv("MONGODB_DATABASE_URL")

		config := LoadConfig()

		assert.Equal(t, customURI, config.MongoURI, "Should use custom MongoDB URI from environment")
		assert.Equal(t, "testdb", config.MongoDatabase, "Should use default MongoDB database name")
	})

	t.Run("empty_values_in_env", func(t *testing.T) {
		os.Setenv("MONGODB_DATABASE_URL", "")
		os.Setenv("MONGODB_DATABASE", "")
		defer func() {
			os.Unsetenv("MONGODB_DATABASE_URL")
			os.Unsetenv("MONGODB_DATABASE")
		}()

		config := LoadConfig()

		assert.Equal(t, "", config.MongoURI, "Should use empty MongoDB URI from environment")
		assert.Equal(t, "", config.MongoDatabase, "Should use empty MongoDB database name from environment")
	})
}

func TestLoadConfigWithEnvFile(t *testing.T) {
	envContent := `MONGODB_DATABASE_URL=mongodb://envfile:27017
MONGODB_DATABASE=envfiledb
`
	err := os.WriteFile(".env", []byte(envContent), 0o644)
	if err != nil {
		t.Fatalf("Failed to create temporary .env file: %v", err)
	}
	defer os.Remove(".env") // Clean up after test

	os.Unsetenv("MONGODB_DATABASE_URL")
	os.Unsetenv("MONGODB_DATABASE")

	config := LoadConfig()

	assert.Equal(t, "mongodb://envfile:27017", config.MongoURI, "Should use MongoDB URI from .env file")
	assert.Equal(t, "envfiledb", config.MongoDatabase, "Should use MongoDB database name from .env file")
}

func TestEnvironmentVariablesPrecedence(t *testing.T) {
	envContent := `MONGODB_DATABASE_URL=mongodb://envfile:27017
MONGODB_DATABASE=envfiledb
`
	err := os.WriteFile(".env", []byte(envContent), 0o644)
	if err != nil {
		t.Fatalf("Failed to create temporary .env file: %v", err)
	}
	defer os.Remove(".env") // Clean up after test

	envURI := "mongodb://envvar:27017"
	envDB := "envvardb"
	os.Setenv("MONGODB_DATABASE_URL", envURI)
	os.Setenv("MONGODB_DATABASE", envDB)
	defer func() {
		os.Unsetenv("MONGODB_DATABASE_URL")
		os.Unsetenv("MONGODB_DATABASE")
	}()

	config := LoadConfig()

	assert.Equal(t, envURI, config.MongoURI, "Environment variable should take precedence over .env file")
	assert.Equal(t, envDB, config.MongoDatabase, "Environment variable should take precedence over .env file")
}
