package config

// Config app
type Config struct {
	Env  string
	Port string
}

// NewConfig create new Config
func NewConfig() *Config {
	config := &Config{
		// Env:  os.Getenv("ENV"),
		// Port: os.Getenv("PORT"),
		Env:  "production",
		Port: "80",
	}

	// if config.Env == "" {
	// 	config.Env = "production"
	// }
	// if config.Env == "production" {
	// 	config.Port = "80"
	// }
	return config
}
