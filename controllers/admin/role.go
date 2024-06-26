package admin

import (
	"XiaoMiStore/dao"
	"XiaoMiStore/logic"
	"XiaoMiStore/models"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
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

func (con RoleController) GetAuthInfo(c *gin.Context) {
	role_id := c.Query("id") //获取角色id
	//根据角色id获取当前角色已有的授权信息，并对复选框进行自行选中
	authInfo := []models.RoleAccess{}
	if err := dao.DB.Where("role_id = ?", role_id).Find(&authInfo).Error; err != nil {
		con.Error(c, "获取授权信息失败", -1, nil)
		return
	}
	con.Success(c, "获取授权信息成功", 0, authInfo)
}

func (con RoleController) Auth(c *gin.Context) {
	//获取前端传过来的角色id，以及需要对齐进行授权的所有授权id
	//假设前端以form表单传过来，授权id是checkbox形式
	role_id := logic.StringToInt(c.PostForm("role_id")) //角色id
	accessIDList := c.PostFormArray("access_node[]")    //[1 2 3 5]

	//下面两步必须用事务一起做
	if err := dao.DB.Transaction(func(tx *gorm.DB) error { //自动事务
		// 在事务中执行一些 DB 操作（从这里开始，您应该使用 'tx' 而不是 'DB'）
		//1. 把与role_id相关的旧的授权数据全部清空
		if err := dao.DB.Where("role_id = ?", role_id).Delete(&models.RoleAccess{}).Error; err != nil {
			con.Error(c, "授权失败", -1, nil)
			return err
		}
		//2. 再增加新的授权数据，
		roleAccessAdd := models.RoleAccess{}
		for _, v := range accessIDList {
			roleAccessAdd.RoleID = role_id
			roleAccessAdd.AccessID = logic.StringToInt(v)
			if err := dao.DB.Create(&roleAccessAdd).Error; err != nil {
				con.Error(c, "授权失败", -1, nil)
				return err
			}
		}

		// 返回 nil 提交事务
		return nil
	}); err != nil {
		con.Error(c, "授权失败", -1, nil)
		return
	}
	con.Success(c, "授权成功", 0, nil)
}
