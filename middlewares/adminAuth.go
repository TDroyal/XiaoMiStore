package middlewares

import (
	"XiaoMiStore/models"
	"encoding/json"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func InitAdminAuthMiddleware(c *gin.Context) {
	//获取userinfo对应的session
	session := sessions.Default(c)
	uinfo := session.Get("userinfo")
	//类型断言，判断uinfo是不是string
	uinfostr, ok := uinfo.(string)
	if !ok {
		uinfostr = string(uinfostr)
	}
	var userinfoStruct models.Admin
	json.Unmarshal([]byte(uinfostr), &userinfoStruct) //把结构体字符串转成结构体

}
