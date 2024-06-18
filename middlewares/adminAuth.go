package middlewares

import (
	"XiaoMiStore/models"
	"encoding/json"
	"fmt"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// 目前所有的用户输入后台管理中心的url都可以直接进入到后台管理中心，因此需要进行权限判断，只有登录的管理员才能进入
func InitAdminAuthMiddleware(c *gin.Context) {
	// 进行权限判断，没有登录的用户 不能进入后台管理中心

	// 1、获取Url访问的地址
	urlPath := c.Request.URL.String() // 如果浏览器访问后台首页http://127.0.0.1/admin/index       则打印  /admin/index
	fmt.Println(urlPath)
	// 2、获取session里面保存的用户信息
	session := sessions.Default(c)
	uinfo := session.Get("userinfo")
	//类型断言，判断uinfo是不是string
	uinfostr, ok := uinfo.(string)
	// 3、判断session中的用户信息是否存在，如果不存在跳转到登录页面，
	if ok {
		var userinfoStruct models.Admin
		err := json.Unmarshal([]byte(uinfostr), &userinfoStruct) //把结构体字符串转成结构体
		if err != nil || userinfoStruct.Username == "" {         //没有登录成功
			fmt.Println(err, userinfoStruct.Username)
			// 访问前端/admin/login页面时，会会从后端路由/admin/generateCaptcha获取验证码
			// 同时，点击登录时，会调用后端的路由/admin/doLogin提交用户信息
			if urlPath != "/admin/login" && urlPath != "/admin/doLogin" && urlPath != "/admin/generateCaptcha" {
				c.Redirect(302, "/admin/login")
			}
		}

	} else { //表示用户没有登录，跳转到前端的登录页面
		// 4、如果session不存在，判断当前访问的URL是否是login logout captcha，
		if urlPath != "/admin/login" && urlPath != "/admin/doLogin" && urlPath != "/admin/generateCaptcha" {
			c.Redirect(302, "/admin/login")
		}
	}

}
