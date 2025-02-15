package config

import "strconv"

type RedisConfig struct {
	Addr     string
	Password string
	DB       int
}

func (c *Config) loadRedisConfig() {

	db, err := strconv.Atoi(getEnv("AV_REDIS_DB", "0"))
	if err != nil {
		c.Logger.Error("Error: AV_REDIS_DB must be an integer")
	}

	rc := RedisConfig{
		DB:       db,
		Addr:     getEnv("AV_REDIS_ADDR", "localhost:6379"),
		Password: getEnv("AV_REDIS_PASSWORD", "password"),
	}
	c.Redis = rc
}
