package mistore

import (
	"XiaoMiStore/dao"
	"XiaoMiStore/logic"
	"XiaoMiStore/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type BuyController struct {
	BaseController
}

// 测试中间件
func (con BuyController) TestBuy(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "测试：用户支付成功",
		"status":  "0",
		"data":    nil,
	})
}

// 结算界面 获取需要结算的商品列表
func (con BuyController) GetCheckedCartList(c *gin.Context) {
	cartList := []models.Cart{}
	if ok := logic.Cookie.Get(c, "cartList", &cartList); !ok {
		con.Error(c, "获取商品列表失败", -1, nil)
		return
	}
	orderList := []models.Cart{}
	var allPrice float64
	var allNum int
	for _, v := range cartList {
		if v.Checked {
			allPrice += v.Price * float64(v.Num)
			allNum += v.Num
			orderList = append(orderList, v)
		}
	}
	con.Success(c, "获取商品列表成功", 0, gin.H{
		"allPrice":  allPrice,
		"allNum":    allNum,
		"orderList": orderList,
	})
}

/*
提交订单执行结算
	1.获取用户信息以及收货地址信息
	2.获取购买商品信息
	3.把订单信息放在订单表，把商品信息放在商品表
	4.删除购物车里面的选中数据

	5.提交订单成功后，前端跳转到支付页面去
	否则前端提示提交订单失败，不去支付页面
*/

// 思考：结算界面提交订单  如何防止订单重复提交？(如何防止用户提交订单成功后，回到上一个页面再次点击提交订单)
/*
	提交订单的时候需要判断：用户需要上传一个签名(这个签名是进入（刷新）结算界面时生成的)，这个签名也是保存在session中的  提交成功之后，  把session中保存的签名删掉即可
	当结算界面的商品列表为空，直接退到首页
*/

// 提交订单  执行结算
func (con BuyController) DoCheckout(c *gin.Context) {
	// 1. 获取用户信息
	user, ok := c.MustGet("user").(models.User) // 类型断言
	if !ok {
		con.Error(c, "获取用户信息失败", -1, nil)
		return
	}

	// 1.获取收获地址
	address := models.Address{UID: user.ID, DefaultAddress: 1}
	if err := dao.DB.Find(&address).Error; err != nil {
		con.Error(c, "获取收获地址失败", -1, nil)
		return
	}

	// 2.获取购买商品信息 (获取购物车中选中的商品)
	cartList := []models.Cart{}
	if ok := logic.Cookie.Get(c, "cartList", &cartList); !ok {
		con.Error(c, "获取商品列表失败", -1, nil)
		return
	}
	orderList := []models.Cart{}
	var allPrice float64
	for _, v := range cartList {
		if v.Checked {
			allPrice += v.Price * float64(v.Num)
			orderList = append(orderList, v)
		}
	}

	// 3.把订单信息放在订单表，把商品信息放在商品表
	order := models.Order{
		UID:         user.ID,
		OrderID:     logic.GetOrderID(),
		AllPrice:    allPrice,
		Name:        address.Name,
		Phone:       address.Phone,
		Address:     address.Address,
		PayStatus:   0, // 未支付
		PayType:     0, // alipay
		OrderStatus: 0, // 已下单
	}

	// 3.把订单信息放在订单表，把商品信息放在商品表 （用事务）
	if err := dao.DB.Transaction(func(tx *gorm.DB) error {
		// 在事务中执行一些 db 操作（从这里开始，您应该使用 'tx' 而不是 'db'）
		if err := tx.Create(&order).Error; err != nil {
			con.Error(c, "创建订单失败", -1, nil)
			return err // 返回任何错误都会回滚事务
		}

		for _, item := range orderList { // 购物车中需要结算的商品
			orderItem := models.OrderItem{
				OrderID:      order.ID,
				UID:          user.ID,
				ProductTitle: item.Title,
				ProductID:    item.ID, // 购物车ID 也就是存储的商品ID
				ProductImage: item.GoodsImage,
				ProductPrice: item.Price,
				ProductNum:   item.Num,
				GoodsVersion: item.GoodsVersion,
				GoodsColor:   item.GoodsColor,
			}
			if err := dao.DB.Create(&orderItem).Error; err != nil {
				con.Error(c, "创建订单失败", -1, nil)
				return err // 返回任何错误都会回滚事务
			}
		}

		// 返回 nil 提交事务
		return nil
	}); err != nil {
		return
	}

	// 4.删除购物车里面的选中数据
	noSelectedCartList := []models.Cart{}
	for _, item := range cartList {
		if !item.Checked {
			noSelectedCartList = append(noSelectedCartList, item)
		}
	}
	logic.Cookie.Set(c, "cartList", &noSelectedCartList)

	con.Success(c, "创建订单成功", 0, nil)
}

// 支付页面，需要获取订单相关的信息
func (con BuyController) GetOrderInfo(c *gin.Context) {
	order_id := c.Query("order_id")

	// 1. 获取用户信息
	user, ok := c.MustGet("user").(models.User) // 类型断言
	if !ok {
		con.Error(c, "获取用户信息失败", -1, nil)
		return
	}

	order := models.Order{}
	dao.DB.Where("id = ?", order_id).Find(&order)
	if user.ID != order.UID {
		con.Error(c, "非法访问", -1, nil)
		return
	}

	// 获取订单对应的商品
	orderItems := []models.OrderItem{}
	dao.DB.Where("order_id = ?", order.ID).Find(&orderItems)

	con.Success(c, "获取订单信息成功", 0, gin.H{
		"order":      order,
		"orderItems": orderItems,
	})
}
