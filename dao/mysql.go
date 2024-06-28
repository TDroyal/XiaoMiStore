package dao

import (
	"XiaoMiStore/models"
	"log"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	DB *gorm.DB
)

// 连接数据库
func InitMySQL() error {
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
	if err := DB.AutoMigrate(&models.Admin{}, &models.Role{}, &models.Access{}, &models.RoleAccess{}, &models.Focus{}); err != nil {
		return err
	}
	return nil
}
