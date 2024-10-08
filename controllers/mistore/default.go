package mistore

import (
	"XiaoMiStore/dao"
	"XiaoMiStore/logic"
	"XiaoMiStore/models"
	"fmt"
	"time"

	"context"

	"github.com/gin-gonic/gin"

	pb "XiaoMiStore/proto/goods" // 把服务端的proto文件拿过来

	"go-micro.dev/v4/logger"
)

type DefaultController struct {
	BaseController
}

// 前台首页获取相应的数据（对相应的数据进行缓存）
// 1.顶部导航   2.轮播图数据   3.分类的数据   4.中间导航   5.底部导航

// 测试结果，目前数据库中就一条数据，如果数据量越大，那么速度越明显
// [0.000ms] [rows:1] SELECT * FROM `nav` WHERE position = 1 and status = 1
// 从mysql中读取首页顶部导航数据
// [GIN] 2024/07/24 - 21:04:11 | 200 |      1.2164ms |       127.0.0.1 | GET      "/"
// 从redis中读取首页顶部导航数据
// [GIN] 2024/07/24 - 21:04:15 | 200 |       392.7µs |       127.0.0.1 | GET      "/"

func (con DefaultController) Index(c *gin.Context) {
	// 1.获取顶部导航
	topNavList := []models.Nav{}
	if rErr := logic.RGet("topNavList", &topNavList); rErr != nil { //要么是redis中未缓存此数据，要么是从redis中读此数据出错
		//那么就从mysql数据库中读取数据
		if err := dao.DB.Where("position = ? and status = ?", 1, 1).Find(&topNavList).Error; err != nil {
			con.Error(c, "获取首页顶部导航数据失败", -1, nil)
			return
		}
		// 如果从mysql读取数据成功，那么进行数据的缓存
		if err := logic.RSet("topNavList", &topNavList, time.Hour*1); err != nil {
			fmt.Println("首页顶部导航数据缓存失败" + err.Error()) //打印到日志中
		}
		fmt.Println("从mysql中读取首页顶部导航数据")
	} else {
		fmt.Println("从redis中读取首页顶部导航数据")
	}
	// 按需要写即可 2.轮播图数据   3.分类的数据   4.中间导航   5.底部导航 （参考76）

	con.Success(c, "获取首页数据成功", 0, gin.H{
		"topNavList": topNavList,
	})
}

// 测试封装的cookie操作
func (con DefaultController) TestCookie(c *gin.Context) {
	type User struct {
		Username string
		Password string
		age      int
	}
	u := User{
		"royal_111",
		"123456",
		18,
	}
	err := logic.Cookie.Set(c, "user", &u)
	var getUser User
	ok := logic.Cookie.Get(c, "user", &getUser)
	con.Success(c, "获取cookie数据成功", 0, gin.H{
		"error": err,
		"ok":    ok,
		"user":  getUser,
	})
}

// 调用goods微服务
var (
	service = "goods"
	version = "latest"
)

func (con DefaultController) AddGoods(c *gin.Context) {
	// 如果未开启微服务，调用这个接口会让程序崩溃，怎么解决 用panic和recover
	defer func() {
		if r := recover(); r != nil {
			logger.Error("Recovered from panic:", r)
			con.Error(c, "An error occurred while processing the request", -1, nil) // 返回适当的错误响应
		}
	}()

	// Create client
	gclient := pb.NewGoodsService(service, logic.GoodsMicroClient)

	// Call service
	rsp, err := gclient.AddGoods(context.Background(), &pb.AddGoodsRequest{
		Title:   "小米手机",
		Price:   9999.9,
		Content: "便宜好用",
	})

	if err != nil {
		// logger.Fatal(err)
		panic(err) // 主动触发 panic
	}

	logger.Info(rsp)

	con.Success(c, rsp.GetMessage(), 0, gin.H{
		"status": rsp.GetSuccess(),
	})
}
