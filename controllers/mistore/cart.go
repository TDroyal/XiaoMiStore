package mistore

import (
	"XiaoMiStore/dao"
	"XiaoMiStore/logic"
	"XiaoMiStore/models"

	"github.com/gin-gonic/gin"
)

// 购物车
type CartController struct {
	BaseController
}

// 获取购物车数据
func (con CartController) GetCartData(c *gin.Context) {
	var cartList []models.Cart
	logic.Cookie.Get(c, "cartList", &cartList)

	var allPrice float64
	for i := 0; i < len(cartList); i++ {
		if cartList[i].Checked {
			allPrice += cartList[i].Price * float64(cartList[i].Num)
		}
	}

	con.Success(c, "获取购物车信息成功", 0, gin.H{
		"cartList": cartList,
		"allPrice": allPrice,
	})
}

/*
购物车数据保持到哪里？

	1：购物车数据保存到本地     （cookie）

	2：购物车数据保存到服务器（mysql）  （必须登录）

	3：没有登录 购物车数据保存到本地，登录成功后购物车数据保存到数据库
*/
func (con CartController) AddCart(c *gin.Context) { // 将商品加入购物车
	/*
		增加购物车的实现逻辑：
			1.获取增加购物车的数据（把哪个商品加入到购物车）
			2.判断购物车有没有数据（cookie）
			3.如果购物车没有任何数据 直接把当前数据写入cookie
			4.如果购物车有数据
				1.判断购物车有没有当前数据
					有当前数据让当前数据的数量+1，然后写入到cookie
				2.如果没有当前数据直接写入cookie
	*/

	// 点击加入购物车 获取商品的 goodsID 以及 colorID  （假设商品只有这两个属性，其它的如法炮制）
	// 1.获取增加购物车的数据，放在结构体里面（把哪个商品加入到购物车）
	goodsID := logic.StringToInt(c.PostForm("goods_id"))
	colorID := logic.StringToInt(c.PostForm("color_id"))

	goods := models.Goods{}
	goodsColor := models.GoodsColor{}
	if err := dao.DB.Where("id = ?", goodsID).Find(&goods).Error; err != nil {
		con.Error(c, "添加购物车失败", -1, nil)
		return
	}
	if err := dao.DB.Where("id = ?", colorID).Find(&goodsColor).Error; err != nil {
		con.Error(c, "添加购物车失败", -1, nil)
		return
	}

	currentData := models.Cart{
		ID:           goodsID,
		Title:        goods.Title,
		Price:        goods.Price,
		GoodsVersion: goods.GoodsVersion,
		Num:          1,
		GoodsColor:   goodsColor.ColorName,
		GoodsImage:   goods.GoodsImg,
		GoodsGift:    goods.GoodsGift, // 格式：1,23,24  商品id
		GoodsAttr:    "",
		Checked:      true, //默认选中
	}

	// 2.判断购物车有没有数据（cookie）
	cartList := []models.Cart{}
	logic.Cookie.Get(c, "cartList", &cartList)
	if len(cartList) > 0 { // 4. 购物车有数据
		// 4.1.判断购物车有没有当前数据
		if (models.Cart{}).HasCartData(cartList, currentData) {
			// 有当前数据让当前数据的数量+1，然后写入到cookie
			for i := len(cartList) - 1; i >= 0; i-- {
				if cartList[i].ID == currentData.ID && cartList[i].GoodsColor == currentData.GoodsColor {
					cartList[i].Num++
					break
				}
			}
		} else { // 4.2.如果没有当前数据直接写入cookie
			cartList = append(cartList, currentData)
		}
		logic.Cookie.Set(c, "cartList", &cartList)
	} else { // 3.如果购物车没有任何数据 直接把当前数据写入cookie
		cartList = append(cartList, currentData)
		_ = logic.Cookie.Set(c, "cartList", &cartList)
	}
	con.Success(c, "商品加入购物车成功", 0, gin.H{
		"cartlist": cartList,
	})
}

/*
	前端在购物车列表页面点击某个商品的-或+按钮触发下面的两个方法
*/

