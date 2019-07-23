package config

import (
	"github.com/caarlos0/env"
	"github.com/pkg/errors"
)

// Global configuration struct.
type Config struct {
	Log LogConf
}

// LogConf - logging configuration struct.
type LogConf struct {
	Level             string `env:"LOG_LEVEL" envDefault:"debug"`
	DevelopmentLogger bool   `env:"DEVELOPMENT_LOGGER" envDefault:"FALSE"`
}

// Gets a configuration struct from an environment.
func GetConfigFromEnv() (Config, error) {
	config := Config{}

	if err := env.Parse(&config.Log); err != nil {
		return Config{}, errors.Wrap(err, "failed to load Log config")
	}

	return config, nil
}
