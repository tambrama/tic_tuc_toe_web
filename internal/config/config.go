package config

import (
	"os"
)

type Config struct {
	ServerPort string
	DB         ConfigDB
	JWT        []byte
}

type ConfigDB struct {
	URL string
}

func NewConfig() *Config {
	return &Config{
		ServerPort: getEnv("SERVER_PORT", "8081"),
		DB: ConfigDB{
			URL: getEnv("DATABASE_URL", ""),
		},
		JWT: []byte(getEnv("JWT_SECRET", "")),
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
