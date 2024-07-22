package admin

import (
	"XiaoMiStore/dao"
	"XiaoMiStore/logic"
	"XiaoMiStore/models"
	"math"
	"net/http"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var wg sync.WaitGroup

type GoodsController struct {
	BaseController
}

// 获取商品列表信息  (分页)
// 前端   xxx/goodslist/?page=1 或者xxx/goodslist/  都是返回第一页的信息
func (con GoodsController) GetGoodsList(c *gin.Context) {

	page := logic.StringToInt(c.DefaultQuery("page", "1"))
	if page == 0 {
		page = 1
	}
	pageSize := 5

	// 如果搜索框中有内容，应该把搜索框中的内容带上
	keyword := c.Query("keyword")
	var condition string
	if len(keyword) > 0 {
		condition = "title like \"%" + keyword + "%\""
	}

	goodsList := []models.Goods{}
	if err := dao.DB.Where(condition).Offset((page - 1) * pageSize).Limit(pageSize).Find(&goodsList).Error; err != nil {
		con.Success(c, "获取商品列表信息失败", -1, nil)
		return
	}

	// 获取商品总数量
	var totalCount int64
	if err := dao.DB.Model(&models.Goods{}).Where(condition).Count(&totalCount).Error; err != nil {
		con.Success(c, "获取商品列表信息失败", -1, nil)
		return
	}

	con.Success(c, "获取商品列表信息成功", 0, gin.H{
		"goodsList":  goodsList,
		"totalCount": totalCount,
		"page":       page,
		"totalPage":  math.Ceil(float64(totalCount) / float64(pageSize)),
	})
}

func (con GoodsCateController) GetGoodsInfo(c *gin.Context) {
	//根据商品ID来获取商品的所有信息
	id := logic.StringToInt(c.Query("id")) //商品ID

	//1. 获取商品数据
	goods := models.Goods{ID: id}
	if err := dao.DB.Find(&goods).Error; err != nil {
		con.Error(c, "获取商品信息失败", -1, nil)
		return
	}

	//2. 获取商品分类
	goodsCateList := []models.GoodsCate{}
	if err := dao.DB.Where("pid = 0").Preload("GoodsCateItems").Find(&goodsCateList).Error; err != nil {
		con.Error(c, "获取商品信息失败", -1, nil)
		return
	}

	//3. 获取所有颜色 以及选中的颜色

	goodsColorSelected := strings.Split(goods.GoodsColor, ",") //"1,2,5"
	goodsColorMap := make(map[string]string)
	for _, v := range goodsColorSelected {
		goodsColorMap[v] = v
	}

	goodsColorList := []models.GoodsColor{}
	if err := dao.DB.Find(&goodsColorList).Error; err != nil {
		con.Error(c, "获取商品信息失败", -1, nil)
		return
	}

	colorLength := len(goodsColorList)
	for i := 0; i < colorLength; i++ {
		if _, ok := goodsColorMap[logic.IntToString(goodsColorList[i].ID)]; ok {
			goodsColorList[i].Checked = true //此颜色的checkbox框被选中
		}
	}

	//4. 商品的图库信息
	goodsImageList := []models.GoodsImage{}
	if err := dao.DB.Where("goods_id = ?", id).Find(&goodsImageList).Error; err != nil {
		con.Error(c, "获取商品信息失败", -1, nil)
		return
	}

	//5. 获取商品类型
	goodsTypeList := []models.GoodsType{}
	if err := dao.DB.Find(&goodsTypeList).Error; err != nil {
		con.Error(c, "获取商品信息失败", -1, nil)
		return
	}

	//6. 获取规格信息
	goodsAttr := []models.GoodsAttr{}
	if err := dao.DB.Where("goods_id = ?", id).Find(&goodsAttr).Error; err != nil {
		con.Error(c, "获取商品信息失败", -1, nil)
		return
	}

	con.Success(c, "获取商品信息成功", 0, gin.H{
		"goods":          goods,
		"goodsCateList":  goodsCateList,
		"goodsColorList": goodsColorList,
		"goodsTypeList":  goodsTypeList,
		"goodsAttr":      goodsAttr,
		"goodsImageList": goodsImageList,
	})
}

// 配合https://froala.com/wysiwyg-editor/docs/options/这个富文本编辑器完成图片上传的
func (con GoodsController) ImageUpload(c *gin.Context) {
	//上传图片
	imgDir, err := logic.UploadImageFile(c, "file") //注意可以在网络里面看到传递的参数为file
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"link": "",
		})
		return
	}

	if logic.GetOssStatus() != 1 { //将图片存到本地，就生成缩略图
		wg.Add(1)
		go func() {
			defer wg.Done()
			logic.ResizeGoodsImage(imgDir)
		}()
		c.JSON(http.StatusOK, gin.H{
			"link": imgDir, //{link: 'path/to/image.jpg'}  需要把后端ip拼接在最前面
		})
	} else { //开启oss，拼接的地址应该加上oss的域名
		domain, _ := logic.GetSettingFromColumn("OssDomain")
		c.JSON(http.StatusOK, gin.H{
			"link": domain + imgDir, //{link: 'path/to/image.jpg'}  需要把oss 域名拼接在最前面
		})
	}

}

