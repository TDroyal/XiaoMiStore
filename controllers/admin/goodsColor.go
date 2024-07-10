package admin

import (
	"XiaoMiStore/dao"
	"XiaoMiStore/models"

	"github.com/gin-gonic/gin"
)

type GoodsColorController struct {
	BaseController
}

// 获取所有的颜色列表（用于前端添加商品时，实现checkbox复选框）
func (con GoodsColorController) GetGoodsColorList(c *gin.Context) {
	goodsColorList := []models.GoodsColor{}
	if err := dao.DB.Find(&goodsColorList).Error; err != nil {
		con.Error(c, "获取商品颜色列表失败", -1, nil)
		return
	}
	con.Success(c, "获取商品颜色列表成功", 0, goodsColorList)
}
