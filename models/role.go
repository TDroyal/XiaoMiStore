package models

import "gorm.io/gorm"

// 角色表（每个管理员对应的角色）
type Role struct {
	gorm.Model
	Title       string
	Description string
	Status      int
}

func (Role) TableName() string {
	return "role"
}
