package config

import "github.com/caarlos0/env/v11"

type Config struct {
	App      App      `envPrefix:"TODO"`
	Postgres Postgres `envPrefix:"POSTGRES"`
}

type App struct {
	Host string `env:"HOST" envDefault:"localhost"`
	Port int    `env:"PORT" envDefault:"8080"`
}

type Postgres struct {
	Host string `env:"HOST" envDefault:"localhost"`
	Port int    `env:"PORT" envDefault:"5432"`
}

func Load() (*Config, error) {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}
