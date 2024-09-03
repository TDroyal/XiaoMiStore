package models

import "time"

// 收获地址表
type Address struct {
	ID             int
	UID            int
	Name           string
	Phone          string
	Address        string
	DefaultAddress int // 1表示是默认地址 0表示不是默认地址
	CreatedAt      time.Time
}

func (a Address) TableName() string {
	return "address"
}
