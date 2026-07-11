package response

import "github.com/gin-gonic/gin"

type Envelope struct {
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

func OK(c *gin.Context, status int, message string, data interface{}) {
	c.JSON(status, Envelope{Data: data, Message: message})
}

func Error(c *gin.Context, status int, message string) {
	c.JSON(status, Envelope{Data: nil, Message: message})
}

func ValidationError(c *gin.Context, status int, message string, fieldErrors map[string]string) {
	c.JSON(status, Envelope{Data: fieldErrors, Message: message})
}
