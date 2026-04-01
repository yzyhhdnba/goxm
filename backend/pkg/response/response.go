package response

import (
	stdhttp "net/http"

	"github.com/gin-gonic/gin"
)

type Envelope struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

func Success(c *gin.Context, data any) {
	c.JSON(stdhttp.StatusOK, Envelope{
		Code:    0,
		Message: "ok",
		Data:    data,
	})
}

func Error(c *gin.Context, status int, code int, message string) {
	c.JSON(status, Envelope{
		Code:    code,
		Message: message,
		Data:    gin.H{},
	})
}
