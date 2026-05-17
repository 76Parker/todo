// Package config provide configurable params from config.yaml
package config

import (
	"os"

	"github.com/76Parker/golib/httplib"
	"github.com/76Parker/golib/loglib"
	"github.com/76Parker/golib/pglib"
	"gopkg.in/yaml.v3"
)

// Validator validate config
type Validator interface {
	Validate(cfg *Config) error
}

// Config provide all configurable params
type Config struct {
	Logger   loglib.SlogConfig `yaml:"logger" validate:"required"`
	Postgres pglib.Config      `yaml:"postgres" validate:"required"`
	HTTP     httplib.Config    `yaml:"http" validate:"required"`
}

// Load is a constructor for Config struct
func Load(configPath string) (Config, error) {
	var cfg Config

	data, err := os.ReadFile(configPath)
	if err != nil {
		return cfg, err
	}

	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return cfg, err
	}
	if err := validateConfig(&cfg); err != nil {
		return cfg, err
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
