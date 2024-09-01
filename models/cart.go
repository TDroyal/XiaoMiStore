package models

// 购物车表
type Cart struct {
	ID           int
	Title        string
	Price        float64
	GoodsVersion string
	UID          int // 用户id
	Num          int
	GoodsGift    string
	GoodsFitting string
	GoodsColor   string
	GoodsImage   string
	GoodsAttr    string
	Checked      bool `gorm:"-"` //加入购物车，默认选择，即要购买
}

func (c Cart) TableName() string {
	return "cart"
}

// 判断购物车列表cartList中有没有当前数据currentData
func (c Cart) HasCartData(cartList []Cart, currentData Cart) bool {
	for i := len(cartList) - 1; i >= 0; i-- {
		if cartList[i].ID == currentData.ID && cartList[i].GoodsColor == currentData.GoodsColor {
			return true
		}
	}
	return false
}
