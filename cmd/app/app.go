package app

import (
	"context"
	"live-chat-kafka/config"
	"live-chat-kafka/logger"
	"log"
	"sync"
)

type App struct {
	cfg *config.EnvConfig
}

func NewApplication(ctx context.Context) *App {

	cfg, err := config.LoadEnvConfig()
	if err != nil {
		log.Fatalf("fail to read config err : %v", err)
	}

	if err := logger.SlogInit(cfg.Logger); err != nil {
		log.Fatalf("fail to init slog err : %v", err)
	}

	return &App{cfg}
}

func (a *App) Start(wg *sync.WaitGroup) {
	defer wg.Done()
}

func (a *App) Stop(ctx context.Context) {
}
