package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config ...
type Config struct {
	AppPort string

	DbHost     string
	DbPort     string
	DbUser     string
	DbPassword string
	DbName     string
	DbLogMode  bool

	JwtSigningKey string
}

// SetConfiguration reads config from env var
// and return a config struct
func SetConfiguration() Config {
	if err := godotenv.Load(); err != nil {
		log.Print("could not load .env file, read from env")
	}

	dbLogMode := false
	dbLogMode, _ = strconv.ParseBool(os.Getenv("DB_LOG_MODE"))

	return Config{
		AppPort: os.Getenv("APP_PORT"),

		DbHost:     os.Getenv("DB_HOST"),
		DbPort:     os.Getenv("DB_PORT"),
		DbUser:     os.Getenv("DB_USER"),
		DbPassword: os.Getenv("DB_PASSWORD"),
		DbName:     os.Getenv("DB_NAME"),
		DbLogMode:  dbLogMode,

		JwtSigningKey: os.Getenv("JWT_SIGNING_KEY"),
	}
}
