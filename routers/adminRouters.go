package routers

import (
	"XiaoMiStore/controllers/admin"
	"XiaoMiStore/middlewares"

	"github.com/gin-gonic/gin"
)

func SetupAdminRouters(r *gin.Engine) {
	adminRouters := r.Group("/admin", middlewares.InitAdminAuthMiddleware)
	{
		// 公共修改所有表中status或者数字(sort等)字段
		adminRouters.POST("/changeStatus", admin.MainController{}.ChangeStatus) //修改status
		adminRouters.POST("/changeNum", admin.MainController{}.ChangeNum)       //修改数字字段

		//管理员登录
		adminRouters.GET("/index", admin.LoginController{}.Index)                      //模拟后台管理首页
		adminRouters.GET("/generateCaptcha", admin.LoginController{}.GenerateACaptcha) //生成图形验证码
		adminRouters.POST("/doLogin", admin.LoginController{}.Login)                   //后台管理员登录路由
		adminRouters.GET("/doLogout", admin.LoginController{}.Logout)                  //后台管理员退出登录路由

		//轮播图管理
		adminRouters.GET("focus/getFocusList", admin.FocusController{}.GetFocusList) //获取轮播图列表信息
		adminRouters.GET("focus/getFocusInfo", admin.FocusController{}.GetFocusInfo) //获取轮播图信息
		adminRouters.POST("focus/add", admin.FocusController{}.Add)                  //添加轮播图
		adminRouters.POST("focus/edit", admin.FocusController{}.Edit)                //编辑轮播图
		adminRouters.POST("focus/delete", admin.FocusController{}.Delete)            //删除轮播图

		//角色管理
		adminRouters.GET("role/getRoleList", admin.RoleController{}.GetRoleList) //获取角色列表信息
		adminRouters.GET("role/getRoleInfo", admin.RoleController{}.GetRoleInfo) //获取角色信息
		adminRouters.POST("role/add", admin.RoleController{}.Add)                //添加角色
		adminRouters.POST("role/edit", admin.RoleController{}.Edit)              //编辑角色
		adminRouters.POST("role/delete", admin.RoleController{}.Delete)          //删除角色

		adminRouters.GET("role/getAuthInfo", admin.RoleController{}.GetAuthInfo) //获取角色已授权信息
		adminRouters.POST("role/auth", admin.RoleController{}.Auth)              //角色授权(利用复选框修改角色的权限：包括添加，修改，删除)

		//管理员管理
		adminRouters.GET("manager/getManagerList", admin.ManagerController{}.GetManagerList) //获取管理员列表
		adminRouters.GET("manager/getManagerInfo", admin.ManagerController{}.GetManagerInfo) //获取角色信息
		adminRouters.POST("manager/add", admin.ManagerController{}.Add)                      //管理员添加
		adminRouters.POST("manager/edit", admin.ManagerController{}.Edit)                    //管理员编辑
		adminRouters.POST("manager/delete", admin.ManagerController{}.Delete)                //管理员删除

		//权限管理
		adminRouters.GET("access/getTopModule", admin.AccessController{}.GetTopModule)   //获取一级模块列表(做下拉框)
		adminRouters.GET("access/getAccessList", admin.AccessController{}.GetAccessList) //获取权限列表
		adminRouters.GET("access/getAccessInfo", admin.AccessController{}.GetAccessInfo) //获取权限信息
		adminRouters.POST("access/add", admin.AccessController{}.Add)                    //权限添加
		adminRouters.POST("access/edit", admin.AccessController{}.Edit)                  //权限编辑
		adminRouters.POST("access/delete", admin.AccessController{}.Delete)              //权限删除

		//商品分类管理
		adminRouters.GET("goodsCate/getGoodsCateList", admin.GoodsCateController{}.GetGoodsCateList) //获取商品分类列表
		adminRouters.GET("goodsCate/getTopCate", admin.GoodsCateController{}.GetTopCate)             //获取一级商品分类列表(做下拉框)
		adminRouters.GET("goodsCate/getGoodsCateInfo", admin.GoodsCateController{}.GetGoodsCateInfo) //获取权限信息
		adminRouters.POST("goodsCate/add", admin.GoodsCateController{}.Add)                          //商品分类的添加
		adminRouters.POST("goodsCate/edit", admin.GoodsCateController{}.Edit)                        //商品分类的编辑
		adminRouters.POST("goodsCate/delete", admin.GoodsCateController{}.Delete)                    //商品分类的删除

		//商品类型管理
		adminRouters.GET("goodsType/getGoodsTypeList", admin.GoodsTypeController{}.GetGoodsTypeList) //获取商品类型列表信息
		adminRouters.GET("goodsType/getGoodsTypeInfo", admin.GoodsTypeController{}.GetGoodsTypeInfo) //获取商品类型信息
		adminRouters.POST("goodsType/add", admin.GoodsTypeController{}.Add)                          //添加商品类型
		adminRouters.POST("goodsType/edit", admin.GoodsTypeController{}.Edit)                        //编辑商品类型
		adminRouters.POST("goodsType/delete", admin.GoodsTypeController{}.Delete)                    //删除商品类型

		// 商品类型属性管理
		adminRouters.GET("goodsTypeAttribute/getGoodsTypeAttributeList", admin.GoodsTypeAttributeController{}.GetGoodsTypeAttributeList) //获取商品类型属性列表信息
		adminRouters.GET("goodsTypeAttribute/getGoodsTypeAttributeInfo", admin.GoodsTypeAttributeController{}.GetGoodsTypeAttributeInfo) //获取商品类型属性信息
		adminRouters.POST("goodsTypeAttribute/add", admin.GoodsTypeAttributeController{}.Add)                                            //添加商品类型属性
		adminRouters.POST("goodsTypeAttribute/edit", admin.GoodsTypeAttributeController{}.Edit)                                          //编辑商品类型属性
		adminRouters.POST("goodsTypeAttribute/delete", admin.GoodsTypeAttributeController{}.Delete)                                      //删除商品类型属性

		// 商品颜色管理  (增删改查需要什么再补什么)
		adminRouters.GET("goodsColor/getGoodsColorList", admin.GoodsColorController{}.GetGoodsColorList) // 获取所有的颜色列表

		//商品管理
		adminRouters.POST("goods/imageUpload", admin.GoodsController{}.ImageUpload)                     //froala富文本编辑器上传图片
		adminRouters.GET("goods/getGoodsList", admin.GoodsController{}.GetGoodsList)                    //获取商品列表信息
		adminRouters.GET("goods/getGoodsInfo", admin.GoodsCateController{}.GetGoodsCateInfo)            //获取商品的所有信息
		adminRouters.POST("goods/add", admin.GoodsController{}.Add)                                     //添加商品(goroutine有无问题？？？)
		adminRouters.POST("goods/edit", admin.GoodsController{}.Edit)                                   //修改商品(goroutine有无问题？？？)
		adminRouters.POST("goods/delete", admin.GoodsController{}.Delete)                               //删除商品
		adminRouters.POST("goods/changeGoodsImageColor", admin.GoodsController{}.ChangeGoodsImageColor) //异步修改商品的图库信息（将图片和图片颜色进行绑定）
		adminRouters.POST("goods/removeGoodsImage", admin.GoodsController{}.RemoveGoodsImage)           //异步删除商品相册信息

		// 导航栏管理
		adminRouters.GET("nav/getNavList", admin.NavController{}.GetNavList) //获取导航栏列表信息
		adminRouters.GET("nav/getNavInfo", admin.NavController{}.GetNavInfo) //获取导航栏信息
		adminRouters.POST("nav/add", admin.NavController{}.Add)              //添加导航栏
		adminRouters.POST("nav/edit", admin.NavController{}.Edit)            //编辑导航栏
		adminRouters.POST("nav/delete", admin.NavController{}.Delete)        //删除导航栏

		// 系统设置管理
		adminRouters.GET("setting/getSettingInfo", admin.SettingController{}.GetSettingInfo) //获取系统设置信息
		adminRouters.POST("setting/edit", admin.SettingController{}.Edit)                    //编辑系统设置信息
	}
}
