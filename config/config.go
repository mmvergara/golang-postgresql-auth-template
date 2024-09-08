package config

import (
	"fmt"
	"log"
	"os"
	"strconv"

	_ "github.com/joho/godotenv/autoload"
)

type Config struct {
	APPENV                   string
	PORT                     int
	JWT_TOKEN_DURATION_HOURS int
	JWT_SECRET               string
	DB_DATABASE              string
	DB_PASSWORD              string
	DB_USERNAME              string
	DB_PORT                  string
	DB_HOST                  string
	DB_SCHEMA                string
}

func NewConfig() *Config {
	requiredEnvVars := []string{
		"DB_DATABASE", "DB_PASSWORD", "DB_USERNAME", "DB_PORT", "DB_HOST", "DB_SCHEMA",
		"JWT_TOKEN_DURATION_HOURS", "JWT_SECRET", "PORT",
	}

	for _, envVar := range requiredEnvVars {
		if os.Getenv(envVar) == "" {
			log.Fatalf("%s required environment variable is not set", envVar)
		}
	}

	app_env := os.Getenv("APP_ENV")
	if app_env == "" {
		app_env = "development"
	}
	fmt.Println("Running in", app_env, "mode")

	jwt_token_duration_hours, err := strconv.Atoi(os.Getenv("JWT_TOKEN_DURATION_HOURS"))
	if err != nil {
		log.Fatal("JWT_TOKEN_DURATION_HOURS environment variable is not a valid integer")
	}

	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		log.Fatal("PORT environment variable is not a valid integer")
	}

	jwt_secret := os.Getenv("JWT_SECRET")

	db_database := os.Getenv("DB_DATABASE")
	db_password := os.Getenv("DB_PASSWORD")
	db_username := os.Getenv("DB_USERNAME")
	db_port := os.Getenv("DB_PORT")
	db_host := os.Getenv("DB_HOST")
	db_schema := os.Getenv("DB_SCHEMA")

	return &Config{
		APPENV:                   app_env,
		DB_DATABASE:              db_database,
		DB_PASSWORD:              db_password,
		DB_USERNAME:              db_username,
		DB_PORT:                  db_port,
		DB_HOST:                  db_host,
		DB_SCHEMA:                db_schema,
		JWT_TOKEN_DURATION_HOURS: jwt_token_duration_hours,
		JWT_SECRET:               jwt_secret,
		PORT:                     port,
	}
}
