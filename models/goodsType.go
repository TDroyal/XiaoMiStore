package models

import "time"

//商品类型  （相关操作类似Role表 角色管理）
type GoodsType struct {
	ID          int    `gorm:"primaryKey"`
	Title       string `gorm:"type:varchar(255);"`
	Description string `gorm:"type:varchar(255)"`
	Status      int    `gorm:"type:tinyint(1)"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (GoodsType) TableName() string {
	return "goods_type"
}
