package mistore

import (
	"XiaoMiStore/dao"
	"XiaoMiStore/models"
	"fmt"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// 管理用户收货地址
type AddressController struct {
	BaseController
}

// 增加收货地址
func (con AddressController) AddAddress(c *gin.Context) {
	// 1. 获取用户信息  以及  前端传过来的表单数据
	user, ok := c.MustGet("user").(models.User) // 类型断言
	if !ok {
		con.Error(c, "获取用户信息失败", -1, nil)
		return
	}
	fmt.Println("=======用户信息======\n", user) // 可以在cookie中获取，也可以在中间件set的map中获取
	/*
		=======用户信息======
		{1 13508385027 e10adc3949ba59abbe56e057f20f883e 2024-09-02 14:52:28 +0800 CST 127.0.0.1  1}
	*/

	name := c.PostForm("name")
	phone := c.PostForm("phone")
	address := c.PostForm("address")

	// 2. 判断收获地址的数量
	var addressNum int64
	if err := dao.DB.Model(&models.Address{}).Where("uid = ?", user.ID).Count(&addressNum).Error; err != nil {
		con.Error(c, "添加地址失败", -1, nil)
		return
	}
	if addressNum > 10 {
		con.Error(c, "地址数已达上限", -1, nil)
		return
	}

	// 3. 更新当前用户的所有收获地址的默认收获地址状态为0
	dao.DB.Model(&models.Address{}).Where("uid = ? and default_address = ?", user.ID, 1).Update("default_address", 0)

	// 4. 增加当前收货地址，让默认收获地址状态是1
	newAddress := models.Address{
		UID:            user.ID,
		Name:           name,
		Phone:          phone,
		Address:        address,
		DefaultAddress: 1, // 默认地址
	}
	if err := dao.DB.Create(&newAddress).Error; err != nil {
		con.Error(c, "添加地址失败", -1, nil)
		return
	}
	con.Success(c, "添加新地址成功", 0, nil)
}

// 编辑收货地址
func (con AddressController) EditAddress(c *gin.Context) {
	// 1. 获取用户信息
	user, ok := c.MustGet("user").(models.User) // 类型断言
	if !ok {
		con.Error(c, "获取用户信息失败", -1, nil)
		return
	}

	// 2. 获取表单信息
	id := c.PostForm("id")
	name := c.PostForm("name")
	phone := c.PostForm("phone")
	address := c.PostForm("address")

	if err := dao.DB.Model(&models.Address{}).Where("id = ? and uid = ?", id, user.ID).Updates(map[string]any{"name": name, "phone": phone, "address": address}).Error; err != nil {
		con.Error(c, "更新地址信息失败", -1, nil)
		return
	}
	con.Success(c, "更新地址信息成功", 0, nil)
}

// 获取一个收获地址
func (con AddressController) GetOneAddress(c *gin.Context) {
	user, ok := c.MustGet("user").(models.User) // 类型断言
	if !ok {
		con.Error(c, "获取用户信息失败", -1, nil)
		return
	}

	id := c.Query("id")
	address := models.Address{}
	if err := dao.DB.Where("id = ? and uid = ?", id, user.ID).Find(&address).Error; err != nil {
		con.Error(c, "获取地址信息失败", -1, nil)
		return
	}
	con.Success(c, "获取地址信息成功", 0, gin.H{
		"address": address,
	})
}

// 获取用户所有的收获地址
func (con AddressController) GetAllAddress(c *gin.Context) {
	// 1. 获取用户信息
	user, ok := c.MustGet("user").(models.User) // 类型断言
	if !ok {
		con.Error(c, "获取用户信息失败", -1, nil)
		return
	}
	addressList := []models.Address{}
	if err := dao.DB.Where("uid = ?", user.ID).Find(&addressList).Error; err != nil {
		con.Error(c, "获取地址列表信息失败", -1, nil)
		return
	}
	con.Success(c, "获取地址列表信息成功", 0, gin.H{
		"addressList": addressList,
	})
}

// 点击切换默认收获地址
func (con AddressController) ChangeDefaultAddress(c *gin.Context) {
	// 1. 获取用户信息
	user, ok := c.MustGet("user").(models.User) // 类型断言
	if !ok {
		con.Error(c, "获取用户信息失败", -1, nil)
		return
	}
	// 2. 获取收获地址id
	id := c.PostForm("id")

	// 3-4 要一起做，事务
	dao.DB.Transaction(func(tx *gorm.DB) error {
		// 在事务中执行一些 DB 操作（从这里开始，您应该使用 'tx' 而不是 'DB'）

		// 3. 将之前的默认收获地址状态设置为0
		if err := tx.Model(&models.Address{}).Where("uid = ? and default_address = 1", user.ID).Update("default_address", 0).Error; err != nil {
			// 返回任何错误都会回滚事务
			con.Error(c, "切换默认地址失败", -1, nil)
			return err
		}
		// 4. 将当前点击的收获地址设为默认地址
		if err := tx.Model(&models.Address{}).Where("id = ? and uid = ?", id, user.ID).Update("default_address", 1).Error; err != nil {
			// 返回任何错误都会回滚事务
			con.Error(c, "切换默认地址失败", -1, nil)
			return err
		}

		// 返回 nil 提交事务
		return nil
	})
	con.Success(c, "成功切换默认收货地址", 0, nil)
}