// 增加购物车数量
func (con CartController) IncCart(c *gin.Context) {
	// 获取商品ID以及商品颜色
	goodsID := logic.StringToInt(c.PostForm("goods_id"))
	goodsColor := c.PostForm("goods_color")

	var number int
	var currentPrice float64
	var allPrice float64

	cartList := []models.Cart{}
	logic.Cookie.Get(c, "cartList", &cartList)
	for i := len(cartList) - 1; i >= 0; i-- {
		if cartList[i].ID == goodsID && cartList[i].GoodsColor == goodsColor {
			cartList[i].Num += 1
			number = cartList[i].Num
			currentPrice = float64(number) * cartList[i].Price
		}
		if cartList[i].Checked {
			allPrice += cartList[i].Price * float64(cartList[i].Num)
		}
	}
	logic.Cookie.Set(c, "cartList", &cartList)

	con.Success(c, "更新购物车成功", 0, gin.H{
		"allPrice":     allPrice,
		"num":          number,       // 当前这个商品的数量
		"currentPrice": currentPrice, // 当前这个商品的总价
	})
}

// 减少购物车数量
func (con CartController) DecCart(c *gin.Context) {
	// 获取商品ID以及商品颜色
	goodsID := logic.StringToInt(c.PostForm("goods_id"))
	goodsColor := c.PostForm("goods_color")

	var number int
	var currentPrice float64
	var allPrice float64

	cartList := []models.Cart{}
	logic.Cookie.Get(c, "cartList", &cartList)
	for i := len(cartList) - 1; i >= 0; i-- {
		if cartList[i].ID == goodsID && cartList[i].GoodsColor == goodsColor {
			if cartList[i].Num > 1 {
				cartList[i].Num -= 1
			}
			number = cartList[i].Num
			currentPrice = float64(number) * cartList[i].Price
		}
		if cartList[i].Checked {
			allPrice += cartList[i].Price * float64(cartList[i].Num)
		}
	}
	logic.Cookie.Set(c, "cartList", &cartList)

	con.Success(c, "更新购物车成功", 0, gin.H{
		"allPrice":     allPrice,
		"num":          number,       // 当前这个商品的数量
		"currentPrice": currentPrice, // 当前这个商品的总价
	})

}

/*
	前端在购物车列表页面点击某个商品的选中按钮触发下面方法
*/

// 改变一个商品的选中状态
func (con CartController) ChangeOneCartCheckedStatus(c *gin.Context) {
	// 获取商品ID以及商品颜色
	goodsID := logic.StringToInt(c.PostForm("goods_id"))
	goodsColor := c.PostForm("goods_color")

	var allPrice float64

	cartList := []models.Cart{}
	logic.Cookie.Get(c, "cartList", &cartList)
	for i := len(cartList) - 1; i >= 0; i-- {
		if cartList[i].ID == goodsID && cartList[i].GoodsColor == goodsColor {
			cartList[i].Checked = !cartList[i].Checked
		}
		if cartList[i].Checked {
			allPrice += cartList[i].Price * float64(cartList[i].Num)
		}
	}
	logic.Cookie.Set(c, "cartList", &cartList)

	con.Success(c, "改变商品的选中状态成功", 0, gin.H{
		"allPrice": allPrice,
	})
}

// 全选反选
func (con CartController) ChangeAllCartCheckedStatus(c *gin.Context) {
	flag := logic.StringToInt(c.PostForm("checkedAll")) // 1表示全选，2表示全不选
	cartList := []models.Cart{}
	logic.Cookie.Get(c, "cartList", &cartList)

	var allPrice float64
	if flag == 1 {
		for i := len(cartList) - 1; i >= 0; i-- {
			cartList[i].Checked = true
			allPrice += float64(cartList[i].Num) * cartList[i].Price
		}
	} else if flag == 2 {
		for i := len(cartList) - 1; i >= 0; i-- {
			cartList[i].Checked = false
		}
	}

	logic.Cookie.Set(c, "cartList", &cartList)
	con.Success(c, "改变商品的选中状态成功", 0, gin.H{
		"allPrice": allPrice,
	})
}

// 删除购物车数据
func (con CartController) DelOneCart(c *gin.Context) {
	// 获取商品ID以及商品颜色
	goodsID := logic.StringToInt(c.PostForm("goods_id"))
	goodsColor := c.PostForm("goods_color")

	var allPrice float64
	cartList := []models.Cart{}
	logic.Cookie.Get(c, "cartList", &cartList)
	var idx int = -1
	for i := len(cartList) - 1; i >= 0; i-- {
		if cartList[i].ID == goodsID && cartList[i].GoodsColor == goodsColor {
			// 删除i这个数据
			idx = i
		}
		if cartList[i].Checked && idx != i {
			allPrice += cartList[i].Price * float64(cartList[i].Num)
		}
	}
	cartList = append(cartList[0:idx], cartList[idx+1:]...) // 删除i这个数据
	logic.Cookie.Set(c, "cartList", &cartList)

	con.Success(c, "改变商品的选中状态成功", 0, gin.H{
		"allPrice": allPrice,
	})
}
