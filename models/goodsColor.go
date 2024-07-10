package models

//商品颜色表  (这个表的增删改查就后续再实现了，主要是不常用)
type GoodsColor struct {
	ID         int    `gorm:"primarKey"`
	ColorName  string `gorm:"varchar(255)"`
	ColorVaule string `gorm:"varchar(255)"`
	Status     int    `gorm:"tinyint(1)"`
	Checked    bool   `gorm:"-"` //表示创建数据库表时，忽略此字段，Checked判断当前颜色是否属于某商品
}

func (GoodsColor) TableName() string {
	return "goods_color"
}
