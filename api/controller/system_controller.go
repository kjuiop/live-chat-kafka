package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"live-chat-kafka/api/form"
	"live-chat-kafka/internal/models"
	"net/http"
)

type SystemController struct {
}

func NewSystemController() *SystemController {
	return &SystemController{}
}

func (s *SystemController) successResponse(c *gin.Context, statusCode int, data interface{}) {
	c.JSON(statusCode, form.ApiResponse{
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

	c.JSON(statusCode, form.ApiResponse{
		ErrorCode: errorCode,
		Message:   models.GetCustomMessage(errorCode),
	})
}

func (s *SystemController) GetHealth(c *gin.Context) {
	s.successResponse(c, http.StatusOK, nil)
	return
}
