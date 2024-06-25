package admin

import (
	"XiaoMiStore/dao"
	"XiaoMiStore/logic"
	"XiaoMiStore/models"
	"fmt"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type ManagerController struct { //管理员管理
	BaseController
}

// 这个是获取了管理员列表的所有信息，传送给前端太多非必要的信息
// func (con ManagerController) GetManagerList(c *gin.Context) {
// 	managerList := []models.Admin{}
// 	if err := dao.DB.Preload("Role").Find(&managerList).Error; err != nil {
// 		con.Error(c, "获取管理员列表失败，请稍后再试", -1, nil)
// 		return
// 	}
// 	// SELECT * FROM `role` WHERE `role`.`id` IN (1,2) AND `role`.`deleted_at` IS NULL
// 	// SELECT * FROM `admin` WHERE `admin`.`deleted_at` IS NULL
// 	// fmt.Printf("%#v", managerList)
// 	con.Success(c, "获取管理员列表成功", 0, managerList)
// }

// dao.DB.Preload("Role").Find(&managerList)为了不把很多管理员列表的无用信息传到前端，浪费带宽，我们考虑只传有用的数据Scan
type managerLists struct {
	ID        uint      `json:"id"`
	Username  string    `json:"username"`
	Mobile    string    `json:"mobile"`
	Email     string    `json:"email"`
	Status    int8      `json:"status"`
	Title     string    `json:"role"`
	CreatedAt time.Time `json:"create_time"`
}

// 更舒服的版本，返回前端指定的字段和数据  两表连查即可
func (con ManagerController) GetManagerList(c *gin.Context) { //后期改成分页
	managerList := []managerLists{}
	//左外连接或者内连接进行连查即可
	if err := dao.DB.Model(&models.Admin{}).Joins("Role").Select("admin.id", "username", "mobile", "email", "admin.status", "role.title", "admin.created_at").Scan(&managerList).Error; err != nil {
		con.Error(c, "获取管理员列表失败，请稍后重试", -1, nil)
		return
	}
	con.Success(c, "获取管理员列表成功", 0, managerList)
}

func (con ManagerController) GetManagerInfo(c *gin.Context) { //   admin/manager/getManagerInfo?id=5
	id := c.Query("id")
	manager := []models.Admin{}
	if err := dao.DB.Preload("Role").Where("id = ?", uint(logic.StringToInt(id))).Find(&manager).Error; err != nil {
		con.Error(c, "获取管理员信息失败，请稍后重试", -1, nil)
		return
	}
	if len(manager) == 0 {
		con.Error(c, "此管理员不存在", -1, nil)
		return
	}

	con.Success(c, "获取管理员信息成功", 0, manager)
}

func (con ManagerController) Add(c *gin.Context) {
	//可以直接用参数绑定c.ShouldBind()来实现
	username := strings.Trim(c.PostForm("username"), " ")
	password := strings.Trim(c.PostForm("password"), " ") //密码不要空格
	mobile := strings.Trim(c.PostForm("mobile"), " ")
	email := strings.Trim(c.PostForm("email"), " ")
	role_id := uint(logic.StringToInt(c.PostForm("role_id")))

	//添加前需要判断管理员username是否存在
	managerList := []models.Admin{}
	if err := dao.DB.Where("username = ?", username).Find(&managerList).Error; err != nil {
		con.Error(c, "添加管理员失败，请稍后重试", -1, nil)
		return
	}
	if len(managerList) > 0 {
		con.Error(c, "此管理员账号已经在", -1, nil)
		return
	}
	fmt.Println(managerList)
	manager := models.Admin{
		Username: username,
		Password: logic.GetMD5(password),
		Mobile:   mobile,
		Email:    email,
		Status:   1,
		RoleID:   role_id,
		IsSuper:  1,
	}
	if err := dao.DB.Create(&manager).Error; err != nil {
		con.Error(c, "添加管理员失败，请稍后重试", -1, nil)
		return
	}
	con.Success(c, "添加管理员成功", 0, nil)
}

// 前后端分离项目，前端一般传到后端的都是json格式的数据，建议采用shouldbind()
func (con ManagerController) Edit(c *gin.Context) {
	id := uint(logic.StringToInt(c.PostForm("id")))
	password := strings.Trim(c.PostForm("password"), " ") //密码不要空格
	mobile := strings.Trim(c.PostForm("mobile"), " ")
	email := strings.Trim(c.PostForm("email"), " ")
	role_id := uint(logic.StringToInt(c.PostForm("role_id")))

	manager := models.Admin{}
	var cnt int64
	if err := dao.DB.Where("id = ?", id).Find(&manager).Count(&cnt).Error; err != nil {
		con.Error(c, "修改管理员信息失败，请稍后重试", -1, nil)
		return
	}

	if cnt == 0 {
		con.Error(c, "此管理员不存在", -1, nil)
		return
	}

	manager.Mobile = mobile
	manager.Email = email
	manager.RoleID = role_id
	//判断密码是否为空，为空表示不修改密码
	if password != "" {
		manager.Password = logic.GetMD5(password)
	}
	if err := dao.DB.Save(&manager).Error; err != nil {
		con.Error(c, "修改管理员信息失败，请稍后重试", -1, nil)
		return
	}
	con.Success(c, "修改管理员信息成功", 0, nil)
}

func (con ManagerController) Delete(c *gin.Context) {
	id := uint(logic.StringToInt(c.PostForm("id")))
	var cnt int64
	if err := dao.DB.Model(&models.Admin{}).Where("id = ?", id).Count(&cnt).Error; err != nil {
		con.Error(c, "删除管理员信息失败，请稍后重试", -1, nil)
		return
	}
	if cnt == 0 {
		con.Error(c, "此管理员不存在", -1, nil)
		return
	}

	if err := dao.DB.Delete(&models.Admin{}, id).Error; err != nil {
		con.Error(c, "删除管理员失败，请稍后重试", -1, nil)
		return
	}
	con.Success(c, "删除管理员成功", 0, nil)
}
