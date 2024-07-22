package config

import (
	"os"
)

type Config struct {
	MongoURI    string
	KafkaBroker string
	JWTSecret   string
}

func LoadConfig() *Config {
	return &Config{
		MongoURI:    getEnv("MONGO_URI", "mongodb://localhost:27017"),
		KafkaBroker: getEnv("KAFKA_BROKER", "localhost:9092"),
		JWTSecret:   getEnv("JWT_SECRET", "your_secret_key"),
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
