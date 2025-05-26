package main

import (
	"fmt"
	"github.com/caarlos0/env/v11"
)

type Config struct {
	Server HTTPServer `envPrefix:"HTTP_SERVER_"`
	DB     DB         `envPrefix:"DB_"`
}

type DB struct {
	User     string `env:"USER,notEmpty"`
	Password string `env:"PASS,notEmpty"`
	Name     string `env:"NAME,notEmpty"`
	Host     string `env:"HOST,notEmpty"`
	Port     int    `env:"PORT,notEmpty"`
	SSLMode  string `env:"SSL_MODE,notEmpty"`
}

type HTTPServer struct {
	Port string `env:"PORT,notEmpty"`
}

func loadConfigFromEnv() (Config, error) {
	cfg, err := env.ParseAs[Config]()
	if err != nil {
		return Config{}, fmt.Errorf("failed to parse config: %w", err)
	}

	return cfg, nil
}
