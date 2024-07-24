package logic

// 操作Redis

import (
	"XiaoMiStore/dao"
	"context"
	"encoding/json"
	"errors"
	"time"
)

var ctx = context.Background()

// redis字符串类型的设置和取值
func RedisSet(key string, val interface{}, expiration time.Duration) error {
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

// 可以在redis中存放 map, 结构体，string等类型的数据

func RSet(key string, val interface{}, expiration time.Duration) error {
	// val 可能是结构体，map，string等类型的数据，需要先把它转为json字符串
	if dao.RedisEnable {
		v, parseErr := json.Marshal(val)
		if parseErr != nil {
			return parseErr
		}

		if err := dao.RDB.Set(ctx, key, string(v), expiration).Err(); err != nil { //0表示在redis中存储，并且没有过期时间
			return err
		}
		return nil
	}
	// redis未开启 redis is closed
	return errors.New("redis is closed")
}

func RGet(key string, obj interface{}) error {
	if dao.RedisEnable {
		valStr, parseErr := dao.RDB.Get(ctx, key).Result()
		if parseErr != nil {
			return parseErr
		}
		// json.Unmarshal将[]byte类型的数据转为Json格式绑定到obj上
		if err := json.Unmarshal([]byte(valStr), &obj); err != nil {
			return err
		}
		return nil
	}
	// redis未开启 redis is closed
	return errors.New("redis is closed")
}
