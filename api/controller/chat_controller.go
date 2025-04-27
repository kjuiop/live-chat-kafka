package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"live-chat-kafka/api/form"
	"live-chat-kafka/internal/domain/chat"
	"live-chat-kafka/internal/domain/chat/types"
	"live-chat-kafka/internal/models"
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

	socket, err := cc.upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		cc.failResponse(c, http.StatusInternalServerError, models.ErrNotFoundChatRoom, fmt.Errorf("not found chat room, roomId : %s, err : %w", roomId, err))
		return
	}

	chatRoom, err := cc.ChatUseCase.GetChatRoom(c, roomId)
	if err != nil {
		cc.failResponse(c, http.StatusBadRequest, models.ErrNotFoundChatRoom, fmt.Errorf("not found chat room, roomId : %s, err : %w", roomId, err))
		return
	}

	if err := cc.ChatUseCase.ServeWs(c, socket, chatRoom, userId); err != nil {
		cc.failResponse(c, http.StatusInternalServerError, models.ErrInternalServerError, fmt.Errorf("failed serve websocket, roomId : %s, err : %w", roomId, err))
		return
	}

	cc.successResponse(c, http.StatusOK, form.APIResponse{
		ErrorCode: 0,
		Message:   "close chat connection",
		Result:    nil,
	})
}
