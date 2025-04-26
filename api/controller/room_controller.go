package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"live-chat-kafka/api/form"
	"live-chat-kafka/config"
	"live-chat-kafka/internal/domain/room"
	"live-chat-kafka/internal/models"
	"net/http"
)

type RoomController struct {
	cfg         config.Policy
	RoomUseCase room.UseCase
}

func NewRoomController(cfg config.Policy, useCase room.UseCase) *RoomController {
	return &RoomController{
		cfg:         cfg,
		RoomUseCase: useCase,
	}
}

func (r *RoomController) successResponse(c *gin.Context, statusCode int, data interface{}) {
	c.JSON(statusCode, form.APIResponse{
		ErrorCode: models.NoError,
		Message:   models.GetCustomMessage(models.NoError),
		Result:    data,
	})
}

func (r *RoomController) failResponse(c *gin.Context, statusCode, errorCode int, err error) {
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

func (r *RoomController) CreateChatRoom(c *gin.Context) {
	req := form.RoomRequest{}
	if err := c.ShouldBind(&req); err != nil {
		r.failResponse(c, http.StatusBadRequest, models.ErrParsing, fmt.Errorf("CreateRoom json parsing err : %w", err))
		return
	}

	roomInfo := room.NewRoomInfo(req, r.cfg.Prefix)
	if err := r.RoomUseCase.CreateChatRoom(c, *roomInfo); err != nil {
		r.failResponse(c, http.StatusInternalServerError, models.ErrRedisHMSETError, fmt.Errorf("CreateRoom HMSET err : %w", err))
		return
	}

	if err := r.RoomUseCase.RegisterRoomId(c, *roomInfo); err != nil {
		r.failResponse(c, http.StatusInternalServerError, models.ErrRedisHMSETError, fmt.Errorf("register room id HMSET err : %w", err))
		return
	}

	roomRes := form.RoomResponse{
		RoomId:       roomInfo.RoomId,
		CustomerId:   roomInfo.CustomerId,
		ChannelKey:   roomInfo.ChannelKey,
		BroadcastKey: roomInfo.BroadcastKey,
		CreatedAt:    roomInfo.CreatedAt,
	}

	r.successResponse(c, http.StatusCreated, roomRes)
}
