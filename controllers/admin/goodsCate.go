package admin

import (
	"XiaoMiStore/dao"
	"XiaoMiStore/logic"
	"XiaoMiStore/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type GoodsCateController struct {
	BaseController
}

func (con GoodsCateController) GetGoodsCateList(c *gin.Context) {
	//获取商品分类(同时加载出来它对应的子分类)
	goodsCateList := []models.GoodsCate{}
	if err := dao.DB.Where("pid = ?", 0).Preload("GoodsCateItems", func(db *gorm.DB) *gorm.DB {
		return db.Order("goods_cate.sort DESC")
	}).Order("sort desc").Find(&goodsCateList).Error; err != nil {
		con.Error(c, "获取商品分类列表失败", -1, nil)
		return
	}
	con.Success(c, "获取商品分类列表成功", 0, goodsCateList)
}

// 加载顶级商品分类列表
func (con GoodsCateController) GetTopCate(c *gin.Context) {
	topCateList := []models.GoodsCate{}
	if err := dao.DB.Where("pid = ?", 0).Find(&topCateList).Error; err != nil {
		con.Error(c, "获取顶级商品分类列表信息失败", -1, nil)
		return
	}
	con.Success(c, "获取顶级商品分类列表信息成功", 0, topCateList)
}

func (con GoodsCateController) GetGoodsCateInfo(c *gin.Context) { //获取商品分类信息
	//前端发出的GET请求携带了querystring:   /xxx/edit?id=55
	id := logic.StringToInt(c.Query("id"))
	goodsCateInfo := models.GoodsCate{ID: id}
	if err := dao.DB.Find(&goodsCateInfo).Error; err != nil {
		con.Error(c, "获取商品分类信息失败", -1, nil)
		return
	}
	con.Success(c, "获取商品分类信息成功", 0, goodsCateInfo)
}

// 新的商品分类信息
type GoodsCateInfo struct {
	Title string `form:"title"`
	// CateImg     string `json:"cate_img"`
	Link        string `form:"link"`
	Template    string `form:"template"`
	Pid         int    `form:"pid"`
	SubTitle    string `form:"sub_title"`
	KeyWords    string `form:"keywords"`
	Description string `form:"description"`
	Sort        int    `form:"sort"`
	Status      int    `form:"status"`
}

func (con GoodsCateController) Add(c *gin.Context) { // 添加商品分类
	//获取前端以postform形式传过来的需要新商品分类信息
	goodsCateInfo := GoodsCateInfo{}
	if err := c.ShouldBind(&goodsCateInfo); err != nil {
		con.Error(c, "添加商品分类失败，请稍后重试", -1, nil)
		return
	}
	// 再上传图片CateImg   postform形式的数据转发过来
	CateImgDir, upload_err := logic.UploadImageFile(c, "cate_img")
	if upload_err != nil {
		con.Error(c, "上传商品分类图片失败，请稍后重试", -1, nil)
		return
	}
	goodsCate := models.GoodsCate{
		Title:       goodsCateInfo.Title,
		CateImg:     CateImgDir,
		Link:        goodsCateInfo.Link,
		Template:    goodsCateInfo.Template,
		Pid:         goodsCateInfo.Pid,
		SubTitle:    goodsCateInfo.SubTitle,
		KeyWords:    goodsCateInfo.KeyWords,
		Description: goodsCateInfo.Description,
		Sort:        goodsCateInfo.Sort,
		Status:      goodsCateInfo.Status,
	}

	if err := dao.DB.Create(&goodsCate).Error; err != nil {
		con.Error(c, "添加商品分类失败，请稍后重试", -1, nil)
		return
	}
	con.Success(c, "添加商品分类成功", 0, nil)
}

// 新的商品分类信息
type EditGoodsCateInfo struct {
	ID          int    `form:"id"`
	Title       string `form:"title"`
	Link        string `form:"link"`
	Template    string `form:"template"`
	Pid         int    `form:"pid"`
	SubTitle    string `form:"sub_title"`
	KeyWords    string `form:"keywords"`
	Description string `form:"description"`
	Sort        int    `form:"sort"`
	Status      int    `form:"status"`
}

func (con GoodsCateController) Edit(c *gin.Context) {
	goodsCateInfo := EditGoodsCateInfo{}
	if err := c.ShouldBind(&goodsCateInfo); err != nil {
		con.Error(c, "修改商品分类信息失败，请稍后重试", -1, nil)
		return
	}

	//获取传过来的图片
	CateImgDir, upload_err := logic.UploadImageFile(c, "cate_img")
	if upload_err != nil {
		con.Error(c, "上传商品分类图片失败，请稍后重试", -1, nil)
		return
	}

	goodsCate := models.GoodsCate{ID: goodsCateInfo.ID}
	goodsCate.Title = goodsCateInfo.Title
	goodsCate.CateImg = CateImgDir //如果未上传图片应该默认图片不修改，参考focus管理的编辑，这里不再实现
	goodsCate.Link = goodsCateInfo.Link
	goodsCate.Template = goodsCateInfo.Template
	goodsCate.Pid = goodsCateInfo.Pid
	goodsCate.SubTitle = goodsCateInfo.SubTitle
	goodsCate.KeyWords = goodsCateInfo.KeyWords
	goodsCate.Description = goodsCateInfo.Description
	goodsCate.Sort = goodsCateInfo.Sort
	goodsCate.Status = goodsCateInfo.Status

	if err := dao.DB.Save(&goodsCate).Error; err != nil {
		con.Error(c, "修改商品分类信息失败，请稍后重试", -1, nil)
		return
	}
	con.Success(c, "修改商品分类信息成功", 0, nil)

}

type GoodsCateID struct {
	ID int `json:"id"`
}

func (con GoodsCateController) Delete(c *gin.Context) {
	//需要判断当前删除的商品分类是否是顶级商品分类(一级商品分类)
	//如果当前删除一级商品分类，下面如果有二级商品分类，就不能删除它，否则可以删除它
	goodsCateID := GoodsCateID{}
	if err := c.ShouldBind(&goodsCateID); err != nil {
		con.Error(c, "删除商品分类失败，请稍后重试", -1, nil)
		return
	}
	goodsCate := models.GoodsCate{ID: goodsCateID.ID}
	if err := dao.DB.Find(&goodsCate).Error; err != nil {
		con.Error(c, "删除商品分类失败，请稍后重试", -1, nil)
		return
	}

	if goodsCate.Pid == 0 { //是一级商品分类
		goodsCateList := []models.GoodsCate{}
		if err := dao.DB.Where("pid = ?", goodsCate.ID).Find(&goodsCateList).Error; err != nil {
			con.Error(c, "删除商品分类失败，请稍后重试", -1, nil)
			return
		}
		if len(goodsCateList) > 0 { //当前一级商品分类下还有二级商品分类
			con.Error(c, "删除失败，此权限下还有子商品分类", -1, nil)
			return
		}
		if err := dao.DB.Delete(&goodsCate).Error; err != nil {
			con.Error(c, "删除商品分类失败，请稍后重试", -1, nil)
			return
		}
	} else { //不是一级商品分类
		if err := dao.DB.Delete(&goodsCate).Error; err != nil {
			con.Error(c, "删除商品分类失败，请稍后重试", -1, nil)
			return
		}
	}
	con.Error(c, "删除商品分类成功", 0, nil)
}
