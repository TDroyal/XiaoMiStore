package logic

import (
	"encoding/json"
	"fmt"

	"github.com/gin-gonic/gin"
)

// 定义结构体  缓存结构体  私有
type ginCookie struct{}

// 写入数据的方法  存储在用户浏览器的
func (cookie ginCookie) Set(c *gin.Context, key string, value any) error {
	if v, err := json.Marshal(value); err != nil {
		return err
	} else {
		// c.SetCookie(key, string(v), 3600, "/", "localhost", false, true) // "127.0.0.1"
		// des加密
		desKey := []byte("12345678") // des加解密算法的密钥，8位
		v, _ = DesEncrypt(v, desKey)
		c.SetCookie(key, string(v), 3600, "/", c.Request.Host, false, true) // "127.0.0.1"
		return nil
	}
}

// 获取数据的方法
func (cookie ginCookie) Get(c *gin.Context, key string, obj any) bool {
	v, err := c.Cookie(key)
	if err == nil && v != "" && v != "[]" {
		// des解密
		desKey := []byte("12345678")
		dv, e := DesDecrypt([]byte(v), desKey)
		if e != nil {
			return false
		}
		if err1 := json.Unmarshal(dv, &obj); err1 != nil {
			fmt.Println(err1)
			return false
		}
		return true
	}
	return false
}

// 删除数据的方法
func (cookie ginCookie) Remove(c *gin.Context, key string) bool {
	c.SetCookie(key, "", -1, "/", c.Request.Host, false, true) // "127.0.0.1"
	return true
}

// 实例化结构体
var Cookie = &ginCookie{}
