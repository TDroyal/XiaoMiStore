package logic

// 操作Redis

import (
	"XiaoMiStore/dao"
	"context"
	"time"
)

var ctx = context.Background()

func RedisSet(key, val string, expiration time.Duration) error {
	if err := dao.RDB.Set(ctx, key, val, expiration).Err(); err != nil { //0表示在redis中存储，并且没有过期时间
		return err
	}
	return nil
}

func RedisGet(key string) (string, error) {
	val, err := dao.RDB.Get(ctx, key).Result()
	if err != nil {
		return "", err
	}
	return val, nil
}

func RedisDel(key string) error {
	err := dao.RDB.Del(ctx, key).Err()
	if err != nil {
		return err
	}
	return nil
}
