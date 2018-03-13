package config

import (
	"os"
)

// Config app
type Config struct {
	Env  string
	Port string
}

// NewConfig create new Config
func NewConfig() *Config {
	config := &Config{
		Env:  os.Getenv("ENV"),
		Port: os.Getenv("PORT"),
	}

	if config.Env == "" {
		config.Env = "production"
	}
	if config.Env == "production" {
		config.Port = "80"
	}
	return config
}
