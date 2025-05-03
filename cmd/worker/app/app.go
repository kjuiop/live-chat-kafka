package app

import (
	"context"
	"fmt"
	"live-chat-kafka/api/controller"
	"live-chat-kafka/api/route"
	"live-chat-kafka/config"
	"live-chat-kafka/internal/database"
	cps "live-chat-kafka/internal/domain/chat/pubsub"
	cu "live-chat-kafka/internal/domain/chat/usecase"
	rps "live-chat-kafka/internal/domain/room/pubsub"
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
	"net"
	"sync"
	"time"
)

type App struct {
	cfg  *config.EnvConfig
	srv  server.Client
	mq   message_queue.Client
	db   database.Client
	su   system.UseCase
	addr string
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

	mq, err := message_queue.NewKafkaClient(cfg.Kafka, true, false)
	if err != nil {
		log.Fatalf("fail to create kafka producer err : %v", err)
	}

	srv := server.NewGinServer(cfg)

	app := &App{
		cfg: cfg,
		srv: srv,
		mq:  mq,
		db:  db,
	}

	app.setupRouter()
	app.registerServer()

	return app
}

func (a *App) Start(wg *sync.WaitGroup) {
	defer wg.Done()
	a.srv.Run()

}

func (a *App) Stop(ctx context.Context) {
	a.su.PublishServerStatusEvent(a.addr, false)

	a.srv.Shutdown(ctx)
	a.db.Close()
	a.mq.Close()
}

func (a *App) setupRouter() {

	timeout := time.Duration(a.cfg.Policy.ContextTimeout) * time.Second

	systemRepo := sr.NewSystemRepository(a.db)

	systemPubSub := sps.NewSystemPubSub(a.cfg.Kafka, a.mq)
	roomPubSub := rps.NewRoomPubSub(a.cfg.Kafka, a.mq)
	chatPubSub := cps.NewChatPubSub(a.cfg.Kafka, a.mq)

	roomRepo := rr.NewRoomRepository(a.db)

	systemUseCase := su.NewSystemUseCase(systemRepo, systemPubSub)
	roomUseCase := ru.NewRoomUseCase(roomRepo, timeout, roomPubSub)
	chatUseCase := cu.NewChatUseCase(roomUseCase, timeout, chatPubSub)

	chatController := controller.NewChatController(chatUseCase)

	router := route.RouterConfig{
		Engine:         a.srv.GetEngine(),
		ChatController: chatController,
	}
	router.WSSetup()

	a.su = systemUseCase
}

func (a *App) registerServer() {

	addrs, err := net.InterfaceAddrs()
	if err != nil {
		log.Fatalf("failed parsing ip address, err : %v", err)
	}

	var ip net.IP
	for _, addr := range addrs {
		if ipNet, ok := addr.(*net.IPNet); ok {
			if !ipNet.IP.IsLoopback() && ipNet.IP.To4() != nil {
				ip = ipNet.IP
				break
			}
		}
	}

	if ip == nil {
		log.Fatalln("no ip address found")
	}

	addr := fmt.Sprintf("%s:%s", ip.String(), a.cfg.Server.Port)
	if err := a.su.SetChatServerInfo(addr, true); err != nil {
		log.Fatalf("failed register server info, address : %s, err : %v", addr, err)
	}
	a.addr = addr
	a.su.PublishServerStatusEvent(addr, true)

}
