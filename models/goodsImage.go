package models

import "time"

//保存商品相册信息
type GoodsImage struct {
	ID        int `gorm:"primaryKey"`
	GoodsID   int `gorm:"index"` //对应goods表的ID
	ImgUrl    string
	ColorID   int `gorm:"index"` //对应goods_color表，和颜色进行关联，点击不同的颜色，切换不同颜色的商品图片
	Sort      int
	Status    int
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (GoodsImage) TableName() string {
	return "goods_image"
}
