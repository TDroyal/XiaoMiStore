package dao

import (
	"XiaoMiStore/models"
	"fmt"
	"log"
	"os"

	"gopkg.in/ini.v1"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	DB *gorm.DB
)

// 原始连接数据库
func OldInitMySQL() error {
	var err error
	dsn := "root:123456@tcp(127.0.0.1:3306)/xiaomistore?charset=utf8mb4&parseTime=True&loc=Local"

	// 配置日志记录器，设置为log.New(os.Stdout, "\r\n", log.LstdFlags)以将日志输出到控制台
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			LogLevel: logger.Info, // 设置日志级别为Info，将打印SQL语句
		},
	)

	if DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: newLogger,
	}); err != nil {
		return err
	}
	DB.Debug() //打印sql语句
	return nil
}

// 连接数据库，从ini文件中读出mysql的相关配置
func InitMySQL() error {

	config, iniErr := ini.Load("conf/app.ini")
	if iniErr != nil {
		fmt.Printf("Fail to read file: %v", iniErr)
		return iniErr
	}

	// 从ini文件中读出配置
	ip := config.Section("mysql").Key("ip").String()
	port := config.Section("mysql").Key("port").String()
	user := config.Section("mysql").Key("user").String()
	password := config.Section("mysql").Key("password").String()
	database := config.Section("mysql").Key("database").String()

	var err error
	// dsn := "root:123456@tcp(127.0.0.1:3306)/xiaomistore?charset=utf8mb4&parseTime=True&loc=Local"
	dsn := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?charset=utf8mb4&parseTime=True&loc=Local", user, password, ip, port, database)
	// 配置日志记录器，设置为log.New(os.Stdout, "\r\n", log.LstdFlags)以将日志输出到控制台
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			LogLevel: logger.Info, // 设置日志级别为Info，将打印SQL语句
		},
	)

	if DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: newLogger,
	}); err != nil {
		return err
	}
	DB.Debug() //打印sql语句
	return nil
}

// 关闭数据库
func CloseMySQL() {
	sqlDB, err := DB.DB()
	if err != nil {
		panic(err)
	}
	sqlDB.Close()
}

// 迁移模型  AutoMigrate 用于自动迁移您的 schema，保持您的 schema 是最新的。
func InitModels() error {
	if err := DB.AutoMigrate(&models.Admin{}, &models.Role{}, &models.Access{}, &models.RoleAccess{}, &models.Focus{}, &models.GoodsCate{}, &models.GoodsType{}, &models.GoodsTypeAttribute{}, &models.GoodsColor{}, &models.Goods{}, &models.GoodsAttr{}, &models.GoodsImage{}, &models.Nav{}, &models.Setting{}, &models.UserTemp{}, &models.User{}); err != nil {
		return err
	}
	return nil
}
