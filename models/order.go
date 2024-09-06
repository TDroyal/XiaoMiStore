package models

import "time"

type Order struct {
	ID          int
	UID         int
	OrderID     string //订单的一个订单号  zz4864865gdaxyszz
	AllPrice    float64
	Name        string // 收货人姓名
	Phone       string
	Address     string
	PayStatus   int // 支付状态     0未支付    1已支付
	PayType     int // 支付类型：0  alipay支付宝，1  wechat微信
	OrderStatus int // 订单状态  0已下单  1已付款  2已配货  3已发货  4交易成功  5退货  6取消
	CreatedAt   time.Time
}

func (o Order) TableName() string {
	return "order"
}
