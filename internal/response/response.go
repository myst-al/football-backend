package response

import (
	apperror "football-backend/internal/errors"

	"github.com/gin-gonic/gin"
)

type APIResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func Success(c *gin.Context, code int, message string, data interface{}) {
	c.JSON(code, APIResponse{
		Code:    code,
		Message: message,
		Data:    data,
	})
}

func Error(c *gin.Context, code int, message string) {
	c.JSON(code, APIResponse{
		Code:    code,
		Message: message,
		Data:    nil,
	})
}

func FromError(c *gin.Context, err error) {
	if appErr, ok := err.(*apperror.AppError); ok {
		Error(c, appErr.Code, appErr.Message)
		return
	}

	Error(c, 500, "Terjadi kesalahan pada server")
}
