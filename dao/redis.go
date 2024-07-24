package dao

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
	"gopkg.in/ini.v1"
)

var (
	RDB         *redis.Client
	ctx         = context.Background()
	RedisEnable bool
)

// 原始连接redis
func OldInitRedis() error {
	RDB = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "123456", // no password set
		DB:       0,        // use default DB
	})

	_, err := RDB.Ping(ctx).Result()
	if err != nil {
		return err
	}
	return nil
}

func InitRedis() error {
	config, iniErr := ini.Load("conf/app.ini")
	if iniErr != nil {
		fmt.Printf("Fail to read file: %v", iniErr)
		return iniErr
	}

	// 从ini文件中读出配置
	ip := config.Section("redis").Key("ip").String()
	port := config.Section("redis").Key("port").String()
	password := config.Section("redis").Key("password").String()
	RedisEnable, _ = config.Section("redis").Key("redisEnable").Bool()

	if RedisEnable { //启动redis服务
		RDB = redis.NewClient(&redis.Options{
			Addr:     ip + ":" + port,
			Password: password, // no password set
			DB:       0,        // use default DB
		})

		_, err := RDB.Ping(ctx).Result()
		if err != nil {
			return err
		}
		return nil
	}
	return nil
}
