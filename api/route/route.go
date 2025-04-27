package route

import (
	"github.com/gin-gonic/gin"
	"live-chat-kafka/api/controller"
)

type RouterConfig struct {
	Engine           *gin.Engine
	SystemController *controller.SystemController
	RoomController   *controller.RoomController
	ChatController   *controller.ChatController
}

func (r *RouterConfig) APISetup() {
	api := r.Engine.Group("/api")
	r.SetupRoomRouter(api)
	r.SetupSystemRouter(api)
}

func (r *RouterConfig) WSSetup() {
	ws := r.Engine.Group("/ws")
	r.SetupWSGroup(ws)
}

func (r *RouterConfig) SetupRoomRouter(api *gin.RouterGroup) {
	api.POST("/rooms", r.RoomController.CreateChatRoom)
	api.GET("/rooms/:room_id", r.RoomController.GetChatRoom)
	api.PUT("/rooms/:room_id", r.RoomController.UpdateChatRoom)
	api.DELETE("/rooms/:room_id", r.RoomController.DeleteChatRoom)
	api.GET("/rooms/id", r.RoomController.GetChatRoomId)
}

func (r *RouterConfig) SetupSystemRouter(api *gin.RouterGroup) {
	api.GET("/system/health-check", r.SystemController.GetHealth)
	api.GET("/system/server-list", r.SystemController.GetServerList)
}

func (r *RouterConfig) SetupWSGroup(ws *gin.RouterGroup) {
	ws.GET("/chat/join/rooms/:room_id/user/:user_id", r.ChatController.ServeWS)
}
