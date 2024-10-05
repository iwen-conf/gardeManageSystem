/*
* Editor:KaiWen
* PATH:utils/response.go
* Description:
 */

// utils/response.go

package utils

import (
	"github.com/gin-gonic/gin"
)

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func ResponseSuccess(c *gin.Context, message string, data interface{}) {
	c.JSON(200, Response{
		Code:    200,
		Message: message,
		Data:    data,
	})
}

func ResponseError(c *gin.Context, message string) {
	c.JSON(400, Response{
		Code:    400,
		Message: message,
		Data:    nil,
	})
}
