package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
)

func RecoveryErrorReport() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				errMsg := fmt.Sprintf("recovered from panic : %v", err)
				slog.Error(errMsg)
				c.AbortWithStatus(http.StatusInternalServerError)
			}
		}()
		c.Next()
	}
}
