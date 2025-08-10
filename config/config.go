package config

import (
	"os"

	"github.com/joho/godotenv"
)

var InitializedConfig = InitConfig()

type Config struct {
	Host string
	Port string
}

func InitConfig() Config {
	godotenv.Load()
	return Config{
		Host: GetEnv("PROJECT_HOST", "localhost"),
		Port: GetEnv("PROJECT_PORT", "8080"),
	}
}

func GetEnv(key, defaultValue string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return defaultValue
}
