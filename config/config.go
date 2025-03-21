package config

import (
	"github.com/kelseyhightower/envconfig"
)

type EnvConfig struct {
	Logger Logger
	Server Server
}

type Server struct {
	Mode           string `envconfig:"LCK_ENV" default:"dev"`
	Port           string `envconfig:"LCK_SERVER_PORT" default:"8090"`
	TrustedProxies string `envconfig:"LCK_TRUSTED_PROXIES" default:"127.0.0.1/32"`
}

type Logger struct {
	Level       string `envconfig:"LCK_LOG_LEVEL" default:"debug"`
	Path        string `envconfig:"LCK_LOG_PATH" default:"./logs/access.log"`
	PrintStdOut bool   `envconfig:"LCK_LOG_STDOUT" default:"true"`
}

func LoadEnvConfig() (*EnvConfig, error) {
	var config EnvConfig
	if err := envconfig.Process("lck", &config); err != nil {
		return nil, err
	}

	if err := config.CheckValid(); err != nil {
		return nil, err
	}

	return &config, nil
}

func (c *EnvConfig) CheckValid() error {
	return nil
}
