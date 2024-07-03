package admin

import (
	"XiaoMiStore/dao"
	"XiaoMiStore/logic"
	"XiaoMiStore/models"

	"github.com/gin-gonic/gin"
)

type GoodsTypeController struct {
	BaseController
}

func (con GoodsTypeController) GetGoodsTypeList(c *gin.Context) { //后期按需求改成分页查询
	goodsTypeList := []models.GoodsType{}
	if err := dao.DB.Find(&goodsTypeList).Error; err != nil {
		con.Error(c, "获取商品类型失败，请稍后再试", -1, nil)
		return
	}
	con.Success(c, "获取商品类型列表成功", 0, goodsTypeList)
}

func (con GoodsTypeController) GetGoodsTypeInfo(c *gin.Context) {
	//获取需要修改的角色id           //首先把角色信息传到对应的表单框中
	id := c.Query("id")
	goodsType := models.GoodsType{}
	goodsType.ID = logic.StringToInt(id)
	if err := dao.DB.Find(&goodsType).Error; err != nil {
		con.Error(c, "获取商品类型信息失败，请稍后重试", -1, nil)
		return
	}
	con.Success(c, "获取商品类型信息成功", 0, goodsType)
}

type GoodsTypeInfo struct {
	Title       string `form:"title" json:"title"`
	Description string `form:"description" json:"description"`
	Status      int    `form:"status" json:"status"`
}

func (con GoodsTypeController) Add(c *gin.Context) {
	goodsTypeInfo := GoodsTypeInfo{}
	if err := c.ShouldBind(&goodsTypeInfo); err != nil {
		con.Error(c, "添加商品类型失败，请稍后重试", -1, nil)
		return
	}
	goodsType := models.GoodsType{
		Title:       goodsTypeInfo.Title,
		Description: goodsTypeInfo.Description,
		Status:      goodsTypeInfo.Status,
	}
	if err := dao.DB.Create(&goodsType).Error; err != nil {
		con.Error(c, "添加商品类型失败，请重试", -1, nil)
		return
	}
	con.Success(c, "添加商品类型成功", 0, nil)
}

type EditGoodsTypeInfo struct {
	ID          int    `form:"id" json:"id"`
	Title       string `form:"title" json:"title"`
	Description string `form:"description" json:"description"`
	Status      int    `form:"status" json:"status"`
}

func (con GoodsTypeController) Edit(c *gin.Context) {
	goodsTypeInfo := EditGoodsTypeInfo{}
	if err := c.ShouldBind(&goodsTypeInfo); err != nil {
		con.Error(c, "修改商品类型信息失败，请稍后重试", -1, nil)
		return
	}

	goodsType := models.GoodsType{
		ID:          goodsTypeInfo.ID,
		Title:       goodsTypeInfo.Title,
		Description: goodsTypeInfo.Description,
		Status:      goodsTypeInfo.Status,
	}

	if err := dao.DB.Save(&goodsType).Error; err != nil {
		con.Error(c, "修改商品类型信息失败，请稍后重试", -1, nil)
		return
	}
	con.Success(c, "修改商品类型信息成功", 0, nil)

}

type GoodsTypeID struct {
	ID int `json:"id"`
}

func (con GoodsTypeController) Delete(c *gin.Context) {
	goodsTypeID := GoodsTypeID{}
	if err := c.ShouldBind(&goodsTypeID); err != nil {
		con.Error(c, "删除商品类型失败，请稍后重试", -1, nil)
		return
	}

	goodsType := models.GoodsType{ID: goodsTypeID.ID}

	if err := dao.DB.Delete(&goodsType).Error; err != nil {
		con.Error(c, "删除商品类型失败，请稍后重试", -1, nil)
		return
	}
	con.Success(c, "删除商品类型成功", 0, nil)
}
