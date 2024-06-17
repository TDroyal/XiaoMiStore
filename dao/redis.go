package dao

import (
	"context"

	"github.com/redis/go-redis/v9"
)

var (
	RDB *redis.Client
	ctx = context.Background()
)

func InitRedis() error {
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
