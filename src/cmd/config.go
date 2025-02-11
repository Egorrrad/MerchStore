package cmd

import "github.com/caarlos0/env/v9"

type Config struct {
	ServerPort string `env:"SERVER_PORT,required"`
	Database   DatabaseConfig
}

type DatabaseConfig struct {
	User     string `env:"DATABASE_USER,required"`
	Password string `env:"DATABASE_PASSWORD,required"`
	Host     string `env:"DATABASE_HOST,required"`
	Port     string `env:"DATABASE_PORT,required"`
	Name     string `env:"DATABASE_NAME,required"`
}

func Load() (*Config, error) {
	cfg := Config{}

	err := env.Parse(&cfg)
	cfg.ServerPort = ":" + cfg.ServerPort
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}
