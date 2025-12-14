package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port               string
	DatabaseURL        string
	KafkaBrokers       []string
	KafkaTopic         string
	MatchmakingTimeout int
}

func Load() *Config {
	godotenv.Load()

	port := getEnv("PORT", "8080")
	// Ensure port has a colon prefix
	if port[0] != ':' {
		port = ":" + port
	}

	return &Config{
		Port:               port,
		DatabaseURL:        getEnv("DATABASE_URL", "postgres://user:pass@localhost/4inrow"),
		KafkaBrokers:       []string{getEnv("KAFKA_BROKERS", "localhost:9092")},
		KafkaTopic:         getEnv("KAFKA_TOPIC", "game-events"),
		MatchmakingTimeout: 10, // seconds
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
