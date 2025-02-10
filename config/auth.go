package config

type AuthConfig struct {
	SecretKey []byte
}

func (c *Config) loadAuthConfig() {

	c.Auth = AuthConfig{
		SecretKey: []byte(getEnv("AV_SECRET", "default")),
	}
}
