package config

import (
	"log"
	"os"
	"strconv"
)

type Config struct {
	Database DatabaseConfig
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
}

func Load() *Config {
	return &Config{
		Database: DatabaseConfig{
			Host:     mustEnv("POSTGRES_HOST"),
			Port:     mustEnv("POSTGRES_PORT"),
			User:     mustEnv("POSTGRES_USER"),
			Password: mustEnv("POSTGRES_PASSWORD"),
			Name:     mustEnv("POSTGRES_DB"),
		},
	}
}

func mustEnv(key string) string {
	val := os.Getenv(key)
	if val == "" {
		log.Fatalf("ENV variable %s is required but not set", key)
	}
	return val
}

func mustEnvInt(key string) int {
	valStr := mustEnv(key)
	val, err := strconv.Atoi(valStr)
	if err != nil {
		log.Fatalf("ENV variable %s must be an integer, got: %s", key, valStr)
	}
	return val
}
