package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBUser     string
	DBPassword string
	DBName     string
	DBHost     string
	DBPort     string
}

type ConfigLoader interface {
	Load() *Config
}

type EnvConfigLoader struct{}

func (e *EnvConfigLoader) Load() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Printf("Error loading .env file: %v", err)
	}

	return &Config{
		DBUser:     os.Getenv("DB_USER"),
		DBPassword: os.Getenv("DB_PASSWORD"),
		DBName:     os.Getenv("DB_NAME"),
		DBHost:     os.Getenv("DB_HOST"),
		DBPort:     os.Getenv("DB_PORT"),
	}
}

func LoadConfig(loader ConfigLoader) *Config {
	return loader.Load()
}
