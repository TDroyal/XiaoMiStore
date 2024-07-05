package admin

import (
	"XiaoMiStore/dao"
	"XiaoMiStore/logic"
	"XiaoMiStore/models"

	"github.com/gin-gonic/gin"
)

type GoodsTypeAttributeController struct {
	BaseController
}

func (con GoodsTypeAttributeController) GetGoodsTypeAttributeList(c *gin.Context) { //获取根据商品类型ID获取其对应的所有属性列表
	cate_id := logic.StringToInt(c.Query("id")) //商品类型ID
	goodsTypeAttributeList := []models.GoodsTypeAttribute{}
	if err := dao.DB.Where("cate_id = ?", cate_id).Find(&goodsTypeAttributeList).Error; err != nil {
		con.Error(c, "获取商品类型属性列表失败", -1, nil)
		return
	}
	con.Error(c, "获取商品类型属性列表成功", 0, goodsTypeAttributeList)
}

func (con GoodsTypeAttributeController) GetGoodsTypeAttributeInfo(c *gin.Context) { //获取某个商品类型属性的信息
	id := logic.StringToInt(c.Query("id")) //商品类型属性ID
	goodsTypeAttribute := models.GoodsTypeAttribute{}
	if err := dao.DB.Where("id = ?", id).Find(&goodsTypeAttribute).Error; err != nil {
		con.Error(c, "获取商品类型属性信息失败", -1, nil)
		return
	}
	con.Error(c, "获取商品类型属性信息成功", 0, goodsTypeAttribute)
}

type GoodsTypeAttrInfo struct {
	Title     string `json:"title" form:"title"`
	CateID    int    `json:"cate_id" form:"cate_id"` //所属商品类型ID  (GoodsType)
	AttrType  int    `json:"attr_type" form:"attr_type"`
	AttrValue string `json:"attr_value" form:"attr_value"`
	Sort      int    `json:"sort" form:"sort"`
	// Status    int    `json:"status" form:"status"`  //插入元素默认status为1
}

func (con GoodsTypeAttributeController) Add(c *gin.Context) { //添加商品类型属性
	goodsTypeAttrInfo := GoodsTypeAttrInfo{}
	if err := c.ShouldBind(&goodsTypeAttrInfo); err != nil {
		con.Error(c, "添加商品类型属性失败", -1, nil)
		return
	}
	goodsTypeAttribute := models.GoodsTypeAttribute{
		Title:     goodsTypeAttrInfo.Title,
		CateID:    goodsTypeAttrInfo.CateID,
		AttrType:  goodsTypeAttrInfo.AttrType,
		AttrValue: goodsTypeAttrInfo.AttrValue,
		Sort:      goodsTypeAttrInfo.Sort,
		Status:    1,
	}
	if err := dao.DB.Create(&goodsTypeAttribute).Error; err != nil {
		con.Error(c, "添加商品类型属性失败", -1, nil)
		return
	}
	con.Success(c, "添加商品类型属性成功", 0, nil)
}

type EditGoodsTypeAttrInfo struct {
	ID        int    `json:"id" form:"id"`
	Title     string `json:"title" form:"title"`
	CateID    int    `json:"cate_id" form:"cate_id"` //所属商品类型ID  (GoodsType)
	AttrType  int    `json:"attr_type" form:"attr_type"`
	AttrValue string `json:"attr_value" form:"attr_value"`
	Sort      int    `json:"sort" form:"sort"`
	Status    int    `json:"status" form:"status"` //插入元素默认status为1
}

func (con GoodsTypeAttributeController) Edit(c *gin.Context) {
	editGoodsTypeAttrInfo := EditGoodsTypeAttrInfo{}
	if err := c.ShouldBind(&editGoodsTypeAttrInfo); err != nil {
		con.Error(c, "修改商品类型属性失败", -1, nil)
		return
	}
	goodsTypeAttribute := models.GoodsTypeAttribute{
		ID:        editGoodsTypeAttrInfo.ID,
		CateID:    editGoodsTypeAttrInfo.CateID,
		Title:     editGoodsTypeAttrInfo.Title,
		AttrType:  editGoodsTypeAttrInfo.AttrType,
		AttrValue: editGoodsTypeAttrInfo.AttrValue,
		Status:    editGoodsTypeAttrInfo.Status,
		Sort:      editGoodsTypeAttrInfo.Sort,
	}
	if err := dao.DB.Save(&goodsTypeAttribute).Error; err != nil {
		con.Error(c, "修改商品类型属性信息失败，请稍后重试", -1, nil)
		return
	}
	con.Success(c, "修改商品类型属性信息成功", 0, nil)

}

type GoodsTypeAttrID struct {
	ID int `json:"id"`
}

func (con GoodsTypeAttributeController) Delete(c *gin.Context) {
	//获取商品类型属性ID
	goodsTypeAttrID := GoodsTypeAttrID{}
	if err := c.ShouldBind(&goodsTypeAttrID); err != nil {
		con.Error(c, "删除商品类型属性失败，请稍后重试", -1, nil)
		return
	}
	goodsTypeAttribute := models.GoodsTypeAttribute{ID: goodsTypeAttrID.ID}

	if err := dao.DB.Delete(&goodsTypeAttribute).Error; err != nil {
		con.Error(c, "删除商品类型属性失败，请稍后重试", -1, nil)
		return
	}
	con.Success(c, "删除商品类型属性成功", 0, nil)

}
