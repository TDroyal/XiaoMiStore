package models

import "time"

type OrderItem struct {
	ID           int
	UID          int
	OrderID      int // 对应的是order表中的ID字段
	ProductTitle string
	ProductID    int
	ProductImage string
	ProductPrice float64
	ProductNum   int
	GoodsVersion string
	GoodsColor   string
	CreatedAt    time.Time
}

func (o OrderItem) TableName() string {
	return "order_item"
}
