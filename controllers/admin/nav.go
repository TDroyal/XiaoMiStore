package admin

import (
	"XiaoMiStore/dao"
	"XiaoMiStore/logic"
	"XiaoMiStore/models"
	"math"

	"github.com/gin-gonic/gin"
)

// 导航栏的增删改查
type NavController struct {
	BaseController
}

// 获取导航栏列表  (分页处理)
func (con NavController) GetNavList(c *gin.Context) {

	page := logic.StringToInt(c.DefaultQuery("page", "1"))
	if page == 0 {
		page = 1
	}
	pageSize := 5

	nav := []models.Nav{}
	if err := dao.DB.Offset((page - 1) * pageSize).Limit(pageSize).Find(&nav).Error; err != nil {
		con.Error(c, "获取导航栏列表失败", -1, nil)
		return
	}

	// 获取商品总数量
	var totalCount int64
	if err := dao.DB.Model(&models.Nav{}).Count(&totalCount).Error; err != nil {
		con.Success(c, "获取导航栏列表失败", -1, nil)
		return
	}

	con.Success(c, "获取导航栏列表成功", 0, gin.H{
		"nav":        nav,
		"totalCount": totalCount,
		"page":       page,
		"totalPage":  math.Ceil(float64(totalCount) / float64(pageSize)),
	})
}

func (con NavController) GetNavInfo(c *gin.Context) { // 获取导航栏信息
	id := logic.StringToInt(c.Query("id"))
	nav := models.Nav{}
	nav.ID = id
	if err := dao.DB.Find(&nav).Error; err != nil {
		con.Error(c, "获取导航栏信息失败", -1, nil)
		return
	}
	con.Success(c, "获取导航栏信息成功", 0, nav)
}

func (con NavController) Add(c *gin.Context) {
	//接收前端传过来的数据，这里支持单文件（图片）上传  对上传的文件进行压缩存储，减轻服务器的带宽和存储压力
	//前端以postform形式将数据传过来
	title := c.PostForm("title")
	link := c.PostForm("link")
	position := logic.StringToInt(c.PostForm("position"))
	isOpennew := logic.StringToInt(c.PostForm("is_opennew"))
	relation := c.PostForm("relation")
	sort := logic.StringToInt(c.PostForm("sort"))
	status := logic.StringToInt(c.PostForm("status"))

	nav := models.Nav{
		Title:     title,
		Link:      link,
		Position:  position,
		IsOpennew: isOpennew,
		Relation:  relation,
		Sort:      sort,
		Status:    status,
	}
	if err := dao.DB.Create(&nav).Error; err != nil {
		con.Error(c, "添加导航栏失败，请稍后重试", -1, nil)
		return
	}
	con.Success(c, "添加导航栏成功", 0, nil)
}

func (con NavController) Edit(c *gin.Context) {
	id := logic.StringToInt(c.PostForm("id"))
	title := c.PostForm("title")
	link := c.PostForm("link")
	position := logic.StringToInt(c.PostForm("position"))
	isOpennew := logic.StringToInt(c.PostForm("is_opennew"))
	relation := c.PostForm("relation")
	sort := logic.StringToInt(c.PostForm("sort"))
	status := logic.StringToInt(c.PostForm("status"))

	if err := dao.DB.Model(&models.Nav{}).Where("id = ?", id).Updates(map[string]interface{}{"title": title, "position": position, "link": link, "sort": sort, "status": status, "is_opennew": isOpennew, "relation": relation}).Error; err != nil {
		con.Error(c, "修改导航栏失败，请稍后重试", -1, nil)
		return
	}

	con.Success(c, "修改导航栏成功", 0, nil)
}

type NavID struct {
	ID int `json:"id"`
}

func (con NavController) Delete(c *gin.Context) {
	nav_id := NavID{}
	if err := c.ShouldBind(&nav_id); err != nil {
		con.Error(c, "删除导航栏失败，请稍后重试", -1, nil)
		return
	}
	if err := dao.DB.Where("id = ?", nav_id.ID).Delete(&models.Nav{}).Error; err != nil {
		con.Error(c, "删除导航栏失败，请稍后重试", -1, nil)
		return
	}
	con.Success(c, "删除导航栏成功", 0, nil)
}
