package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"live-chat-kafka/api/form"
	"live-chat-kafka/internal/domain/chat"
	"live-chat-kafka/internal/domain/chat/types"
	"live-chat-kafka/internal/models"
	"log/slog"
	"net/http"
)

type ChatController struct {
	upgrader    *websocket.Upgrader
	ChatUseCase chat.UseCase
}

func NewChatController(chatUseCase chat.UseCase) *ChatController {
	return &ChatController{
		upgrader: &websocket.Upgrader{
			ReadBufferSize:  types.SocketBufferSize,
			WriteBufferSize: types.MessageBufferSize,
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
		ChatUseCase: chatUseCase,
	}
}

func (cc *ChatController) successResponse(c *gin.Context, statusCode int, data interface{}) {
	c.JSON(statusCode, form.APIResponse{
		ErrorCode: models.NoError,
		Message:   models.GetCustomMessage(models.NoError),
		Result:    data,
	})
}

func (cc *ChatController) failResponse(c *gin.Context, statusCode, errorCode int, err error) {

	logMessage := models.GetCustomErrMessage(errorCode, err.Error())
	c.Errors = append(c.Errors, &gin.Error{
		Err:  fmt.Errorf(logMessage),
		Type: gin.ErrorTypePrivate,
	})

	c.JSON(statusCode, form.APIResponse{
		ErrorCode: errorCode,
		Message:   models.GetCustomMessage(errorCode),
	})
}

func (cc *ChatController) ServeWS(c *gin.Context) {

	roomId := c.Param("room_id")
	userId := c.Param("user_id")

	logger := slog.With("roomId", roomId, "userId", userId)

	socket, err := cc.upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		cc.failResponse(c, http.StatusInternalServerError, models.ErrNotFoundChatRoom, fmt.Errorf("not found chat room, roomId : %s, err : %w", roomId, err))
		return
	}

	chatRoom, err := cc.ChatUseCase.GetChatRoom(c, roomId)
	if err != nil {
		logger.Error("not found chat room", "roomId", roomId, "userId", userId, "err", err)
		if err := socket.Close(); err != nil {
			logger.Warn("socket close error", "roomId", roomId, "userId", userId, "err", err)
		}
		return
	}

	if err := cc.ChatUseCase.ServeWs(c, socket, chatRoom, userId); err != nil {
		logger.Error("failed to serve websocket", "roomId", roomId, "userId", userId, "err", err)
		return
	}
}
