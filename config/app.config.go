package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config represents the main configuration structure for the application.
type Config struct {
	AppPort        string `json:"app_port"`
	AppEnv         string `json:"app_env"`
	AppAPIBasePath string `json:"app_api_base_path"`
	AppJWTSecret   string `json:"app_jst_secret"`

	DBUrl string `json:"db_url"`
}

// LoadConfig loads the configuration from the .env file or environment variables. It returns a pointer to the Config struct. If the .env file is not found, it logs a warning and uses environment variables instead. If any error occurs during the loading process, it logs an error and returns nil.
func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: No .env file found. Using environment variables.")
	}

	return &Config{
		AppPort:        os.Getenv("APP_PORT"),
		AppEnv:         os.Getenv("APP_ENV"),
		AppAPIBasePath: os.Getenv("APP_API_BASE_PATH"),
		AppJWTSecret:   os.Getenv("APP_JWT_SECRET"),

		DBUrl: os.Getenv("DB_URL"),
	}
}
