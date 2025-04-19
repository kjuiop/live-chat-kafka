package config

import (
	"github.com/kelseyhightower/envconfig"
)

type APIEnvConfig struct {
	Logger Logger
	Server APIServer
	Redis  Redis
	Kafka  Kafka
}

type WorkerEnvConfig struct {
	Logger Logger
	Server WorkerServer
	Redis  Redis
	Kafka  Kafka
}

type APIServer struct {
	Mode           string `envconfig:"LCK_ENV" default:"dev"`
	Port           string `envconfig:"LCK_SERVER_PORT" default:"8090"`
	TrustedProxies string `envconfig:"LCK_TRUSTED_PROXIES" default:"127.0.0.1/32"`
}

type WorkerServer struct {
	Mode           string `envconfig:"LCK_ENV" default:"dev"`
	Port           string `envconfig:"LCK_SERVER_PORT" default:"8090"`
	TrustedProxies string `envconfig:"LCK_TRUSTED_PROXIES" default:"127.0.0.1/32"`
}

type Logger struct {
	Level       string `envconfig:"LCK_LOG_LEVEL" default:"debug"`
	Path        string `envconfig:"LCK_LOG_PATH" default:"./logs/access.log"`
	PrintStdOut bool   `envconfig:"LCK_LOG_STDOUT" default:"true"`
}

type Redis struct {
	Mode     string `envconfig:"LCK_REDIS_MODE" default:"single"`
	Addr     string `envconfig:"LCK_REDIS_ADDR" default:":6379"`
	Password string `envconfig:"LCK_REDIS_PASSWORD"`
	Masters  string `envconfig:"LCK_REDIS_MASTERS"`
	PoolSize int    `envconfig:"LCK_REDIS_POOL_SIZE" default:"100"`
}

type Kafka struct {
	URL             string `envconfig:"LCK_KAFKA_URL" default:"localhost:9292"`
	GroupID         string `envconfig:"LCK_KAFKA_GROUP_ID" default:"chat-consumer-1"`
	ClientID        string `envconfig:"LCK_KAFKA_CLIENT_ID" default:"chat-producer-1"`
	ConsumerTimeout int    `envconfig:"LCK_KAFKA_CONSUMER_TIMEOUT" default:"1000"`
}

func LoadAPIEnvConfig() (*APIEnvConfig, error) {
	var config APIEnvConfig
	if err := envconfig.Process("lck", &config); err != nil {
		return nil, err
	}

	if err := config.CheckValid(); err != nil {
		return nil, err
	}

	return &config, nil
}

func (a *APIEnvConfig) CheckValid() error {
	return nil
}

func LoadWorkerEnvConfig() (*WorkerEnvConfig, error) {
	var config WorkerEnvConfig
	if err := envconfig.Process("lck", &config); err != nil {
		return nil, err
	}

	if err := config.CheckValid(); err != nil {
		return nil, err
	}

	return &config, nil
}

func (w *WorkerEnvConfig) CheckValid() error {
	return nil
}