// *********************这种方法要弄明白*********

// 1. 通用信息（例如商品标题、所属分类、商品价格等存在goods表里面）

// 2. 规格包装（商品类型属性表GoodsTypeAttribute生成的表单信息）（例如商品类型、基本信息、性能等放在goods_attr表里面）

// 3. 商品相册（放在goods_image表里面）

func (con GoodsController) Add(c *gin.Context) {
	//最原始的点击按钮（type="submit"）提交form表单，如果前端一个form表单中有多个相同name="attr_id_list"的input框
	// attrIdList := c.PostFormArray("attr_id_list")

	// // 上传的多个商品属性ID以及属性ID对应的值
	// attrIdList := c.PostFormArray("attr_id_list")
	// attrValueList := c.PostFormArray("attr_value_list")
	// // 上传的多个商品图片（"多个图片上传成功后，前端得到了这些图片的url绑定在了隐藏的input输入框中，再传到后端"）
	// goodsImageList := c.PostFormArray("goods_image_list")

	// 1.获取表单提交过来的数据  可以改成shouldbind
	title := c.PostForm("title")
	subTitle := c.PostForm("sub_title")
	goodsSn := c.PostForm("goods_sn")
	cateID := logic.StringToInt(c.PostForm("cate_id"))
	goodsNumber := logic.StringToInt(c.PostForm("goods_number"))
	//注意小数点
	marketPrice := logic.StringToFloat(c.PostForm("market_price"))
	price := logic.StringToFloat(c.PostForm("price"))

	relationGoods := c.PostForm("relation_goods")
	goodsAttr := c.PostForm("goods_attr") //更多属性，额外的属性
	goodsVersion := c.PostForm("goods_version")
	goodsGift := c.PostForm("goods_gift")
	goodsFitting := c.PostForm("goods_fitting")
	//获取的是切片
	goodsColorArr := c.PostFormArray("goods_color") //checkbox的名称都是goods_color

	goodsKeywords := c.PostForm("goods_keywords")
	goodsDesc := c.PostForm("goods_desc")
	goodsContent := c.PostForm("goods_content")
	isHot := logic.StringToInt(c.PostForm("is_hot"))
	isBest := logic.StringToInt(c.PostForm("is_best"))
	isNew := logic.StringToInt(c.PostForm("is_new"))
	goodsTypeID := logic.StringToInt(c.PostForm("goods_type_id"))
	sort := logic.StringToInt(c.PostForm("sort"))
	status := logic.StringToInt(c.PostForm("status"))

	// 2.获取颜色信息 把颜色转化为字符串   //存到goods表的goods_color string中
	goodsColorStr := strings.Join(goodsColorArr, ",")

	// 3.上传图片 生成缩略图   商品的主图片 goods表中的GoodsImg
	goodsImg, upload_err := logic.UploadImageFile(c, "goods_img")
	if upload_err != nil {
		con.Error(c, "添加商品失败", -1, nil)
		return
	}

	// 生成缩略图
	if logic.GetOssStatus() != 1 { //将图片存到本地，就生成缩略图
		wg.Add(1)
		go func() {
			defer wg.Done()
			logic.ResizeGoodsImage(goodsImg)
		}()
	}

	// 下面3步要么都成功，要么都不成功，事务维护数据库的一致性

	// 4.增加商品数据
	goods := models.Goods{
		Title:         title,
		SubTitle:      subTitle,
		GoodsSn:       goodsSn,
		CateID:        cateID,
		ClickCount:    100,
		GoodsNumber:   goodsNumber,
		MarketPrice:   marketPrice,
		Price:         price,
		RelationGoods: relationGoods,
		GoodsAttr:     goodsAttr,
		GoodsVersion:  goodsVersion,
		GoodsGift:     goodsGift,
		GoodsFitting:  goodsFitting,
		GoodsKeywords: goodsKeywords,
		GoodsDesc:     goodsDesc,
		GoodsContent:  goodsContent,
		IsHot:         isHot,
		IsBest:        isBest,
		IsNew:         isNew,
		GoodsTypeID:   goodsTypeID,
		Sort:          sort,
		Status:        status,
		GoodsColor:    goodsColorStr,
		GoodsImg:      goodsImg,
	}

	// 内部的报错会传给Transaction的返回值吗？不行应该怎么解决
	if err := dao.DB.Transaction(func(tx *gorm.DB) error {
		// 返回任何错误都会回滚事务

		// 4.增加商品数据
		if err := tx.Create(&goods).Error; err != nil {
			con.Error(c, "添加商品失败", -1, nil)
			return err
		}

		// 5.增加图库信息
		// 上传的多个商品图片（"多个图片上传成功后，前端得到了这些图片的url绑定在了隐藏的input输入框中，再传到后端"）
		errChan := make(chan error, 1)
		wg.Add(1)
		go func() {
			defer wg.Done() // 在函数结束时调用 Done() 方法
			goodsImageList := c.PostFormArray("goods_image_list")
			for _, v := range goodsImageList {
				goodsImgObj := models.GoodsImage{
					GoodsID: goods.ID,
					ImgUrl:  v,
					Sort:    10,
					Status:  1,
				}
				if err := tx.Create(&goodsImgObj).Error; err != nil {
					con.Error(c, "添加商品失败", -1, nil)
					errChan <- err
					return
				}
			}
		}()

		// 6.增加规格包装
		// 上传的多个商品属性ID以及属性ID对应的值  goods_attr添加数据
		wg.Add(1)
		go func() {
			attrIdList := c.PostFormArray("attr_id_list")
			attrValueList := c.PostFormArray("attr_value_list")
			length := len(attrIdList)
			for i := 0; i < length; i++ {
				defer wg.Done() // 在函数结束时调用 Done() 方法
				// 根据attrIdList获取商品类型属性表(goodsTypeAttribute)的数据
				GoodsTypeAttributeID := logic.StringToInt(attrIdList[i])
				GoodsTypeAttributeObj := models.GoodsTypeAttribute{ID: GoodsTypeAttributeID}
				if err := tx.Find(&GoodsTypeAttributeObj).Error; err != nil {
					con.Error(c, "添加商品失败", -1, nil)
					errChan <- err
					return
				}

				//给商品属性里面增加数据
				goodsAttrObj := models.GoodsAttr{
					GoodsID:         goods.ID,
					AttributeTitle:  GoodsTypeAttributeObj.Title,
					AttributeType:   GoodsTypeAttributeObj.AttrType,
					AttributeID:     GoodsTypeAttributeObj.ID,
					AttributeCateID: GoodsTypeAttributeObj.CateID,
					AttributeValue:  attrValueList[i],
					Status:          1,
					Sort:            10,
				}
				if err := tx.Create(&goodsAttrObj).Error; err != nil {
					con.Error(c, "添加商品失败", -1, nil)
					errChan <- err
					return
				}
			}
		}()

		// 等待所有协程完成
		go func() {
			wg.Wait()
			close(errChan) // 关闭通道
		}()

		// 从通道中接收错误信息
		for err := range errChan { // 通道将保持打开状态，range 循环将一直等待，直到通道关闭或有新的值被发送到通道中。
			if err != nil {
				con.Error(c, "添加商品失败", -1, nil)
				return err
			}
		}

		// 返回 nil 提交事务
		return nil
	}); err != nil {
		con.Error(c, "添加商品失败", -1, nil)
		return
	}

	con.Success(c, "添加商品成功", 0, nil)
}

