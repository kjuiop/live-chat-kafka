package app

import (
	"context"
	"live-chat-kafka/api/controller"
	"live-chat-kafka/api/route"
	"live-chat-kafka/config"
	"live-chat-kafka/internal/database"
	sps "live-chat-kafka/internal/domain/system/pubsub"
	sr "live-chat-kafka/internal/domain/system/repository"
	su "live-chat-kafka/internal/domain/system/usecase"
	"live-chat-kafka/internal/message_queue"
	"live-chat-kafka/internal/server"
	"live-chat-kafka/logger"
	"log"
	"sync"
)

type App struct {
	cfg *config.EnvConfig
	srv server.Client
	mq  message_queue.Client
	db  database.Client
}

func NewApplication(ctx context.Context) *App {

	cfg, err := config.LoadEnvConfig()
	if err != nil {
		log.Fatalf("fail to read config err : %v", err)
	}

	if err := logger.SlogInit(cfg.Logger); err != nil {
		log.Fatalf("fail to init slog err : %v", err)
	}

	db, err := database.NewRedisSingleClient(ctx, cfg.Redis)
	if err != nil {
		log.Fatalf("fail to init redis err : %v", err)
	}

	mq, err := message_queue.NewKafkaConsumerClient(cfg.Kafka)
	if err != nil {
		log.Fatalf("fail to connect kafka client")
	}

	srv := server.NewGinServer(cfg)

	app := &App{
		cfg,
		srv,
		mq,
		db,
	}
	app.setupRouter()

	return app
}

func (a *App) Start(wg *sync.WaitGroup) {
	defer wg.Done()
	a.srv.Run()
}

func (a *App) Stop(ctx context.Context) {
	a.srv.Shutdown(ctx)
}

func (a *App) setupRouter() {

	systemRepo := sr.NewSystemRepository(a.db)

	systemPubSub := sps.NewSystemPubSub(a.cfg.Kafka, a.mq)

	systemUseCase := su.NewSystemUseCase(systemRepo, systemPubSub)

	systemController := controller.NewSystemController(systemUseCase)

	router := route.RouterConfig{
		Engine:           a.srv.GetEngine(),
		SystemController: systemController,
	}
	router.APISetup()
}
