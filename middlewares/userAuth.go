package middlewares

import (
	"XiaoMiStore/logic"
	"XiaoMiStore/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func InitUserAuthMiddleware(c *gin.Context) {
	// 判断用户有没有登录
	user := models.User{}
	isLogin := logic.Cookie.Get(c, "userinfo", &user)
	if !isLogin || len(user.Phone) != 11 { // 未登录的话
		c.JSON(http.StatusOK, gin.H{
			"message": "用户未登录",
			"status":  "-1",
			"data":    nil,
		})
		c.Abort()
		return
	}
	c.Set("user", user) //供后续handler使用
}
