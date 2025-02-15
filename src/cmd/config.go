package cmd

import (
	"github.com/caarlos0/env/v9"
)

type Config struct {
	Service  ServiceConfig
	Database DatabaseConfig
	Cache    CacheConfig
}

type DatabaseConfig struct {
	User     string `env:"DATABASE_USER,required"`
	Password string `env:"DATABASE_PASSWORD,required"`
	Host     string `env:"DATABASE_HOST,required"`
	Port     string `env:"DATABASE_PORT,required"`
	Name     string `env:"DATABASE_NAME,required"`
}

type CacheConfig struct {
	Host string `env:"CACHE_HOST,required"`
	Port string `env:"CACHE_PORT,required"`
}

type ServiceConfig struct {
	ServerPort string `env:"SERVER_PORT,required"`
	LogLevel   string `env:"LOG_LEVEL,required"`
	SecretKey  string `env:"SECRET_KEY,required"`
}

func Load() (*Config, error) {
	cfg := Config{}

	err := env.Parse(&cfg)
	cfg.Service.ServerPort = ":" + cfg.Service.ServerPort
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}
