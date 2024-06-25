package routers

import (
	"XiaoMiStore/controllers/admin"
	"XiaoMiStore/middlewares"

	"github.com/gin-gonic/gin"
)

func SetupAdminRouters(r *gin.Engine) {
	adminRouters := r.Group("/admin", middlewares.InitAdminAuthMiddleware)
	{
		//管理员登录
		adminRouters.GET("/index", admin.LoginController{}.Index)                      //模拟后台管理首页
		adminRouters.GET("/generateCaptcha", admin.LoginController{}.GenerateACaptcha) //生成图形验证码
		adminRouters.POST("/doLogin", admin.LoginController{}.Login)                   //后台管理员登录路由
		adminRouters.GET("/doLogout", admin.LoginController{}.Logout)                  //后台管理员退出登录路由

		//轮播图管理
		adminRouters.POST("focus/add", admin.FocusController{}.Add)       //添加轮播图
		adminRouters.POST("focus/edit", admin.FocusController{}.Edit)     //编辑轮播图
		adminRouters.POST("focus/delete", admin.FocusController{}.Delete) //删除轮播图

		//角色管理
		adminRouters.GET("role/getRoleList", admin.RoleController{}.GetRoleList) //获取角色列表信息
		adminRouters.GET("role/getRoleInfo", admin.RoleController{}.GetRoleInfo) //获取角色信息
		adminRouters.POST("role/add", admin.RoleController{}.Add)                //添加角色
		adminRouters.POST("role/edit", admin.RoleController{}.Edit)              //编辑角色
		adminRouters.POST("role/delete", admin.RoleController{}.Delete)          //删除角色

		//管理员管理
		adminRouters.GET("manager/getManagerList", admin.ManagerController{}.GetManagerList) //获取管理员列表
		adminRouters.GET("manager/getManagerInfo", admin.ManagerController{}.GetManagerInfo) //获取角色信息
		adminRouters.POST("manager/add", admin.ManagerController{}.Add)                      //管理员添加
		adminRouters.POST("manager/edit", admin.ManagerController{}.Edit)                    //管理员编辑
		adminRouters.POST("manager/delete", admin.ManagerController{}.Delete)                //管理员删除
	}
}
