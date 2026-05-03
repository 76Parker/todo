package config

import (
	"github.com/caarlos0/env/v11"
)

type Validator interface {
	Validate(cfg *Config) error
}

type Config struct {
	App      App      `envPrefix:"TODO"`
	Postgres Postgres `envPrefix:"POSTGRES"`
}

type App struct {
	Host string `env:"HOST" envDefault:"localhost" validate:"required"`
	Port int    `env:"PORT" envDefault:"8080" validate:"required,allowed_ports"`
}

type Postgres struct {
	Host string `env:"HOST" envDefault:"localhost" validate:"required"`
	Port int    `env:"PORT" envDefault:"5432" validate:"required,allowed_ports"`
}

func Load() (*Config, error) {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, err
	}
	if err := validateConfig(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}

func validateConfig(cfg *Config) error {
	v, err := NewValidator()
	if err != nil {
		return err
	}
	if err := v.Validate(cfg); err != nil {
		return err
	}
	return nil
}
