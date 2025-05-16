package config

import (
	"github.com/tharun-rs/rprox/logger"

	"github.com/caarlos0/env/v6"
)

type Config struct {
	AppPort  string `env:"PORT" envDefault:":8080"`
	RedisURL string `env:"REDIS_URL" envDefault:"redis://localhost:6379"`
	RedisPass string `env:"REDIS_PASS" envDefault:""`
	RedisDB   int    `env:"REDIS_DB" envDefault:"0"`
}

var Cfg Config

func Init() {
	if err := env.Parse(&Cfg); err != nil {
		logger.Log.Errorf("Failed to parse env: %v", err)
		return
	}
	logger.Log.Infof("Loaded config — PORT: %s, REDIS_URL: %s", Cfg.AppPort, Cfg.RedisURL)
}
