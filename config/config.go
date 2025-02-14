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
		c.Logger.Debug(".env file not found, only ENV variables are used")
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
