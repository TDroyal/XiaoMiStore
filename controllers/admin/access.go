package admin

import (
	"XiaoMiStore/dao"
	"XiaoMiStore/logic"
	"XiaoMiStore/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AccessController struct { //权限管理
	BaseController
}

// 不同身份的管理员登录成功后台后，根据管理员角色左侧显示不同的导航栏
// 同时权限列表也可以用这个接口获取的信息进行渲染
func (con AccessController) GetAccessList(c *gin.Context) {
	//获取顶级模块(一级导航栏，同时加载出来它对应的二级导航type=2栏以及操作权限type=3)
	topModuleList := []models.Access{}
	if err := dao.DB.Where("module_id = ?", 0).Preload("AccessList", func(db *gorm.DB) *gorm.DB {
		return db.Order("access.sort DESC")
	}).Order("sort desc").Find(&topModuleList).Error; err != nil {
		con.Error(c, "获取模块信息失败", -1, nil)
		return
	}
	con.Success(c, "获取模块信息成功", 0, topModuleList)
}

func (con AccessController) GetTopModule(c *gin.Context) { // 获取一级模块列表(做下拉框)
	//获取顶级模块
	topModuleList := []models.Access{}
	if err := dao.DB.Where("module_id = ?", 0).Find(&topModuleList).Error; err != nil {
		con.Error(c, "获取顶级模块信息失败", -1, nil)
		return
	}
	con.Success(c, "获取顶级模块信息成功", 0, topModuleList)
}

func (con AccessController) GetAccessInfo(c *gin.Context) { //获取权限信息
	//前端发出的GET请求携带了querystring:   /xxx/edit?id=55
	id := logic.StringToInt(c.Query("id"))
	accessInfo := models.Access{ID: id}
	if err := dao.DB.Find(&accessInfo).Error; err != nil {
		con.Error(c, "获取权限信息失败", -1, nil)
		return
	}
	con.Success(c, "获取权限信息成功", 0, accessInfo)
}

// 新的权限信息
type AccessInfo struct {
	ModuleName  string `json:"module_name"`
	Type        int    `json:"type"`
	ActionName  string `json:"action_name"`
	Url         string `json:"url"`
	ModuleID    int    `json:"module_id"`
	Sort        int    `json:"sort"`
	Description string `json:"description"`
	Status      int    `json:"status"`
}

func (con AccessController) Add(c *gin.Context) { //添加权限
	//获取前端以Json形式传过来的需要新添加权限信息
	accessInfo := AccessInfo{}
	if err := c.ShouldBind(&accessInfo); err != nil {
		con.Error(c, "添加权限失败，请稍后重试", -1, nil)
		return
	}
	access := models.Access{
		ModuleName:  accessInfo.ModuleName,
		Type:        accessInfo.Type,
		ActionName:  accessInfo.ActionName,
		Url:         accessInfo.Url,
		ModuleID:    accessInfo.ModuleID,
		Sort:        accessInfo.Sort,
		Description: accessInfo.Description,
		Status:      accessInfo.Status,
	}
	if err := dao.DB.Create(&access).Error; err != nil {
		con.Error(c, "添加权限失败，请稍后重试", -1, nil)
		return
	}
	con.Success(c, "添加权限成功", 0, nil)
}

// 前端传过来的需要修改的权限信息
type EditAccessInfo struct {
	ID          int    `json:"id"`
	ModuleName  string `json:"module_name"`
	Type        int    `json:"type"`
	ActionName  string `json:"action_name"`
	Url         string `json:"url"`
	ModuleID    int    `json:"module_id"`
	Sort        int    `json:"sort"`
	Description string `json:"description"`
	Status      int    `json:"status"`
}

func (con AccessController) Edit(c *gin.Context) {
	//一些具体的细节具体应用时再补全即可，例如哪些数据必须得到binding:"required"等
	accessInfo := EditAccessInfo{}
	if err := c.ShouldBind(&accessInfo); err != nil {
		con.Error(c, "修改权限信息失败，请稍后重试", -1, nil)
		return
	}
	access := models.Access{ID: accessInfo.ID}

	// 不需要查找，性能降低了
	// if err := dao.DB.Find(&access).Error; err != nil {
	// 	con.Error(c, "修改权限信息失败，请稍后重试", -1, nil)
	// 	return
	// }
	access.ModuleName = accessInfo.ModuleName
	access.Type = accessInfo.Type
	access.ActionName = accessInfo.ActionName
	access.Url = accessInfo.Url
	access.ModuleID = accessInfo.ModuleID
	access.Sort = accessInfo.Sort
	access.Description = accessInfo.Description
	access.Status = accessInfo.Status

	if err := dao.DB.Save(&access).Error; err != nil {
		con.Error(c, "修改权限信息失败，请稍后重试", -1, nil)
		return
	}
	con.Success(c, "修改权限信息成功", 0, nil)
}

type AccessID struct {
	ID int `json:"id"`
}

func (con AccessController) Delete(c *gin.Context) { //前端以JSON格式将id传到后端
	//需要判断当前删除的模块是否是顶级模块(一级列表)
	//如果当前删除一级列表，下面如果有二级列表或者操作，就不能删除它，否则可以删除它
	accessID := AccessID{}
	if err := c.ShouldBind(&accessID); err != nil {
		con.Error(c, "删除权限失败，请稍后重试", -1, nil)
		return
	}
	access := models.Access{ID: accessID.ID}
	if err := dao.DB.Find(&access).Error; err != nil {
		con.Error(c, "删除权限失败，请稍后重试", -1, nil)
		return
	}

	if access.ModuleID == 0 { //是一级列表
		accessList := []models.Access{}
		if err := dao.DB.Where("module_id = ?", access.ID).Find(&accessList).Error; err != nil {
			con.Error(c, "删除权限失败，请稍后重试", -1, nil)
			return
		}
		if len(accessList) > 0 { //当前一级列表下还有二级列表或者操作
			con.Error(c, "删除失败，此权限下还有子菜单或子操作", -1, nil)
			return
		}
		if err := dao.DB.Delete(&access).Error; err != nil {
			con.Error(c, "删除权限失败，请稍后重试", -1, nil)
			return
		}
	} else { //不是一级列表
		if err := dao.DB.Delete(&access).Error; err != nil {
			con.Error(c, "删除权限失败，请稍后重试", -1, nil)
			return
		}
	}
	con.Error(c, "删除权限成功", 0, nil)
}
