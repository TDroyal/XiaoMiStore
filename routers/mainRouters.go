package routers

import "github.com/gin-gonic/gin"

func SetupRouter() *gin.Engine { //父路由
	r := gin.Default()

	//把所有的子路由注册都放在这里来
	SetupDefaultRouters(r)
	SetupAdminRouters(r)

	return r
}
