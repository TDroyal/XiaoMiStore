package logic

import (
	"crypto/md5"
	"fmt"
)

// 获得字符串的md5值
func GetMD5(str string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(str)))
}
