package routers

import (
	// "XiaoMiStore/models"
	// "encoding/gob"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
)

// session.Set("userinfo", userinfo[0])   //设置session   //之前报错：gob: type not registered for interface: models.Admin
// 解决方法1
// func init() { //通过在init()函数中注册models.Admin类型，您将确保gob编码器能够正确地序列化和反序列化该类型的对象。
// 	gob.Register(models.Admin{})
// }

func SetupRouter() *gin.Engine { //父路由
	r := gin.Default()

	//配置redis存储的session中间件  基于 Redis 存储 Session
	// https://github.com/gin-contrib/sessions
	store, err := redis.NewStore(10, "tcp", "localhost:6379", "123456", []byte("secret_royal_111"))
	if err != nil {
		panic(err)
	}
	r.Use(sessions.Sessions("mysession", store))

	//把所有的子路由注册都放在这里来
	SetupDefaultRouters(r)
	SetupAdminRouters(r)

	return r
}
