package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseUrl              string
	DatabaseName             string
	DatabaseUsername         string
	DatabasePassword         string
	ServerMode               string
	ServerPort               string
	AccessTokenKey           string
	RefreshTokenKey          string
	EncryptionKey            string
	AccessTokenExpiredMinute string
	ApiKey                   string
	AllowOrigin              string
}

func LoadConfig() (*Config, error) {
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "dev"
	}

	envFile := fmt.Sprintf(".env.%s", env)

	err := godotenv.Load(envFile)
	if err != nil {
		log.Fatalf("Error loading %s file: %v", envFile, err)
	}

	config := &Config{
		DatabaseUrl:              os.Getenv("DATABASE_URL"),
		DatabaseName:             os.Getenv("DATABASE_DB_NAME"),
		DatabaseUsername:         os.Getenv("DATABASE_USERNAME"),
		DatabasePassword:         os.Getenv("DATABASE_PASSWORD"),
		ServerMode:               os.Getenv("SERVER_MODE"),
		ServerPort:               os.Getenv("SERVER_PORT"),
		AccessTokenKey:           os.Getenv("ACCESS_TOKEN_KEY"),
		RefreshTokenKey:          os.Getenv("REFRESH_TOKEN_32_BYTES_KEY"),
		EncryptionKey:            os.Getenv("ENCRYPTION_DATA_32_BYTES_KEY"),
		AccessTokenExpiredMinute: os.Getenv("ACCESS_TOKEN_EXPIRED_MINUTE"),
		ApiKey:                   os.Getenv("API_KEY"),
		AllowOrigin:              os.Getenv("ALLOW_ORIGIN"),
	}

	return config, nil
}
