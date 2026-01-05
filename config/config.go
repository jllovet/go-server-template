package config

import (
	"os"

	"github.com/joho/godotenv"
)

var InitializedConfig = InitConfig()

type Config struct {
	Host        string
	Port        string
	DatabaseURL string
}

func InitConfig() Config {
	godotenv.Load()
	return Config{
		Host:        GetEnv("PROJECT_HOST", "localhost"),
		Port:        GetEnv("PROJECT_PORT", "8080"),
		DatabaseURL: GetEnv("DATABASE_URL", ""),
	}
}

func GetEnv(key, defaultValue string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return defaultValue
}
