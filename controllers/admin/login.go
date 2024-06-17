package admin

import (
	"XiaoMiStore/logic"
	"net/http"

	"github.com/gin-gonic/gin"
)

type LoginController struct {
	BaseController
}

func (con LoginController) Index(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "ok",
		"status":  0,
		"data":    nil,
	})
}

func (con LoginController) Login(c *gin.Context) {
	// 获取前端表单传过来的账户，密码，图形验证码id，以及用户输入的验证码   也可以用c.ShouldBind()
	captchaID := c.PostForm("CaptchaID")
	vertifyValue := c.PostForm("vertifyValue")
	verify_result := logic.VerifyCaptcha(captchaID, vertifyValue) //验证结果
	if !verify_result {
		c.JSON(http.StatusOK, gin.H{
			"status":  -1,
			"message": "验证码错误",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  0,
		"message": verify_result,
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
