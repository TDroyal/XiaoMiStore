package models

import "time"

//角色表role和权限表access多对多
type RoleAccess struct {
	RoleID    int `gorm:"primaryKey"`
	AccessID  int `gorm:"primaryKey"`
	CreatedAt time.Time
}

func (RoleAccess) TableName() string {
	return "role_access"
}
