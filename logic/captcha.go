package logic

// 封装图形验证码
// 官方文档 https://github.com/mojocn/base64Captcha
import (
	"context"
	"fmt"
	"image/color"
	"time"

	pb "XiaoMiStore/proto/captcha"

	"github.com/mojocn/base64Captcha"
)

// 创建默认store，将生成的验证码保存在服务器
// var store = base64Captcha.DefaultMemStore //后续可以把验证码保存在redis里面

// 自定义 store 需要实现 Store 这个接口
// type Store interface {
// 	// Set sets the digits for the captcha id.
// 	Set(id string, value string)

// 	// Get returns stored digits for the captcha id. Clear indicates
// 	// whether the captcha must be deleted from the store.
// 	Get(id string, clear bool) string

//     //Verify captcha's answer directly
// 	Verify(id, answer string, clear bool) bool
// }

type RedisMemStore struct{}

// 拼接一个前缀
const (
	CaptchaKeyPrefix string = "captcha:"
)

func (r RedisMemStore) Set(id string, value string) error { // 图形验证码存储在redis中，两分钟过期
	var key = CaptchaKeyPrefix + id
	if err := RedisSet(key, value, 2*time.Minute); err != nil {
		return err
	}
	return nil
}

func (r RedisMemStore) Get(id string, clear bool) string {
	var key = CaptchaKeyPrefix + id
	val, err := RedisGet(key)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	if clear { // 取值之后就从redis中删除此验证码
		err = RedisDel(key)
		if err != nil {
			fmt.Println(err)
			return ""
		}
	}
	return val
}

func (r RedisMemStore) Verify(id, answer string, clear bool) bool {
	val := r.Get(id, clear)
	return val == answer
}

// 配置RedisStore      RedisStore实现base64Captcha.Store接口（需要实现其中的三个方法）
var store base64Captcha.Store = RedisMemStore{}

// 获取验证码   //生成字符串的验证码
func GenerateCaptchaLocal() (id, b64s, answer string, err error) {
	//配置验证码信息   https://captcha.mojotv.cn/.netlify/functions/captcha
	DriverString := base64Captcha.DriverString{
		Height:          40,
		Width:           100,
		NoiseCount:      0,
		ShowLineOptions: 2 | 4,
		Length:          4,
		Source:          "1234567890qwertyuioplkjhgfdsazxcvbnm",
		BgColor: &color.RGBA{
			R: 3,
			G: 102,
			B: 214,
			A: 125,
		},
		Fonts: []string{"3Dumb.ttf"},
	}

	var driver base64Captcha.Driver = DriverString.ConvertFonts()
	c := base64Captcha.NewCaptcha(driver, store)
	id, b64s, answer, err = c.Generate() //生成图片的base64编码   图片对应的答案字符串

	return id, b64s, answer, err
}

// 验证验证码
func VerifyCaptchaLocal(id, verify_string string) bool { //verify_string是客户端传过来的字符串
	verify_pass := store.Verify(id, verify_string, true) //是否验证通过
	return verify_pass
}

// 新增远程验证码微服务
var (
	service = "captcha"
	version = "latest"
)

func GenerateCaptcha() (id, b64s, answer string, err error) {
	//配置验证码信息   https://captcha.mojotv.cn/.netlify/functions/captcha
	// Create client
	cClient := pb.NewCaptchaService(service, CaptchaMicroClient)

	// Call service
	res, err := cClient.GenerateCaptcha(context.Background(), &pb.GenerateCaptchaRequest{
		Height: 40,
		Width:  100,
		Length: 4,
	})

	return res.Id, res.B64S, res.Answer, err
}

// 验证验证码
func VerifyCaptcha(id, verify_string string) bool { //verify_string是客户端传过来的字符串
	// Create client
	cClient := pb.NewCaptchaService(service, CaptchaMicroClient)
	// Call service
	res, _ := cClient.VerifyCaptcha(context.Background(), &pb.VerifyCaptchaRequest{
		Id:           id,
		VerifyString: verify_string,
	})
	return res.VerifyPass
}
