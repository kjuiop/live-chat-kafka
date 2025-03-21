package route

import (
	"github.com/gin-gonic/gin"
	"live-chat-kafka/api/controller"
)

type RouterConfig struct {
	Engine           *gin.Engine
	SystemController *controller.SystemController
}

func (r *RouterConfig) APISetup() {
	api := r.Engine.Group("/api")
	r.SetupSystemRouter(api)
}

func (r *RouterConfig) SetupSystemRouter(api *gin.RouterGroup) {
	api.GET("/system/health-check", r.SystemController.GetHealth)
}
