package models

import (
	"gorm.io/gorm"
)

//轮播图表

type Focus struct {
	gorm.Model
	Title     string `gorm:"type:varchar(255)"` //轮播图的标题
	FocusType int    `gorm:"type:tinyint(1)"`   //类型表示此轮播图在pc端显示，还是移动端显示   1表示网站，2表示app，3表示小程序
	FocusImg  string `gorm:"type:varchar(255)"` //轮播图的图片的url
	Link      string `gorm:"type:varchar(255)"` //用户点击轮播图的跳转地址
	Sort      int    `gorm:"type:int"`
	Status    int    `gorm:"type:tinyint(1)"`
}

// 修改默认表名为user
func (u Focus) TableName() string {
	return "focus"
}
