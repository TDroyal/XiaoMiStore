package admin

import (
	"XiaoMiStore/dao"
	"XiaoMiStore/logic"
	"XiaoMiStore/models"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
)

type RoleController struct { //角色管理
	BaseController
}

func (con RoleController) GetRoleList(c *gin.Context) { //后期按需求改成分页查询
	roleList := []models.Role{}
	if err := dao.DB.Find(&roleList).Error; err != nil {
		con.Error(c, "获取角色列表失败，请稍后再试", -1, nil)
		return
	}
	con.Success(c, "获取角色列表成功", 0, roleList)
}

func (con RoleController) Add(c *gin.Context) { //添加角色
	title := strings.Trim(c.PostForm("title"), " ")             //角色名称  (去除输入字符串中的空格)
	description := strings.Trim(c.PostForm("description"), " ") //角色描述

	role := models.Role{}
	role.Title = title
	role.Description = description
	role.Status = 1

	if err := dao.DB.Create(&role).Error; err != nil {
		con.Error(c, "增加角色失败，请重试", -1, nil)
		return
	}
	con.Success(c, "增加角色成功", 0, nil)
}

func (con RoleController) GetRoleInfo(c *gin.Context) {
	//获取需要修改的角色id           //首先把角色信息传到对应的表单框中
	id := c.Query("id")

	role := models.Role{}
	role.ID = uint(logic.StringToInt(id))
	if err := dao.DB.Find(&role).Error; err != nil {
		con.Error(c, "获取角色信息失败，请稍后重试", -1, nil)
		return
	}
	con.Success(c, "获取角色信息成功", 0, role)
}

func (con RoleController) Edit(c *gin.Context) {
	id := c.PostForm("id")
	title := c.PostForm("title")
	description := c.PostForm("description")
	role := models.Role{}
	role.ID = uint(logic.StringToInt(id))

	var cnt int64
	if err := dao.DB.Find(&role).Count(&cnt).Error; err != nil {
		con.Error(c, "修改角色信息失败，请稍后重试", -1, nil)
		return
	}

	if cnt == 0 {
		con.Error(c, "此角色不存在", -1, nil)
		return
	}

	if err := dao.DB.Model(&role).Updates(map[string]interface{}{"title": title, "description": description}).Error; err != nil {
		con.Error(c, "修改角色信息失败，请稍后重试", -1, nil)
		fmt.Println(role)
		return
	}
	fmt.Println(role)
	con.Success(c, "修改角色信息成功", 0, nil)
}

func (con RoleController) Delete(c *gin.Context) {
	id := uint(logic.StringToInt(c.PostForm("id")))
	role := models.Role{}
	role.ID = id

	var cnt int64
	if err := dao.DB.Model(&models.Role{}).Where("id = ?", id).Count(&cnt).Error; err != nil {
		con.Error(c, "删除角色信息失败，请稍后重试", -1, nil)
		return
	}

	if cnt == 0 {
		con.Error(c, "此角色不存在", -1, nil)
		return
	}

	if err := dao.DB.Delete(&role).Error; err != nil {
		con.Error(c, "删除角色失败，请稍后重试", -1, nil)
		return
	}
	con.Success(c, "删除角色成功", 0, nil)
}
