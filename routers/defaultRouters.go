package routers

import (
	"XiaoMiStore/controllers/admin"

	"github.com/gin-gonic/gin"
)

func SetupDefaultRouters(r *gin.Engine) {
	defaultRouters := r.Group("/")
	{
		defaultRouters.GET("/", admin.LoginController{}.Index)
	}
}
