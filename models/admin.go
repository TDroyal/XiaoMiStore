package models

import (
	"gorm.io/gorm"
)

type Admin struct { // 管理员表   默认表名是 `admins`
	gorm.Model
	Username string `gorm:"type:varchar(255);"`
	Password string `gorm:"type:varchar(32);"`
	Mobile   string `gorm:"type:varchar(11);"`
	Email    string `gorm:"type:varchar(255);"`
	Status   int8   `gorm:"type:tinyint(1)"`
	RoleID   uint
	IsSuper  int8 `gorm:"type:tinyint(1)"`
}

// 修改默认表名为user
func (u Admin) TableName() string {
	return "admin"
}
