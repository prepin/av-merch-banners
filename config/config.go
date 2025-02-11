package config

import (
	"av-merch-shop/pkg/logging"
	"log/slog"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Logger *slog.Logger
	Server ServerConfig
	Auth   AuthConfig
	DB     DBConfig
}

func Load() *Config {
	config := &Config{
		Logger: logging.NewLogger(),
	}
	config.load()
	return config
}

func (c *Config) load() {
	err := godotenv.Load()
	if err != nil {
		c.Logger.Warn("Warning: .env file not found or error loading:", "error", err)
	}

	c.loadDBConfig()
	c.loadAuthConfig()
	c.loadServerConfig()
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
