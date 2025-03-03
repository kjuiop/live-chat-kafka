package server

import (
	"context"
	"github.com/gin-gonic/gin"
)

type Client interface {
	Run()
	Shutdown(ctx context.Context)
	GetEngine() *gin.Engine
}
