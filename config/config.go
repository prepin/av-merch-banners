package config

import (
	"av-merch-shop/pkg/logging"
	"log/slog"
)

type Config struct {
	Logger *slog.Logger
}

func Load() *Config {
	config := &Config{
		Logger: logging.NewLogger(),
	}
	return config
}
