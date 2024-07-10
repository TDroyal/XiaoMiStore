package models

import "time"

// 存放物品的规格包装
type GoodsAttr struct {
	ID              int    `gorm:"primaryKey"`
	GoodsID         int    `gorm:"index"` //对应goods表的ID
	AttributeCateID int    `gorm:"index"` //对应goods_type表的ID
	AttributeID     int    `gorm:"index"` //对应goods_type_attribute表的ID
	AttributeTitle  string // 存的和goods_type_attribute表中的title一样的东西
	AttributeType   int    //对应goods_type_attribute表的attr_type  (用于判断生成文本框还是下拉框)
	AttributeValue  string //每个属性对应具体的值
	Status          int
	Sort            int
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

func (GoodsAttr) TableName() string {
	return "goods_attr"
}
