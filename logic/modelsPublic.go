package logic

import (
	"XiaoMiStore/dao"
	"XiaoMiStore/models"
)

// 封装模型models的一些方法

/*
根据商品分类获取推荐商品
	@param {Number} cateID - 分类ID  如何是顶级分类，还要获取子分类的所有数据
	@param {String} GoodsType - hot best new all
	@param {Number} limitNum - 数量
*/

func GetGoodsByCategory(cateID int, goodsType string, limitNum int) ([]models.Goods, error) {
	goodsList := []models.Goods{}

	//判断cate_id是否是顶级分类
	var isTopCate int64
	if err := dao.DB.Model(&models.GoodsCate{}).Where("id = ? and pid = 0", cateID).Count(&isTopCate).Error; err != nil {
		return nil, err
	}

	var tempSlice []int //存顶级分类以及下面二级分类的ID
	tempSlice = append(tempSlice, cateID)
	if isTopCate == 1 { //是顶级分类
		// 获取顶级分类下面的二级分类
		goodsCateList := []models.GoodsCate{}
		if err := dao.DB.Where("pid = ?", cateID).Find(&goodsCateList).Error; err != nil {
			return nil, err
		}

		for i := 0; i < len(goodsCateList); i++ {
			tempSlice = append(tempSlice, goodsCateList[i].ID)
		}
	}

	goods_type := "" //all
	if goodsType == "hot" {
		goods_type = "is_hot = 1"
	} else if goodsType == "best" {
		goods_type = "is_best = 1"
	} else if goodsType == "new" {
		goods_type = "is_new = 1"
	}

	if err := dao.DB.Where("cate_id = in and "+goods_type, tempSlice).Order("sort DESC").Limit(limitNum).Find(&goodsList).Error; err != nil {
		return nil, err
	}
	return goodsList, nil
}
