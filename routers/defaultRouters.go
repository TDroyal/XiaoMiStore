package routers

import (
	"XiaoMiStore/controllers/mistore"
	"XiaoMiStore/middlewares"

	"github.com/gin-gonic/gin"
)

func SetupDefaultRouters(r *gin.Engine) {
	defaultRouters := r.Group("/default")
	{
		defaultRouters.GET("/", mistore.DefaultController{}.Index)
		defaultRouters.GET("/testCookie", mistore.DefaultController{}.TestCookie)                                   // 测试封装的cookie操作
		defaultRouters.GET("/getCategoryListAndGoodsList", mistore.ProductController{}.GetCategoryListAndGoodsList) // 用户点击首页左侧商品分类获取对应的商品分类列表以及商品列表
		defaultRouters.GET("/getGoodsDetailInfo", mistore.ProductController{}.GetGoodsDetailInfo)                   // 获取商品详情信息
		defaultRouters.GET("/getImageList", mistore.ProductController{}.GetImageList)                               // 根据商品id和颜色id获取商品图库
		defaultRouters.POST("/markdown_to_html", mistore.ProductController{}.MarkDownToHTML)                        // 将markdown文本转为html文本的接口（写着测试用的）

		// 测试调用goods微服务接口
		defaultRouters.POST("/addgoods", mistore.DefaultController{}.AddGoods)

		// 购物车目前都是在cookie中实现的，用户未登录
		defaultRouters.GET("/cart/getCartData", mistore.CartController{}.GetCartData)                                // 获取购物车数据
		defaultRouters.POST("/cart/addCart", mistore.CartController{}.AddCart)                                       // 添加购物车
		defaultRouters.POST("/cart/incCart", mistore.CartController{}.IncCart)                                       //增加购物车中商品的数量
		defaultRouters.POST("/cart/decCart", mistore.CartController{}.DecCart)                                       // 减少购物车中商品的数量
		defaultRouters.POST("/cart/changeOneCartCheckedStatus", mistore.CartController{}.ChangeOneCartCheckedStatus) // 改变购物车列表中一个商品的选中状态
		defaultRouters.POST("/cart/changeAllCartCheckedStatus", mistore.CartController{}.ChangeAllCartCheckedStatus) // 改变购物车列表中所有商品的选中状态（全选反选）
		defaultRouters.POST("/cart/delOneCart", mistore.CartController{}.DelOneCart)                                 // 删除购物车中的某条数据

		// 用户登录/注册的相关路由
		defaultRouters.GET("/pass/generateCaptcha", mistore.PassController{}.GenerateACaptcha)         //生成图形验证码
		defaultRouters.POST("/pass/sendCode", mistore.PassController{}.SendCode)                       // 向用户手机号发送验证码
		defaultRouters.POST("/pass/vertifyStep1ToStep2", mistore.PassController{}.VertifyStep1ToStep2) // 注册步骤一跳转到步骤二需要进行验证（防止用户直接打开注册的步骤二页面）
		defaultRouters.POST("/pass/vertifySmsCode", mistore.PassController{}.VertifySmsCode)           // 验证短信验证码是否正确
		defaultRouters.POST("/pass/vertifyStep2ToStep3", mistore.PassController{}.VertifyStep2ToStep3) // 注册步骤二跳转到步骤三需要进行验证（防止用户直接打开注册的步骤三页面）
		defaultRouters.POST("/pass/doRegister", mistore.PassController{}.DoRegister)                   // 输入密码完成注册
		defaultRouters.POST("/pass/login", mistore.PassController{}.Login)                             // 用户登录

		// 结算价钱的路由需要用户登录了才能执行
		defaultRouters.POST("/buy/test", middlewares.InitUserAuthMiddleware, mistore.BuyController{}.TestBuy)                         // 测试中间件 （登录的用户才会执行到后面的handler去）
		defaultRouters.GET("/buy/getCheckedCartList", middlewares.InitUserAuthMiddleware, mistore.BuyController{}.GetCheckedCartList) // 进入结算界面  获取需要结算的商品列表
		defaultRouters.POST("/buy/doCheckout", middlewares.InitUserAuthMiddleware, mistore.BuyController{}.DoCheckout)                // 在结算界面  提交订单
		defaultRouters.GET("/buy/getOrderInfo", middlewares.InitUserAuthMiddleware, mistore.BuyController{}.GetOrderInfo)             // 支付页面，获取订单相关的信息

		// 操作收货地址相关路由
		defaultRouters.GET("/address/getOneAddress", middlewares.InitUserAuthMiddleware, mistore.AddressController{}.GetOneAddress)                // 获取一个收货地址
		defaultRouters.GET("/address/getAllAddress", middlewares.InitUserAuthMiddleware, mistore.AddressController{}.GetAllAddress)                // 获取所有收货地址
		defaultRouters.POST("/address/add", middlewares.InitUserAuthMiddleware, mistore.AddressController{}.AddAddress)                            // 增加收货地址
		defaultRouters.POST("/address/edit", middlewares.InitUserAuthMiddleware, mistore.AddressController{}.EditAddress)                          // 编辑收货地址
		defaultRouters.POST("/address/changeDefaultAddress", middlewares.InitUserAuthMiddleware, mistore.AddressController{}.ChangeDefaultAddress) // 切换默认收获地址
	}
}
