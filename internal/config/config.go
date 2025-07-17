package config

import (
	"os"

	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

type Config struct {
	DB DBConfig
}

type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
}

func Load(log *zap.Logger, env string) *Config {
	if err := godotenv.Load(env); err != nil {
		log.Warn("No .env file found, using system env")
	}

	return &Config{
		DB: DBConfig{
			Host:     getEnv("DB_HOST", log),
			Port:     getEnv("DB_PORT", log),
			User:     getEnv("DB_USER", log),
			Password: getEnv("DB_PASSWORD", log),
			Name:     getEnv("DB_NAME", log),
		},
	}
}

func getEnv(key string, log *zap.Logger) string {
	if val, exists := os.LookupEnv(key); exists {
		return val
	}
	log.Error("Обязательная переменная окружения не установлена", zap.String("key", key))
	panic("missing required environment variable: " + key)
}
