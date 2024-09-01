package routers

import (
	"XiaoMiStore/controllers/mistore"

	"github.com/gin-gonic/gin"
)

func SetupDefaultRouters(r *gin.Engine) {
	defaultRouters := r.Group("/")
	{
		defaultRouters.GET("/", mistore.DefaultController{}.Index)
		defaultRouters.GET("/testCookie", mistore.DefaultController{}.TestCookie)                                   // 测试封装的cookie操作
		defaultRouters.GET("/getCategoryListAndGoodsList", mistore.ProductController{}.GetCategoryListAndGoodsList) // 用户点击首页左侧商品分类获取对应的商品分类列表以及商品列表
		defaultRouters.GET("/getGoodsDetailInfo", mistore.ProductController{}.GetGoodsDetailInfo)                   // 获取商品详情信息
		defaultRouters.GET("/getImageList", mistore.ProductController{}.GetImageList)                               // 根据商品id和颜色id获取商品图库
		defaultRouters.POST("/markdown_to_html", mistore.ProductController{}.MarkDownToHTML)                        // 将markdown文本转为html文本的接口（写着测试用的）

		// 购物车目前都是在cookie中实现的，用户未登录
		defaultRouters.GET("/cart/getCartData", mistore.CartController{}.GetCartData)                                // 获取购物车数据
		defaultRouters.POST("/cart/addCart", mistore.CartController{}.AddCart)                                       // 添加购物车
		defaultRouters.POST("/cart/incCart", mistore.CartController{}.IncCart)                                       //增加购物车中商品的数量
		defaultRouters.POST("/cart/decCart", mistore.CartController{}.DecCart)                                       // 减少购物车中商品的数量
		defaultRouters.POST("/cart/changeOneCartCheckedStatus", mistore.CartController{}.ChangeOneCartCheckedStatus) // 改变购物车列表中一个商品的选中状态
		defaultRouters.POST("/cart/changeAllCartCheckedStatus", mistore.CartController{}.ChangeAllCartCheckedStatus) // 改变购物车列表中所有商品的选中状态（全选反选）
		defaultRouters.POST("/cart/delOneCart", mistore.CartController{}.DelOneCart)                                 // 删除购物车中的某条数据
	}
}
