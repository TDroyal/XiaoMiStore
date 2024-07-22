package models

// 网站的系统设置，基本信息也存储起来（数据库中只有一条数据）
type Setting struct {
	ID              int    `form:"id"`
	SiteTitle       string `form:"site_title"` //网站名称 有助于seo优化
	SiteLogo        string //网站logo
	SiteKeywords    string `form:"site_keywords"`    // 有助于seo优化
	SiteDescription string `form:"site_description"` // 有助于seo优化
	NoPicture       string //有些图片不显示时，显示的默认图片
	SiteIcp         string `form:"site_icp"`        //备案号
	SiteTel         string `form:"site_tel"`        //联系方式
	SearchKeywords  string `form:"search_keywords"` //搜索关键词
	TongjiCode      string `form:"tongji_code"`     //统计代码

	//对象存储相关的 (在阿里云的oss服务里面都能查找到)
	Appid      string `form:"app_id"`      // oss的AccessKey ID   OSS_ACCESS_KEY_ID
	AppSecret  string `form:"app_secret"`  // oss的AccessKey Secret  OSS_ACCESS_KEY_SECRET
	EndPoint   string `form:"end_point"`   // Bucket对应的Endpoint  EndPoint
	BucketName string `form:"bucket_name"` // Bucket的name  BucketName
	OssStatus  int    `form:"oss_status"`  // 1 表示开启   0表示不开启  弃用这个，直接设置在后端配置文件里面即可
	OssDomain  string `form:"oss_domain"`  // 存储与oss域名进行绑定的域名，用于拼接图片地址用的
	// 如果开启了oss对象存储服务，那么把图片上传至对象存储云服务器，否则还是存到本地

	// 系统生成哪些缩略图尺寸 以逗号隔开
	ThumbnailSize string `form:"thumbnail_size"` //200,300,500生成200*200，300*300，500*500的缩略图，用go image实现
}
