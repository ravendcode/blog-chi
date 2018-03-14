package config

import (
	"flag"
	"os"
)

// Config app
type Config struct {
	Env  string
	Port string
}

// NewConfig create new Config
func NewConfig() *Config {
	env := flag.String("env", os.Getenv("ENV"), "server environment")
	port := flag.String("port", os.Getenv("PORT"), "server listening port")

	flag.Parse()

	config := &Config{
		Env:  *env,
		Port: *port,
	}

	return config
}