func (con GoodsController) Edit(c *gin.Context) {

	// 1.获取表单提交过来的数据  可以改成shouldbind
	id := logic.StringToInt(c.PostForm("id"))
	title := c.PostForm("title")
	subTitle := c.PostForm("sub_title")
	goodsSn := c.PostForm("goods_sn")
	cateID := logic.StringToInt(c.PostForm("cate_id"))
	goodsNumber := logic.StringToInt(c.PostForm("goods_number"))
	//注意小数点
	marketPrice := logic.StringToFloat(c.PostForm("market_price"))
	price := logic.StringToFloat(c.PostForm("price"))

	relationGoods := c.PostForm("relation_goods")
	goodsAttr := c.PostForm("goods_attr") //更多属性，额外的属性
	goodsVersion := c.PostForm("goods_version")
	goodsGift := c.PostForm("goods_gift")
	goodsFitting := c.PostForm("goods_fitting")
	//获取的是切片
	goodsColorArr := c.PostFormArray("goods_color") //checkbox的名称都是goods_color

	goodsKeywords := c.PostForm("goods_keywords")
	goodsDesc := c.PostForm("goods_desc")
	goodsContent := c.PostForm("goods_content")
	isHot := logic.StringToInt(c.PostForm("is_hot"))
	isBest := logic.StringToInt(c.PostForm("is_best"))
	isNew := logic.StringToInt(c.PostForm("is_new"))
	goodsTypeID := logic.StringToInt(c.PostForm("goods_type_id"))
	sort := logic.StringToInt(c.PostForm("sort"))
	status := logic.StringToInt(c.PostForm("status"))

	// 2.获取颜色信息 把颜色转化为字符串   //存到goods表的goods_color string中
	goodsColorStr := strings.Join(goodsColorArr, ",")

	// 3.上传图片 生成缩略图   商品的主图片 goods表中的GoodsImg
	goodsImg, upload_err := logic.UploadImageFile(c, "goods_img")
	if upload_err != nil {
		con.Error(c, "添加商品失败", -1, nil)
		return
	}

	// 生成缩略图
	if logic.GetOssStatus() != 1 { //将图片存到本地，就生成缩略图
		wg.Add(1)
		go func() {
			defer wg.Done()
			logic.ResizeGoodsImage(goodsImg)
		}()
	}

	// 下面3步要么都成功，要么都不成功，事务维护数据库的一致性

	// 4.修改商品数据
	goods := models.Goods{
		ID:            id,
		Title:         title,
		SubTitle:      subTitle,
		GoodsSn:       goodsSn,
		CateID:        cateID,
		GoodsNumber:   goodsNumber,
		MarketPrice:   marketPrice,
		Price:         price,
		RelationGoods: relationGoods,
		GoodsAttr:     goodsAttr,
		GoodsVersion:  goodsVersion,
		GoodsGift:     goodsGift,
		GoodsFitting:  goodsFitting,
		GoodsKeywords: goodsKeywords,
		GoodsDesc:     goodsDesc,
		GoodsContent:  goodsContent,
		IsHot:         isHot,
		IsBest:        isBest,
		IsNew:         isNew,
		GoodsTypeID:   goodsTypeID,
		Sort:          sort,
		Status:        status,
		GoodsColor:    goodsColorStr,
		GoodsImg:      goodsImg,
	}

	// 内部的报错会传给Transaction的返回值吗？不行应该怎么解决
	if err := dao.DB.Transaction(func(tx *gorm.DB) error {
		// 返回任何错误都会回滚事务

		// 4.增加商品数据
		if err := tx.Save(&goods).Error; err != nil {
			con.Error(c, "修改商品失败", -1, nil)
			return err
		}

		// 5.修改图库 增加图库信息
		// 上传的多个商品图片（"多个图片上传成功后，前端得到了这些图片的url绑定在了隐藏的input输入框中，再传到后端"）
		errChan := make(chan error, 1)
		wg.Add(1)
		go func() {
			defer wg.Done() // 在函数结束时调用 Done() 方法
			goodsImageList := c.PostFormArray("goods_image_list")
			for _, v := range goodsImageList {
				goodsImgObj := models.GoodsImage{
					GoodsID: goods.ID,
					ImgUrl:  v,
					Sort:    10,
					Status:  1,
				}
				if err := tx.Create(&goodsImgObj).Error; err != nil {
					con.Error(c, "修改商品失败", -1, nil)
					errChan <- err
					return
				}
			}
		}()

		// 6.修改规格包装 (要先删除原来的，再添加新的)
		// 上传的多个商品属性ID以及属性ID对应的值  goods_attr添加数据

		wg.Add(1)
		go func() {
			goodsAttrObj := models.GoodsAttr{}
			dao.DB.Where("goods_id = ?", goods.ID).Delete(goodsAttrObj)

			attrIdList := c.PostFormArray("attr_id_list")
			attrValueList := c.PostFormArray("attr_value_list")
			length := len(attrIdList)
			for i := 0; i < length; i++ {
				defer wg.Done() // 在函数结束时调用 Done() 方法
				// 根据attrIdList获取商品类型属性表(goodsTypeAttribute)的数据
				GoodsTypeAttributeID := logic.StringToInt(attrIdList[i])
				GoodsTypeAttributeObj := models.GoodsTypeAttribute{ID: GoodsTypeAttributeID}
				if err := tx.Find(&GoodsTypeAttributeObj).Error; err != nil {
					con.Error(c, "修改商品失败", -1, nil)
					errChan <- err
					return
				}

				//给商品属性里面增加数据
				goodsAttrObj := models.GoodsAttr{
					GoodsID:         goods.ID,
					AttributeTitle:  GoodsTypeAttributeObj.Title,
					AttributeType:   GoodsTypeAttributeObj.AttrType,
					AttributeID:     GoodsTypeAttributeObj.ID,
					AttributeCateID: GoodsTypeAttributeObj.CateID,
					AttributeValue:  attrValueList[i],
					Status:          1,
					Sort:            10,
				}
				if err := tx.Create(&goodsAttrObj).Error; err != nil {
					con.Error(c, "修改商品失败", -1, nil)
					errChan <- err
					return
				}
			}
		}()

		// 等待所有协程完成
		go func() {
			wg.Wait()
			close(errChan) // 关闭通道
		}()

		// 从通道中接收错误信息
		for err := range errChan { // 通道将保持打开状态，range 循环将一直等待，直到通道关闭或有新的值被发送到通道中。
			if err != nil {
				con.Error(c, "修改商品失败", -1, nil)
				return err
			}
		}

		// 返回 nil 提交事务
		return nil
	}); err != nil {
		con.Error(c, "修改商品失败", -1, nil)
		return
	}

	con.Success(c, "修改商品成功", 0, nil)
}

