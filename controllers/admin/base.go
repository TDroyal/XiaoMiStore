package admin

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type BaseController struct{}

// 将成功和失败的返回状态和数据封装在一个basecontroller里面
func (con BaseController) Success(c *gin.Context, message string, status int, data interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"message": message,
		"status":  status,
		"data":    data,
	})
}

func (con BaseController) Error(c *gin.Context, message string, status int, data interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"message": message,
		"status":  status,
		"data":    data,
	})
}
