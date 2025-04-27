package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"live-chat-kafka/api/form"
	"live-chat-kafka/internal/domain/system"
	"live-chat-kafka/internal/models"
	"net/http"
)

type SystemController struct {
	SystemUseCase system.UseCase
}

func NewSystemController(useCase system.UseCase) *SystemController {
	return &SystemController{
		SystemUseCase: useCase,
	}
}

func (s *SystemController) successResponse(c *gin.Context, statusCode int, data interface{}) {
	c.JSON(statusCode, form.APIResponse{
		ErrorCode: models.NoError,
		Message:   models.GetCustomMessage(models.NoError),
		Result:    data,
	})
}

func (s *SystemController) failResponse(c *gin.Context, statusCode, errorCode int, err error) {

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

func (s *SystemController) GetHealth(c *gin.Context) {
	c.String(http.StatusOK, "ok")
}

func (s *SystemController) GetServerList(c *gin.Context) {

	list, err := s.SystemUseCase.GetServerList()
	if err != nil {
		s.failResponse(c, http.StatusInternalServerError, models.ErrInternalServerError, fmt.Errorf("get server list occur err : %w", err))
		return
	}

	if len(list) == 0 {
		s.successResponse(c, http.StatusOK, nil)
		return
	}

	s.successResponse(c, http.StatusOK, list)
}
