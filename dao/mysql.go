package dao

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	DB *gorm.DB
)

// 连接数据库
func InitMySQL() error {
	var err error
	dsn := "root:123456@tcp(127.0.0.1:3306)/xiaomistore?charset=utf8mb4&parseTime=True&loc=Local"
	if DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{}); err != nil {
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
