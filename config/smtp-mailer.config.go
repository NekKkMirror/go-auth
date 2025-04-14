package config

import (
	"log"
	"os"
)

// SMTPConfig contains SMTP settings.
type SMTPConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	From     string
}

// LoadSMTPConfig loads SMTP data from environment variables.
func LoadSMTPConfig() *SMTPConfig {
	cfg := &SMTPConfig{
		Host:     os.Getenv("SMTP_HOST"),
		Port:     os.Getenv("SMTP_PORT"),
		Username: os.Getenv("SMTP_USER"),
		Password: os.Getenv("SMTP_PASS"),
		From:     os.Getenv("SMTP_FROM"),
	}

	if cfg.Host == "" || cfg.Port == "" || cfg.Username == "" || cfg.Password == "" || cfg.From == "" {
		log.Fatal("SMTP configuration is invalid or missing required environment variables")
	}

	return cfg
}
