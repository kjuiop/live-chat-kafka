package app

import (
	"context"
	"live-chat-kafka/api/controller"
	"live-chat-kafka/api/route"
	"live-chat-kafka/config"
	"live-chat-kafka/internal/database"
	rr "live-chat-kafka/internal/domain/room/repository"
	ru "live-chat-kafka/internal/domain/room/usecase"
	"live-chat-kafka/internal/domain/system"
	sps "live-chat-kafka/internal/domain/system/pubsub"
	sr "live-chat-kafka/internal/domain/system/repository"
	su "live-chat-kafka/internal/domain/system/usecase"
	"live-chat-kafka/internal/message_queue"
	"live-chat-kafka/internal/server"
	"live-chat-kafka/logger"
	"log"
	"log/slog"
	"sync"
	"time"
)

type App struct {
	cfg *config.EnvConfig
	srv server.Client
	mq  message_queue.Client
	db  database.Client
	su  system.UseCase
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
		cfg: cfg,
		srv: srv,
		mq:  mq,
		db:  db,
	}
	app.setupRouter()
	if err := app.initProcess(); err != nil {
		log.Fatalf("failed initialized process, err : %v", err)
	}

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

	timeout := time.Duration(a.cfg.Policy.ContextTimeout) * time.Second

	systemRepo := sr.NewSystemRepository(a.db)
	roomRepo := rr.NewRoomRepository(a.db)

	systemPubSub := sps.NewSystemPubSub(a.cfg.Kafka, a.mq)

	systemUseCase := su.NewSystemUseCase(systemRepo, systemPubSub)
	roomUseCase := ru.NewRoomUseCase(roomRepo, timeout)

	systemController := controller.NewSystemController(systemUseCase)
	roomController := controller.NewRoomController(a.cfg.Policy, roomUseCase)

	router := route.RouterConfig{
		Engine:           a.srv.GetEngine(),
		SystemController: systemController,
		RoomController:   roomController,
	}
	router.APISetup()

	a.su = systemUseCase
}

func (a *App) LoopServerInfo(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		select {
		case <-ctx.Done():
			slog.Debug("close Loop Sub Kafka goroutine")
			return
		default:
			event, err := a.su.LoopSubKafka(a.cfg.Kafka.ConsumerTimeout)
			if err != nil {
				slog.Error("received event error", "error", err)
				continue
			}

			if event == nil {
				continue
			}

			slog.Debug("received event", "event", string(event.Value))
		}
	}
}

func (a *App) initProcess() error {

	if err := a.su.RegisterSubTopic("chat"); err != nil {
		return err
	}

	return nil
}
