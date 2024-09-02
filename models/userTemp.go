package models

import "time"

type UserTemp struct {
	ID        int
	Ip        string
	Phone     string
	SendCount int    // 每日发送短信次数（不能超过n，否则不发送验证码）
	CreateDay string // 20240902 用于后续统计 每天一个用户注册发送短信的次数
	CreatedAt time.Time
	Sign      string // 签名，从一个页面跳转到另一个页面时候的一个签名验证
}

func (u UserTemp) TableName() string {
	return "user_temp"
}
