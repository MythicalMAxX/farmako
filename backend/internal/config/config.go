package config

import (
	"encoding/json"
	"os"
)

// Config represents the application configuration
type Config struct {
	Server   ServerConfig   `json:"server"`
	Database DatabaseConfig `json:"database"`
}

// ServerConfig represents server configuration
type ServerConfig struct {
	Port string `json:"port"`
}

// DatabaseConfig represents database configuration
type DatabaseConfig struct {
	Driver   string `json:"driver"`
	Host     string `json:"host"`
	Port     string `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Name     string `json:"name"`
	SSLMode  string `json:"sslmode"`
}

// Load loads configuration from a file and environment variables
func Load(filename string) (*Config, error) {
	// Get default config
	config := DefaultConfig()

	// Try to load from file
	file, err := os.Open(filename)
	if err == nil {
		defer file.Close()
		decoder := json.NewDecoder(file)
		if err := decoder.Decode(config); err != nil {
			return nil, err
		}
	}

	// Override with environment variables if present
	if port := os.Getenv("SERVER_PORT"); port != "" {
		config.Server.Port = port
	}

	if driver := os.Getenv("DB_DRIVER"); driver != "" {
		config.Database.Driver = driver
	}

	if host := os.Getenv("DB_HOST"); host != "" {
		config.Database.Host = host
	}

	if port := os.Getenv("DB_PORT"); port != "" {
		config.Database.Port = port
	}

	if user := os.Getenv("DB_USER"); user != "" {
		config.Database.User = user
	}

	if password := os.Getenv("DB_PASSWORD"); password != "" {
		config.Database.Password = password
	}

	if name := os.Getenv("DB_NAME"); name != "" {
		config.Database.Name = name
	}

	if sslMode := os.Getenv("DB_SSLMODE"); sslMode != "" {
		config.Database.SSLMode = sslMode
	}

	return config, nil
}

// DefaultConfig returns default configuration
func DefaultConfig() *Config {
	return &Config{
		Server: ServerConfig{
			Port: "8081",
		},
		Database: DatabaseConfig{
			Driver:   "sqlite",
			Host:     "",
			Port:     "",
			User:     "",
			Password: "",
			Name:     "coupon_system.db",
			SSLMode:  "",
		},
	}
}
