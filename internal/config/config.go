package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	MongoURI      string
	MongoDatabase string
}

func LoadConfig() *Config {
	err := godotenv.Load()

	if err != nil {
		log.Println("⚠️  Nenhum arquivo .env encontrado, usando variáveis do sistema.")
	}

	config := &Config{
		MongoURI:      getEnv("MONGODB_DATABASE_URL", "mongodb://localhost:27017"),
		MongoDatabase: getEnv("MONGODB_DATABASE", "testdb"),
	}

	return config
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
