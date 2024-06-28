// 轮播图管理，显示在PC端的首页

package admin

import (
	"XiaoMiStore/dao"
	"XiaoMiStore/logic"
	"XiaoMiStore/models"

	"github.com/gin-gonic/gin"
)

type FocusController struct {
	BaseController
}

// 获取轮播图列表
func (con FocusController) GetFocusList(c *gin.Context) {
	focus := []models.Focus{}
	if err := dao.DB.Find(&focus).Error; err != nil {
		con.Error(c, "获取轮播图列表失败", -1, nil)
		return
	}
	con.Success(c, "获取轮播图列表成功", 0, focus)
}

func (con FocusController) GetFocusInfo(c *gin.Context) { // 获取轮播图信息
	id := uint(logic.StringToInt(c.Query("id")))
	focus := models.Focus{}
	focus.ID = id
	if err := dao.DB.Find(&focus).Error; err != nil {
		con.Error(c, "获取轮播图信息失败", -1, nil)
		return
	}
	con.Success(c, "获取轮播图信息成功", 0, focus)
}

func (con FocusController) Add(c *gin.Context) {
	//接收前端传过来的数据，这里支持单文件（图片）上传  对上传的文件进行压缩存储，减轻服务器的带宽和存储压力
	//前端以postform形式将数据传过来
	title := c.PostForm("title")
	focus_type := logic.StringToInt(c.PostForm("focus_type"))
	link := c.PostForm("link")
	sort := logic.StringToInt(c.PostForm("sort"))
	status := logic.StringToInt(c.PostForm("status"))

	//上传单个图片  focus_img是前端传过来的图片对应的name="focus_img"
	focus_img, upload_err := logic.UploadImageFile(c, "focus_img")
	if upload_err != nil {
		con.Error(c, "上传轮播图失败，请稍后重试", -1, nil)
		return
	}

	focus := models.Focus{
		Title:     title,
		FocusType: focus_type,
		FocusImg:  focus_img,
		Link:      link,
		Sort:      sort,
		Status:    status,
	}
	if err := dao.DB.Create(&focus).Error; err != nil {
		con.Error(c, "上传轮播图失败，请稍后重试", -1, nil)
		return
	}
	con.Success(c, "添加轮播图成功", 0, nil)
}

func (con FocusController) Edit(c *gin.Context) {
	id := uint(logic.StringToInt(c.PostForm("id")))
	title := c.PostForm("title")
	focus_type := logic.StringToInt(c.PostForm("focus_type"))
	link := c.PostForm("link")
	sort := logic.StringToInt(c.PostForm("sort"))
	status := logic.StringToInt(c.PostForm("status"))

	//上传单个图片  focus_img是前端传过来的图片对应的name="focus_img"
	focus_img, upload_err := logic.UploadImageFile(c, "focus_img")
	if upload_err != nil {
		con.Error(c, "修改轮播图失败，请稍后重试", -1, nil)
		return
	}

	if focus_img == "" {
		if err := dao.DB.Model(&models.Focus{}).Where("id = ?", id).Updates(map[string]interface{}{"title": title, "focus_type": focus_type, "link": link, "sort": sort, "status": status}).Error; err != nil {
			con.Error(c, "修改轮播图失败，请稍后重试", -1, nil)
			return
		}
	} else {
		if err := dao.DB.Model(&models.Focus{}).Where("id = ?", id).Updates(map[string]interface{}{"title": title, "focus_type": focus_type, "link": link, "sort": sort, "status": status, "focus_img": focus_img}).Error; err != nil {
			con.Error(c, "修改轮播图失败，请稍后重试", -1, nil)
			return
		}
	}

	con.Success(c, "修改轮播图成功", 0, nil)
}

type FocusID struct {
	ID uint `json:"id"`
}

func (con FocusController) Delete(c *gin.Context) {
	focus_id := FocusID{}
	if err := c.ShouldBind(&focus_id); err != nil {
		con.Error(c, "删除轮播图失败，请稍后重试", -1, nil)
		return
	}

	//根据自己的需求   看看需不需要删除图片 os.Remove("static/xxx/xxx.png")

	if err := dao.DB.Where("id = ?", focus_id.ID).Delete(&models.Focus{}).Error; err != nil {
		con.Error(c, "删除轮播图失败，请稍后重试", -1, nil)
		return
	}
	con.Success(c, "删除轮播图成功", 0, nil)
}
