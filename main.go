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

	defer dao.CloseMySQL()

	// 连接Redis数据库
	if err := dao.InitRedis(); err != nil {
		panic(err)
	}

	//注册路由
	r := routers.SetupRouter()

	r.Run(":80")
}
