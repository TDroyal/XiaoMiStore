package mistore

import (
	"XiaoMiStore/logic"
	"XiaoMiStore/models"
	"net/http"

	"github.com/gin-gonic/gin"
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
