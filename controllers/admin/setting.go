package admin

import (
	"XiaoMiStore/dao"
	"XiaoMiStore/logic"
	"XiaoMiStore/models"

	"github.com/gin-gonic/gin"
)

type SettingController struct {
	BaseController
}

// 获取系统的设置信息
func (con SettingController) GetSettingInfo(c *gin.Context) {
	setting := models.Setting{}
	if err := dao.DB.First(&setting).Error; err != nil {
		con.Error(c, "获取系统的设置信息失败", -1, nil)
		return
	}
	con.Success(c, "获取系统的设置信息成功", 0, setting)
}

// 编辑系统设置信息
func (con SettingController) Edit(c *gin.Context) {
	setting := models.Setting{ID: 1}
	if err := dao.DB.Find(&setting).Error; err != nil {
		con.Error(c, "修改系统的设置信息失败", -1, nil)
		return
	}

	// file图片不支持shouldBind解析，不解析图片
	if err := c.ShouldBind(&setting); err != nil {
		con.Error(c, "修改系统的设置信息失败", -1, nil)
		return
	}

	// 上传图片 logo
	siteLogo, upload_err1 := logic.UploadImageFile(c, "site_logo")
	if upload_err1 != nil {
		con.Error(c, "修改系统的设置信息失败", -1, nil)
		return
	}

	if siteLogo != "" { //为空表示没有上传图片，不要把原来的图片覆盖了
		setting.SiteLogo = siteLogo
	}

	// 上传默认图片 no_picture
	noPicture, upload_err2 := logic.UploadImageFile(c, "no_picture")
	if upload_err2 != nil {
		con.Error(c, "修改系统的设置信息失败", -1, nil)
		return
	}

	if noPicture != "" {
		setting.NoPicture = noPicture
	}

	// 执行修改
	if err := dao.DB.Save(&setting).Error; err != nil {
		con.Error(c, "修改系统的设置信息失败", -1, nil)
		return
	}

	con.Success(c, "修改系统的设置信息成功", 0, nil)
}
