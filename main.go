package main

import (
	"XiaoMiStore/dao"
	"XiaoMiStore/routers"
)

func main() {
	// 连接MySQL数据库
	if err := dao.InitMySQL(); err != nil {
		panic(err)
	}

	//main结束时关闭连接数据库
	defer dao.CloseMySQL()

	// 模型绑定
	if err := dao.InitModels(); err != nil {
		panic(err)
	}

	// 连接Redis数据库
	if err := dao.InitRedis(); err != nil {
		panic(err)
	}

	//注册路由
	r := routers.SetupRouter()

	r.Run(":80")
}
