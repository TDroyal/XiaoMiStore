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
				c.Redirect(302, "/admin/login") //重定向的是后端的ip，这种不是前后端分离，前后端分离项目，应该往前端返回一个信息告诉前端重定向到登录页面的路由
				c.Abort()                       //不调用后续的handler处理函数
			}
		} else { //前后端分离的项目，应该由前端判断权限问题。

			//管理员登录成功，做权限判断，没有对应权限的管理员访问相应的页面（url）时应该被拒绝
			//1、获取当前管理员对应角色的权限

			//2、再拿到管理员当前访问的地址

			//3、再看该url对应的权限id是否在当前管理员的权限里面
			c.Next()
		}

	} else { //表示用户没有登录，跳转到前端的登录页面
		// 4、如果session不存在，判断当前访问的URL是否是login logout captcha，
		if urlPath != "/admin/login" && urlPath != "/admin/doLogin" && urlPath != "/admin/generateCaptcha" {
			c.Redirect(302, "/admin/login")
			c.Abort()
		}
	}

}
