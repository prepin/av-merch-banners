package config

import (
	"strconv"
)

type ServerConfig struct {
	Port           string
	ReadTimeout    int
	WriteTimeout   int
	RequestTimeout int
}

func (c *Config) loadServerConfig() {
	readTimeout, err := strconv.Atoi(getEnv("AV_SERVER_READ_TIMEOUT", "5"))
	if err != nil {
		c.Logger.Error("Error: AV_SERVER_READ_TIMEOUT must be an integer")
	}

	writeTimeout, err := strconv.Atoi(getEnv("AV_SERVER_WRITE_TIMEOUT", "5"))
	if err != nil {
		c.Logger.Error("Error: AV_SERVER_WRITE_TIMEOUT must be an integer")
	}

	requestTimeout, err := strconv.Atoi(getEnv("AV_REQUEST_TIMEOUT", "50"))
	if err != nil {
		c.Logger.Error("Error: AV_REQUEST_TIMEOUT must be an integer")
	}

	c.Server = ServerConfig{
		Port:           getEnv("AV_SERVER_PORT", ":8080"),
		ReadTimeout:    readTimeout,
		WriteTimeout:   writeTimeout,
		RequestTimeout: requestTimeout,
	}
}
