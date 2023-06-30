package config

import (
	"errors"
	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Db struct {
		Mongo struct {
			Host       string `env:"HOST" env-default:"127.0.0.1"`
			Port       int    `env:"PORT" env-default:"27017"`
			Collection string `env:"COLLECTION" env-default:"products"`
			Database   string `env:"DATABASE" env-default:"cost_control"`
		}
	}
}

func GetConfig() (*Config, error) {
	var cfg Config

	err := cleanenv.ReadEnv(&cfg)
	if err != nil {
		return nil, errors.New("no config found")
	}

	return &cfg, nil
}
