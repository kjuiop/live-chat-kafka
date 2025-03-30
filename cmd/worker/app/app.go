package app

import (
	"context"
	"live-chat-kafka/config"
	"live-chat-kafka/internal/server"
	"live-chat-kafka/logger"
	"log"
	"sync"
)

type App struct {
	srv server.Client
}

func NewApplication(ctx context.Context) *App {

	cfg, err := config.LoadEnvConfig()
	if err != nil {
		log.Fatalf("fail to read config err : %v", err)
	}

	if err := logger.SlogInit(cfg.Logger); err != nil {
		log.Fatalf("fail to init slog err : %v", err)
	}

	srv := server.NewGinServer(cfg)

	return &App{
		srv: srv,
	}
}

func (a *App) Start(wg *sync.WaitGroup) {
	defer wg.Done()
	a.srv.Run()
}

func (a *App) Stop(ctx context.Context) {
	a.srv.Shutdown(ctx)
}
