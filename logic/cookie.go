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
		c.SetCookie(key, string(v), 3600, "/", "localhost", false, true) // "127.0.0.1"
		return nil
	}
}

// 获取数据的方法
func (cookie ginCookie) Get(c *gin.Context, key string, obj any) bool {
	v, err := c.Cookie(key)
	if err == nil && v != "" && v != "[]" {
		if err1 := json.Unmarshal([]byte(v), &obj); err1 != nil {
			fmt.Println(err1)
			return false
		}
		return true
	}
	return false
}

// 实例化结构体
var Cookie = &ginCookie{}
