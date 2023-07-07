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
	SignedKey        string `env:"SIGNED_KEY" env-required:"true"`
	ExpiredAtMinutes uint   `env:"TOKEN_EXPIRED_AT_MINUTES" env-default:"60"`
	Rest             struct {
		Port int `env:"REST_SERVER_PORT" env-default:"10000"`
	}
	Rpc struct {
		Address string `env:"RPC_ADDRESS" env-default:":5300"`
	}
	TelegramBotToken string `env:"TELEGRAM_BOT_TOKEN" env-required:"true"`
}

func GetConfig() (*Config, error) {
	var cfg Config
	err := cleanenv.ReadEnv(&cfg)
	if err != nil {
		return nil, errors.New("no config found")
	}

	return &cfg, nil
}
