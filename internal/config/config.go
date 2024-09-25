package config

import (
	"github.com/caarlos0/env/v11"
)

type Agent struct {
	AppConfig
	AgentConfig
}

type WebServer struct {
	AppConfig
	WebServerConfig
}

func Load[T any]() (T, error) {
	var cfg T
	if err := env.Parse(&cfg); err != nil {
		return cfg, err
	}

	return cfg, nil
}
