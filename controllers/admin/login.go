package admin

import (
	"XiaoMiStore/dao"
	"XiaoMiStore/logic"
	"XiaoMiStore/models"
	"encoding/json"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type LoginController struct {
	BaseController
}

func (con LoginController) Index(c *gin.Context) {
	var admin = &models.Admin{
		Username: "royal",
		Password: logic.GetMD5("123456"),
		Mobile:   "18275385029",
		Email:    "123456@qq.com",
		Status:   1,
		RoleID:   2,
		IsSuper:  0,
	}
	dao.DB.Create(admin)

	c.JSON(http.StatusOK, gin.H{
		"message": "ok",
		"status":  0,
		"data":    nil,
	})
}

func (con LoginController) GenerateACaptcha(c *gin.Context) {
	id, b64s, answer, err := logic.GenerateCaptcha()

	if err != nil { //错误后续处理
		panic(err)
	}
	c.JSON(http.StatusOK, gin.H{
		"id":           id, //返回给前端的这个id可以存放在一个input框中，其中这个input框隐藏掉即可(type="hidden")，用户提交登录的时候需要将这个id一起传到后端做验证
		"image_base64": b64s,
		"answer":       answer,
	})
}

func (con LoginController) Login(c *gin.Context) {
	// 获取前端表单传过来的账户，密码，图形验证码id，以及用户输入的验证码   也可以用c.ShouldBind()
	username := c.PostForm("username")
	password := c.PostForm("password")
	captchaID := c.PostForm("CaptchaID")
	vertifyValue := c.PostForm("vertifyValue")

	//1. 验证验证码是否正确 (防止一直攻击，请求数据库)
	verifyCode_result := logic.VerifyCaptcha(captchaID, vertifyValue) //验证码的验证结果
	if !verifyCode_result {
		c.JSON(http.StatusOK, gin.H{
			"status":  -1,
			"message": "验证码错误",
			"data":    nil,
		})
		return
	}

	//2.验证用户账号密码是否存在
	var userinfo = []models.Admin{} //定义为切片方便用len方法判断长度
	if err := dao.DB.Where("username = ? and password = ?", username, logic.GetMD5(password)).First(&userinfo).Error; err != nil || len(userinfo) != 1 {
		c.JSON(http.StatusOK, gin.H{
			"status":  -1,
			"message": "用户名或者密码错误",
			"data":    err.Error(),
		})
		return
	}

	//3.登录成功，保存用户信息
	// 后台管理系统的用户信息为了安全起见，是将用户信息保存到session里面的，session保存在服务器上    后续改为JWT
	// 官网地址  https://github.com/gin-contrib/sessions

	session := sessions.Default(c)
	// 注意：session.set无法直接保存结构体对应的切片   把结构体转换为json字符串
	userinfoSlice, _ := json.Marshal(userinfo[0])
	// session.Set("userinfo", userinfo[0])   //设置session   //之前报错：gob: type not registered for interface: models.Admin
	session.Set("userinfo", string(userinfoSlice))
	if err := session.Save(); err != nil { //保存session
		c.JSON(http.StatusOK, gin.H{
			"status":  -1,
			"message": "系统出错，请稍后再试",
			"data":    err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  0,
		"message": "登录成功",
		"data":    nil,
	})
}

func (con LoginController) Logout(c *gin.Context) {
	// 1.销毁session
	session := sessions.Default(c)
	// 它只会从当前请求的会话中删除"userinfo"键，而不会从Redis中删除对应的会话数据。(后续思考怎么解决)
	session.Delete("userinfo")             //删除session
	if err := session.Save(); err != nil { //保存session
		c.JSON(http.StatusOK, gin.H{
			"status":  -1,
			"message": "系统出错，请稍后再试",
			"data":    err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  0,
		"message": "成功退出登录",
		"data":    nil,
	})
}
