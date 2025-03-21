package server

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"live-chat-kafka/api/middleware"
	"live-chat-kafka/config"
	"log"
	"log/slog"
	"net/http"
	"strings"
)

type Gin struct {
	srv    *http.Server
	router *gin.Engine
	cfg    config.Server
}

func NewGinServer(cfg *config.EnvConfig) Client {

	serverCfg := cfg.Server
	router := getGinEngine(serverCfg.Mode)

	if err := router.SetTrustedProxies(strings.Split(serverCfg.TrustedProxies, ",")); err != nil {
		log.Fatalf("failed set trust proxies, err : %v", err)
	}

	router.Use(middleware.LoggingMiddleware)
	router.Use(middleware.RecoveryErrorReport())
	router.Use(middleware.SetCorsPolicy())

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", serverCfg.Port),
		Handler: router,
	}

	return &Gin{
		srv:    srv,
		cfg:    cfg.Server,
		router: router,
	}
}

func (g *Gin) Run() {
	err := g.srv.ListenAndServe()
	if errors.Is(err, http.ErrServerClosed) {
		slog.Debug("server close")
	} else {
		slog.Error("run server error", "error", err)
	}
}

func (g *Gin) Shutdown(ctx context.Context) {
	if err := g.srv.Shutdown(ctx); err != nil {
		slog.Error("error during server shutdown", "error", err)
	}
}

func (g *Gin) GetEngine() *gin.Engine {
	return g.router
}

func getGinEngine(mode string) *gin.Engine {
	switch mode {
	case "prod":
		return gin.New()
	case "test":
		gin.SetMode(gin.TestMode)
		return gin.Default()
	default:
		return gin.Default()
	}
}
