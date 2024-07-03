package models

import "time"

//商品分类 (小米商城首页左侧的分类：例如手机、电视、家电等等)
type GoodsCate struct {
	ID          int    `gorm:"primaryKey"`
	Title       string `gorm:"type:varchar(255);"` //标题 例如手机、电视、家电等等
	CateImg     string `gorm:"type:varchar(255);"` //分类的图片
	Link        string `gorm:"type:varchar(255)"`  //分类对应的跳转地址
	Template    string `gorm:"type:varchar(255)"`  //自定义分类对应的模板  为空加载默认模板，不为空加载自定义模板
	Pid         int    //0表示是一个顶级分类  pid如果和此表中的ID相等，就表示这个ID下对应二级分类
	SubTitle    string `gorm:"type:varchar(255)"`  //seo的标题
	KeyWords    string `gorm:"type:varchar(255)"`  //seo关键词
	Description string `gorm:"type:varchar(1024)"` //seo描述
	Status      int    `gorm:"type:tinyint(1)"`    //状态
	Sort        int    //排序
	CreatedAt   time.Time
	UpdatedAt   time.Time
	//也权限表Access一样，也是一个自关联的表
	GoodsCateItems []GoodsCate `gorm:"foreignKey:Pid; references:ID"`
}

func (GoodsCate) TableName() string {
	return "goods_cate"
}
