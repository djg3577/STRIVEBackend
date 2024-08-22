package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	DBUser        string
	DBPassword    string
	DBName        string
	DBHost        string
	DBPort        string
	RedisHost     string
	RedisPassword string
	RedisDB       int
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
		DBUser:        os.Getenv("DB_USER"),
		DBPassword:    os.Getenv("DB_PASSWORD"),
		DBName:        os.Getenv("DB_NAME"),
		DBHost:        os.Getenv("DB_HOST"),
		DBPort:        os.Getenv("DB_PORT"),
		RedisHost:     os.Getenv("REDIS_HOST"),
		RedisPassword: os.Getenv("REDIS_PASSWORD"),
		RedisDB:       convertStringToInt(os.Getenv("REDIS_DB")),
	}
}
func convertStringToInt(StringToConvert string) int {
	intValue, err := strconv.Atoi(StringToConvert)
	if err != nil {
		log.Printf("Warning: Environment variable %s is not a valid integer", StringToConvert)
	}
	return intValue
}

func LoadConfig(loader ConfigLoader) *Config {
	return loader.Load()
}
