package routers

import (
	"XiaoMiStore/controllers/mistore"

	"github.com/gin-gonic/gin"
)

func SetupDefaultRouters(r *gin.Engine) {
	defaultRouters := r.Group("/")
	{
		defaultRouters.GET("/", mistore.DefaultController{}.Index)
		defaultRouters.GET("/getCategoryListAndGoodsList", mistore.ProductController{}.GetCategoryListAndGoodsList) // 用户点击首页左侧商品分类获取对应的商品分类列表以及商品列表
		defaultRouters.GET("/getGoodsDetailInfo", mistore.ProductController{}.GetGoodsDetailInfo)                   // 获取商品详情信息
		defaultRouters.GET("/getImageList", mistore.ProductController{}.GetImageList)                               // 根据商品id和颜色id获取商品图库
		defaultRouters.POST("/markdown_to_html", mistore.ProductController{}.MarkDownToHTML)                        // 将markdown文本转为html文本的接口（写着测试用的）
	}
}
