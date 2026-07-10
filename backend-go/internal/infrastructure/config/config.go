package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	httpPort string
	jwtSecret string
}

func Load() *Config {
	godotenv.Load()

	httpPort := os.Getenv("HTTP_PORT")
	if httpPort == "" {
		httpPort = "8080"
	}

	return &Config{
		httpPort:  httpPort,
		jwtSecret: os.Getenv("JWT_SECRET"),
	}
}

func (c *Config) GetJWTSecret() string {
	return c.jwtSecret
}

func (c *Config) GetHTTPPort() string {
	return c.httpPort
}
