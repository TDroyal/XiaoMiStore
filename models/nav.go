package models

import "time"

// 商城首页的各种导航
type Nav struct {
	ID        int    `gorm:"primaryKey"`
	Title     string //导航的标题
	Link      string //导航的链接地址
	Position  int    //导航的位置  1表示是一个最顶部的导航  2表示是一个中间导航  3表示是一个底部导航
	IsOpennew int    //当前这个连接地址是否在新窗口打开 1为否  2为是
	Relation  string //配置中部导航的关联商品，将关联的商品id放在这里  1,3,5,7
	Sort      int
	Status    int
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (Nav) TableName() string {
	return "nav"
}
