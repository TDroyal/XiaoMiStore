package logic

import (
	"fmt"

	"github.com/go-micro/plugins/v4/registry/consul"
	"go-micro.dev/v4"
	"go-micro.dev/v4/client"
	"go-micro.dev/v4/registry"
	"gopkg.in/ini.v1"
)

var (
	GoodsMicroClient   client.Client
	CaptchaMicroClient client.Client
	RoleMicroClient    client.Client
)

// 其它包引用了logic这个包，就会自动执行此init函数

// 下面部分代码是go-micro框架中的脚手架命令生成的
/*
go-micro: https://github.com/micro/go-micro
cli: https://github.com/go-micro/cli
*/

func init() {
	InitGoodsMicro()   // 初始化goods微服务
	InitCaptchaMicro() //初始化captcha微服务
	InitRoleMicro()    // 初始化role微服务
}

// 初始化goods微服务
func InitGoodsMicro() {
	config, iniErr := ini.Load("conf/app.ini")
	if iniErr != nil {
		fmt.Printf("Fail to read file: %v", iniErr)
		return
	}

	// 从ini文件中读出配置
	ip := config.Section("consul1").Key("ip").String()
	port := config.Section("consul1").Key("port").String()

	// 配置consul服务发现   自动实现负载均衡
	consulReg := consul.NewRegistry(
		registry.Addrs(fmt.Sprintf("%s:%s", ip, port)), //建议一个微服务连接一个consul client
	)

	// Create service

	srv := micro.NewService()

	srv.Init(
		micro.Registry(consulReg),
	)

	GoodsMicroClient = srv.Client()
}

// 初始化captcha验证码微服务
func InitCaptchaMicro() {
	config, iniErr := ini.Load("conf/app.ini")
	if iniErr != nil {
		fmt.Printf("Fail to read file: %v", iniErr)
		return
	}

	// 从ini文件中读出配置
	ip := config.Section("consul1").Key("ip").String()
	port := config.Section("consul1").Key("port").String()

	// 配置consul服务发现   自动实现负载均衡
	consulReg := consul.NewRegistry(
		registry.Addrs(fmt.Sprintf("%s:%s", ip, port)), //建议一个微服务连接一个consul client
	)
	// Create service

	srv := micro.NewService()
	srv.Init(
		micro.Registry(consulReg),
	)

	CaptchaMicroClient = srv.Client()
}

// 初始化role管理员角色管理微服务
func InitRoleMicro() {
	config, iniErr := ini.Load("conf/app.ini")
	if iniErr != nil {
		fmt.Printf("Fail to read file: %v", iniErr)
		return
	}

	// 从ini文件中读出配置
	ip := config.Section("consul1").Key("ip").String()
	port := config.Section("consul1").Key("port").String()

	// 配置consul服务发现   自动实现负载均衡
	consulReg := consul.NewRegistry(
		registry.Addrs(fmt.Sprintf("%s:%s", ip, port)), //建议一个微服务连接一个consul client
	)
	// Create service

	srv := micro.NewService()
	srv.Init(
		micro.Registry(consulReg),
	)

	RoleMicroClient = srv.Client()
}