// 删除商品
func (con GoodsController) Delete(c *gin.Context) {
	id := logic.StringToInt(c.PostForm("id"))
	// 删除goods表，goods_attr表，goods_image表  事务

	// 事务保证3个同时删除了
	dao.DB.Transaction(func(tx *gorm.DB) error {
		// 返回任何错误都会回滚事务
		goods := models.Goods{ID: id}
		if err := tx.Delete(&goods).Error; err != nil {
			return err
		}

		if err := tx.Where("goods_id = ?", id).Delete(&models.GoodsAttr{}).Error; err != nil {
			return err
		}

		if err := tx.Where("goods_id = ?", id).Delete(&models.GoodsImage{}).Error; err != nil {
			return err
		}

		// 返回 nil 提交事务
		return nil
	})
}

// 异步修改商品相册图片和颜色绑定信息
func (con GoodsController) ChangeGoodsImageColor(c *gin.Context) {
	goods_image_id := logic.StringToInt(c.PostForm("goods_image_id"))
	color_id := logic.StringToInt(c.PostForm("color_id"))
	if err := dao.DB.Model(&models.GoodsImage{}).Where("id = ?", goods_image_id).Update("color_id", color_id).Error; err != nil {
		con.Error(c, "商品图库颜色绑定修改失败", -1, nil)
		return
	}
	con.Success(c, "商品图库颜色绑定修改成功", 0, nil)
}

// 异步删除商品相册信息
func (con GoodsController) RemoveGoodsImage(c *gin.Context) {
	goods_image_id := logic.StringToInt(c.PostForm("goods_image_id"))

	if err := dao.DB.Delete(&models.GoodsImage{}, goods_image_id).Error; err != nil {
		con.Error(c, "删除商品相册信息失败", -1, nil)
		return
	}
	con.Success(c, "删除商品相册信息成功", 0, nil)
}
