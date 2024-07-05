package models

import "time"

//商品类型属性
type GoodsTypeAttribute struct {
	ID        int    `gorm:"primaryKey"`
	CateID    int    `gorm:"index"` //对应goods_type表的ID，
	Title     string //属性名称
	AttrType  int    //属性类型代表该属性后面的录入方式  1表示单行文本框，2表示多行文本框，3表示下拉列表
	AttrValue string //对应3下拉列表（一行代表一个可选择）,录入下拉列表的选项(可选值列表)
	Status    int
	Sort      int
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (GoodsTypeAttribute) TableName() string {
	return "goods_type_attribute"
}
