package mistore

import (
	"XiaoMiStore/dao"
	"XiaoMiStore/logic"
	"XiaoMiStore/models"
	"fmt"
	"net/http"
	"regexp"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// 用户登录 注册的相关控制器

type PassController struct {
	BaseController
}

func (con PassController) GenerateACaptcha(c *gin.Context) {
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

// 注册第一步，发送验证码
func (con PassController) SendCode(c *gin.Context) {
	phone := c.PostForm("phone")
	captchaID := c.PostForm("CaptchaID")
	vertifyValue := c.PostForm("vertifyValue") // 用户输入的图形验证码

	// 1. 验证图形验证码是否正确
	if captchaID == "resend" { // 界面二可以重新获取验证码，但是此时没有图形验证码了
		// 界面二发送验证码
		session := sessions.Default(c)
		sessionVertifyCode := session.Get("vertifyValue")
		sessionVertifyCodeStr, ok := sessionVertifyCode.(string)
		if !ok || sessionVertifyCodeStr != vertifyValue {
			con.Error(c, "请从第一步开始注册，未输入图形验证码", -1, nil)
			return
		}
	} else {
		if result := logic.VerifyCaptcha(captchaID, vertifyValue); !result {
			con.Error(c, "输入的图形验证码错误", -1, nil)
			return
		}
	}

	// 2. 判断手机格式是否合法
	pattern := `^[\d]{11}$` // 正则表达式
	reg := regexp.MustCompile(pattern)
	if result := reg.MatchString(phone); !result {
		con.Error(c, "手机号格式错误，请输入正确的手机号", -1, nil)
		return
	}

	// 3. 验证手机号是否注册过
	var count int64
	if err := dao.DB.Model(&models.User{}).Where("phone = ?", phone).Count(&count).Error; err != nil {
		con.Error(c, "获取验证码失败，请稍后再试", -1, nil)
		return
	}
	if count > 0 {
		con.Error(c, "此手机号已经注册", -1, nil)
		return
	}

	// 4. 判断当前ip地址今天发送短信的次数  防止一个ip用无数个手机号码进行短信轰炸
	currentDay := logic.GetDate() // 20240902
	ip := c.ClientIP()

	var sendCount int64
	if err := dao.DB.Model(&models.UserTemp{}).Where("ip = ? and create_day = ?", ip, currentDay).Count(&sendCount).Error; err != nil {
		con.Error(c, "获取验证码失败，请稍后再试", -1, nil)
		return
	}
	if sendCount >= 5 {
		con.Error(c, "当前ip今日发送短信已达上限", -1, nil)
		return
	}

	// 5. 验证当前手机号今天发送的次数是否合法
	var userTemp []models.UserTemp
	if err := dao.DB.Where("phone = ? and create_day = ?", phone, currentDay).Find(&userTemp).Error; err != nil {
		con.Error(c, "获取验证码失败，请稍后再试", -1, nil)
		return
	}

	if len(userTemp) > 0 {
		if userTemp[0].SendCount >= 3 {
			con.Error(c, "当前手机号今日发送短信已达上限", -1, nil)
			return
		}
		// 1. 生成手机短信验证码
		smsCode := logic.GetRandomNum(6)
		// 调用第三方接口，向此手机号发送验证码
		// todo
		// 参考 https://www.yunpian.com/official/document/sms/zh_CN/introduction_download

		// 2. 服务器保存验证码（保存在session或者redis都可以）
		session := sessions.Default(c)
		session.Set("smsCode", smsCode)
		fmt.Println("手机验证码===", smsCode)
		session.Set("vertifyValue", vertifyValue)
		session.Save()

		// 3. 更新发送短信的次数
		oneUserTemp := models.UserTemp{}
		dao.DB.Where("id = ?", userTemp[0].ID).Find(&oneUserTemp)
		oneUserTemp.SendCount += 1
		dao.DB.Save(&oneUserTemp)
		sign := logic.GetMD5(phone + currentDay) // 主要用于页面跳转，后面根据签名去数据库查询获取手机号
		con.Success(c, "发送短信成功", 0, gin.H{
			"sign": sign,
		})
		return
	} else { // 今天没有发送过短信
		// 1. 生成手机短信验证码
		smsCode := logic.GetRandomNum(6)
		// 调用第三方接口，向此手机号发送验证码
		// todo
		// 参考 https://www.yunpian.com/official/document/sms/zh_CN/introduction_download

		// 2. 服务器保存验证码（保存在session中）
		// 每个用户的 session 数据都是与其请求相关联的，因此不同用户之间的 session 数据是隔离的，不会相互冲突。
		// Gin 使用的默认 session 存储引擎是基于 cookie 的，每个用户的 session 数据都会保存在其自己的 cookie 中，并通过 session ID 进行标识。这样确保了不同用户之间的 session 数据不会相互干扰。
		// 每个用户会创建一个session对象，此对象的key会返回到浏览器（客户端）存在cookie中 ，此对象的value会保存到服务器 在redis中存放在一块单独的地方  浏览器下次访问时会携带 key(cookie)，找到对应的 session(value)
		// 例如：
		// 127.0.0.1:6379> get session_3CB3ODF2U2JMXDO7EMTTDBBREH5OQLJGPU4V73VWYNRPA3COZUYQ
		// "\r\x7f\x04\x01\x02\xff\x80\x00\x01\x10\x01\x10\x00\x00M\xff\x80\x00\x02\x06string\x0c\x0e\x00\x0cvertifyValue\x06string\x0c\x06\x00\x049ins\x06string\x0c\t\x00\asmsCode\x06string\x0c\b\x00\x06814052"

		session := sessions.Default(c)
		// session.Options(sessions.Options{  // 20秒过后，redis中的下面存储的session信息均丢失
		// 	MaxAge: 20,
		// })
		session.Set("smsCode", smsCode)
		fmt.Println("手机验证码===", smsCode)
		session.Set("vertifyValue", vertifyValue)
		session.Save()

		// 3. 记录发送短信的次数
		sign := logic.GetMD5(phone + currentDay) // 主要用于页面跳转，后面根据签名去数据库查询获取手机号
		oneUserTemp := models.UserTemp{
			Ip:        ip,
			Phone:     phone,
			SendCount: 1,
			CreateDay: currentDay,
			Sign:      sign,
		}
		if err := dao.DB.Create(&oneUserTemp).Error; err != nil {
			con.Error(c, "获取验证码失败，请稍后再试", -1, nil)
			return
		}
		con.Success(c, "发送短信成功", 0, gin.H{
			"sign": sign,
		})
		return
	}

}

// 从注册的第一步界面跳到第二步界面(用于页面跳转)，需要利用sign和图形验证码进行验证，验证通过后，才能进入第二步的界面
// 进入步骤二页面时，前端需要将sign和图形验证码传到后端，后端用以验证，验证失败的话，前端重新定位到步骤一页面；验证成功的话，前端渲染步骤二页面
// 进入步骤二页面之前，需要调用VertifyStep1ToStep2
func (con PassController) VertifyStep1ToStep2(c *gin.Context) {
	sign := c.PostForm("sign")
	vertifyValue := c.PostForm("vertifyValue") // 用户输入的图形验证码

	session := sessions.Default(c)
	sessionVertifyCode := session.Get("vertifyValue")
	sessionVertifyCodeStr, ok := sessionVertifyCode.(string)

	// 1. 验证图形验证码是否正确
	if !ok || sessionVertifyCodeStr != vertifyValue {
		con.Error(c, "请获取手机验证码，再进入此页面", -1, nil)
		return
	}

	// 2. 获取sign 判断sign是否合法
	userTemp := []models.UserTemp{}
	dao.DB.Where("sign = ?", sign).Find(&userTemp)
	if len(userTemp) > 0 {
		con.Success(c, "成功进入步骤二", 0, gin.H{
			"phone": userTemp[0].Phone,
			"sign":  sign,
		})
	} else {
		con.Error(c, "请获取手机验证码，再进入此页面", -1, nil)
		return
	}

}

// 注册第二步，验证用户输入的短信验证码是否正确
func (con PassController) VertifySmsCode(c *gin.Context) { //等价于 VertifyStep2ToStep3
	sign := c.PostForm("sign")
	smsCode := c.PostForm("smsCode") // 用户输入的手机验证码

	// 1. 获取sign 判断sign是否合法
	userTemp := []models.UserTemp{}
	dao.DB.Where("sign = ?", sign).Find(&userTemp)
	if len(userTemp) == 0 {
		con.Error(c, "非法请求", -1, nil)
		return
	}

	// 2. 验证手机验证码是否正确  以及是否过期
	session := sessions.Default(c)
	sessionVertifySmsCode := session.Get("smsCode")
	sessionVertifySmsCodeStr, ok := sessionVertifySmsCode.(string)
	if !ok || sessionVertifySmsCodeStr != smsCode {
		con.Error(c, "短信验证码错误", -1, nil)
		return
	}

	// 3. 判断验证码有没有过期   3分钟 （不要把redis中的session数据删掉，应该手动判断）
	nowTime := time.Now()
	timeDiff := nowTime.Sub(userTemp[0].CreatedAt)
	if timeDiff.Minutes() >= 3 { // 以分钟为单位
		con.Error(c, "短信验证码过期", -1, nil)
		return
	}
	con.Success(c, "短信验证码验证通过", 0, gin.H{
		"phone":   userTemp[0].Phone,
		"sign":    sign,
		"smsCode": smsCode,
	}) // 进入到页面三，输入密码即可注册成功。
}

// 进入步骤三页面之前，需要调用VertifyStep2ToStep3
func (con PassController) VertifyStep2ToStep3(c *gin.Context) {
	sign := c.PostForm("sign")
	smsCode := c.PostForm("smsCode") // 用户输入的短信验证码

	// 1. 验证用户输入的短信验证码是否正确
	session := sessions.Default(c)
	sessionVertifySmsCode := session.Get("smsCode")
	sessionVertifySmsCodeStr, ok := sessionVertifySmsCode.(string)
	if !ok || sessionVertifySmsCodeStr != smsCode {
		con.Error(c, "非法进入此页面", -1, nil)
		return
	}

	// 2. 获取sign 判断sign是否合法
	userTemp := []models.UserTemp{}
	dao.DB.Where("sign = ?", sign).Find(&userTemp)
	if len(userTemp) > 0 {
		con.Success(c, "成功进入步骤三", 0, gin.H{
			"phone":   userTemp[0].Phone,
			"sign":    sign,
			"smsCode": smsCode,
		})
	} else {
		con.Error(c, "非法进入此页面", -1, nil)
		return
	}

}

// 注册第三步，输入两次密码完成注册
func (con PassController) DoRegister(c *gin.Context) {
	// 1. 获取前端传过来的数据
	sign := c.PostForm("sign")
	smsCode := c.PostForm("smsCode")
	password := c.PostForm("password")
	rpassword := c.PostForm("rpassword")
	// 2. 验证smsCode是否合法
	session := sessions.Default(c)
	sessionVertifySmsCode := session.Get("smsCode")
	sessionVertifySmsCodeStr, ok := sessionVertifySmsCode.(string)
	if !ok || sessionVertifySmsCodeStr != smsCode {
		con.Error(c, "注册失败", -1, nil)
		return
	}
	// 3. 验证密码是否合法
	if password != rpassword {
		con.Error(c, "注册失败，两次密码不一致", -1, nil)
		return
	}
	if len(password) < 6 {
		con.Error(c, "密码不合法", -1, nil)
		return
	}

	// 4. 验证sign是否合法
	userTemp := []models.UserTemp{}
	dao.DB.Where("sign = ?", sign).Find(&userTemp)
	if len(userTemp) > 0 {
		// 5. 完成注册

		phone := userTemp[0].Phone
		user := models.User{
			Phone:    phone,
			Password: logic.GetMD5(password),
			LastIp:   c.ClientIP(),
			Status:   1,
		}
		if err := dao.DB.Create(&user).Error; err != nil {
			con.Error(c, "注册失败", -1, nil)
			return
		}

		// 执行登录： 将用户登录信息保存至cookie  （不要将密码等信息保存进去，只保存sign等非敏感信息即可）
		logic.Cookie.Set(c, "userinfo", user)
		con.Success(c, "注册成功", 0, nil)
	} else {
		con.Error(c, "非法进入此页面", -1, nil)
		return
	}

}

// 登录
/*
手机号
密码
图形验证码

将用户信息存储在cookie中
*/
func (con PassController) Login(c *gin.Context) {
	phone := c.PostForm("phone")
	password := c.PostForm("password")
	captchaID := c.PostForm("CaptchaID")
	vertifyValue := c.PostForm("vertifyValue") // 用户输入的图形验证码

	// 1. 验证图形验证码
	if result := logic.VerifyCaptcha(captchaID, vertifyValue); !result {
		con.Error(c, "图形验证码错误", -1, nil)
		return
	}

	// 2. 验证账户密码
	user := []models.User{}
	if err := dao.DB.Where("phone = ? and password = ?", phone, logic.GetMD5(password)).Find(&user).Error; err != nil {
		con.Error(c, "登录失败", -1, nil)
		return
	}
	if len(user) == 0 {
		con.Error(c, "手机号或者密码错误", -1, nil)
		return
	}
	// 3. 执行登录（把用户信息写入cookie）
	logic.Cookie.Set(c, "userinfo", &user[0])
	con.Success(c, "用户登录成功", 0, nil)
}

// 退出登录，就是删除cookie里面的用户信息
func (con PassController) LogOut(c *gin.Context) {
	logic.Cookie.Remove(c, "userinfo")
}
