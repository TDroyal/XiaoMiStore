package models

import (
	"time"

	"gorm.io/gorm"
)

type Goods struct {
	ID            int     `gorm:"primaryKey"`
	Title         string  `gorm:"type:varchar(255)"`
	SubTitle      string  `gorm:"type:varchar(255)"`
	GoodsSn       string  `gorm:"type:varchar(255)"`
	CateID        int     //所属的分类ID  参考GoodsCate表
	ClickCount    int     //点击数量
	GoodsNumber   int     //库存
	Price         float64 `gorm:"type:decimal(10,2)"` //价格
	MarketPrice   float64 `gorm:"type:decimal(10,2)"` //市场价格（原价）
	RelationGoods string  `gorm:"type:varchar(255)"`  //关联商品 1,23,24  表示这些商品和当前商品关联   比如小米9-8G-256G  小米9-4G-128G
	GoodsAttr     string  `gorm:"type:varchar(1024)"` //商品属性   额外的属性  格式： 颜色：红色，白色 | 尺寸：41，42，43
	GoodsColor    string  `gorm:"type:varchar(255)"`  //商品颜色
	GoodsVersion  string  `gorm:"type:varchar(255)"`  //商品版本
	GoodsImg      string  `gorm:"type:varchar(255)"`  //商品图片
	GoodsGift     string  `gorm:"type:varchar(255)"`  //商品赠品
	GoodsFitting  string  `gorm:"type:varchar(255)"`  //商品配件
	GoodsKeywords string  `gorm:"type:varchar(255)"`  //商品关键词
	GoodsDesc     string  `gorm:"type:varchar(255)"`  //商品对应的描述
	GoodsContent  string  `gorm:"type:text"`          //商品详情
	// IsDelete      int     `gorm:"type:tinyint"`
	DeletedAt   gorm.DeletedAt `gorm:"index"`
	IsHot       int            `gorm:"type:tinyint"`
	IsBest      int            `gorm:"type:tinyint"`
	IsNew       int            `gorm:"type:tinyint"`
	GoodsTypeID int            //商品类型ID 对应GoodsType表
	Sort        int
	Status      int `gorm:"type:tinyint"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (Goods) TableName() string {
	return "goods"
}
