package mistore

import (
	"XiaoMiStore/dao"
	"XiaoMiStore/logic"
	"XiaoMiStore/models"
	"math"
	"strings"

	"github.com/gin-gonic/gin"
)

type ProductController struct {
	BaseController
}

// 点击首页左侧的商品分类   获取对应的（商品分类列表、商品列表）数据、商品列表分页（点击）
// 根据分类id获取商品列表数据，如果是顶级分类获取所有的二级分类；如果是二级分类把所有的兄弟分类获取了
func (con ProductController) GetCategoryListAndGoodsList(c *gin.Context) {
	cate_id := logic.StringToInt(c.Query("id")) //分类ID    //前端往后端发出GET请求，带上了?id=xxx

	page := logic.StringToInt(c.DefaultQuery("page", "1"))
	if page == 0 {
		page = 1
	}
	pageSize := 5
	// 获取当前的分类
	currentCate := models.GoodsCate{}
	if err := dao.DB.Where("id = ?", cate_id).Find(&currentCate).Error; err != nil {
		con.Error(c, "获取分类列表失败", -1, nil)
		return
	}

	// 所有相关的分类
	subCate := []models.GoodsCate{}
	var tempSlice []int       // 存相关二级分类的id
	if currentCate.Pid == 0 { //当前点击的是顶级分类
		// 获取其所有的二级分类
		if err := dao.DB.Where("pid = ?", currentCate.ID).Find(subCate).Error; err != nil { //本身就是引用类型，不需要Find(&subCate)
			con.Error(c, "获取分类列表失败", -1, nil)
			return
		}
		for i := len(subCate) - 1; i >= 0; i-- {
			tempSlice = append(tempSlice, subCate[i].ID)
		}
	} else { // 点击的是二级分类
		// 获取其所有的兄弟分类（也是同一个顶级分类下的二级分类）
		if err := dao.DB.Where("pid = ?", currentCate.Pid).Find(subCate).Error; err != nil {
			con.Error(c, "获取分类列表失败", -1, nil)
			return
		}
	}
	tempSlice = append(tempSlice, cate_id)
	goodsList := []models.Goods{}
	if err := dao.DB.Where("cate_id in ?", tempSlice).Offset((page - 1) * pageSize).Limit(pageSize).Find(goodsList).Error; err != nil {
		con.Error(c, "获取商品列表失败", -1, nil)
		return
	}
	// 获取商品总数量
	var totalCount int64
	if err := dao.DB.Model(&models.Goods{}).Where("cate_id in ?", tempSlice).Count(&totalCount).Error; err != nil {
		con.Error(c, "获取商品列表信息失败", -1, nil)
		return
	}

	// 拿给前端渲染即可
	con.Success(c, "获取商品分类以及对应的商品列表信息成功", 0, gin.H{
		"page":        page,
		"totalCount":  totalCount,
		"totalPage":   math.Ceil(float64(totalCount) / float64(pageSize)),
		"currentCate": currentCate,
		"subCate":     subCate,
		"goodsList":   goodsList,
	})
}

// 获取商品的详情信息   // 前端点击某个商品，进入详情页面，对应的路由为 xxx/detail?id=19
func (con ProductController) GetGoodsDetailInfo(c *gin.Context) {
	id := logic.StringToInt(c.Query("id"))
	// 1. 获取商品信息
	goods := models.Goods{ID: id}
	if err := dao.DB.Find(&goods).Error; err != nil {
		con.Error(c, "获取商品详情信息失败", -1, nil)
		return
	}

	// 2. 获取关联商品 RelationGoods      string  格式：1,23,24  表示这些商品和当前商品关联
	// relationGoods := []models.Goods{}
	var relationGoods []models.Goods
	// var relationGoodsID []int
	relationGoodsIDString := strings.Split(goods.RelationGoods, ",")
	// for i := len(relationGoodsIDString) - 1; i >= 0; i-- {
	// 	relationGoodsID = append(relationGoodsID, logic.StringToInt(relationGoodsIDString[i]))
	// }
	if err := dao.DB.Model(&models.Goods{}).Where("id in ?", relationGoodsIDString).Select("id, title, price, goods_version").Find(&relationGoods).Error; err != nil {
		con.Error(c, "获取商品详情信息失败", -1, nil)
		return
	}

	// 3. 获取关联赠品 GoodsGift    格式：1,23,24  商品id
	goodsGift := []models.Goods{}
	goodsGiftID := strings.Split(goods.GoodsGift, ",")
	if err := dao.DB.Model(&models.Goods{}).Where("id in ?", goodsGiftID).Select("id, title, price, goods_version").Find(&goodsGift).Error; err != nil {
		con.Error(c, "获取商品详情信息失败", -1, nil)
		return
	}

	// 4. 获取关联颜色 GoodsColor   数据库中存储格式也是：1,2,3
	goodsColor := []models.GoodsColor{}
	goodsColorID := strings.Split(goods.GoodsColor, ",")
	if err := dao.DB.Model(&models.GoodsColor{}).Where("id in ?", goodsColorID).Find(&goodsColor).Error; err != nil {
		con.Error(c, "获取商品详情信息失败", -1, nil)
		return
	}

	// 5. 获取关联配件 GoodsFitting  格式：1,23,24  商品id
	goodsFitting := []models.Goods{}
	goodsFittingID := strings.Split(goods.GoodsFitting, ",")
	if err := dao.DB.Model(&models.Goods{}).Where("id in ?", goodsFittingID).Select("id, title, price, goods_version").Find(&goodsFitting).Error; err != nil {
		con.Error(c, "获取商品详情信息失败", -1, nil)
		return
	}

	// 6. 获取商品关联的图片  GoodsImage
	goodsImage := []models.GoodsImage{}
	if err := dao.DB.Where("goods_id = ?", goods.ID).Limit(6).Find(&goodsImage).Error; err != nil {
		con.Error(c, "获取商品详情信息失败", -1, nil)
		return
	}

	// 7. 获取规格参数信息 GoodsAttr
	goodsAttr := []models.GoodsAttr{}
	if err := dao.DB.Where("goods_id = ?", goods.ID).Find(&goodsAttr).Error; err != nil {
		con.Error(c, "获取商品详情信息失败", -1, nil)
		return
	}

	con.Success(c, "获取商品详情信息成功", 0, gin.H{
		"goods":         goods,
		"relationGoods": relationGoods,
		"goodsGift":     goodsGift,
		"goodsColor":    goodsColor,
		"goodsFitting":  goodsFitting,
		"goodsImage":    goodsImage,
		"goodsAttr":     goodsAttr,
	})
}

// 根据商品id和颜色id获取商品图库
func (con ProductController) GetImageList(c *gin.Context) {
	goodsID := logic.StringToInt(c.Query("goods_id"))
	colorID := logic.StringToInt(c.Query("color_id"))

	// 查询商品图库信息
	goodsImage := models.GoodsImage{}
	if err := dao.DB.Model(&models.GoodsImage{}).Where("goods_id = ? and color_id = ?", goodsID, colorID).Find(&goodsImage); err != nil {
		con.Error(c, "获取商品图库信息失败", -1, nil)
		return
	}
	con.Success(c, "获取商品图库信息成功", 0, goodsImage)
}

func (con ProductController) MarkDownToHTML(c *gin.Context) {
	md := c.PostForm("markdown_text")
	html_string := logic.FormatAttr(md)
	con.Success(c, "转换成功", 0, html_string)
}
