package config

import (
	"github.com/joho/godotenv"
	"github.com/spf13/cast"

	"log"
	"os"
)

type Config struct {
	USER_PORT string
	POST_PORT string

	DB_PORT     string
	DB_HOST     string
	DB_USER     string
	DB_PASSWORD string
	DB_NAME     string
}

func Load() Config {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("Error loading .env file")
	}

	config := Config{}

	config.USER_PORT = cast.ToString(coalesce("USER_PORT", ":50050"))
	config.POST_PORT = cast.ToString(coalesce("POST_PORT", ":7070"))
	config.DB_PORT = cast.ToString(coalesce("DB_PORT", "5432"))
	config.DB_HOST = cast.ToString(coalesce("DB_HOST", "localhost"))
	config.DB_USER = cast.ToString(coalesce("DB_USER", "postgres"))
	config.DB_PASSWORD = cast.ToString(coalesce("DB_PASSWORD", "dodi"))
	config.DB_NAME = cast.ToString(coalesce("DB_NAME", "cm"))

	return config
}

func coalesce(key string, defaultValue interface{}) interface{} {
	val, exists := os.LookupEnv(key)

	if exists {
		return val
	}

	return defaultValue
}
