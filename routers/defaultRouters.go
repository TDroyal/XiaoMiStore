package routers

import (
	"XiaoMiStore/controllers/mistore"

	"github.com/gin-gonic/gin"
)

func SetupDefaultRouters(r *gin.Engine) {
	defaultRouters := r.Group("/")
	{
		defaultRouters.GET("/", mistore.DefaultController{}.Index)
	}
}
