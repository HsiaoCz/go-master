package db

import (
	"context"

	"github.com/go-redis/redis/v8"
)

var RDB *redis.Client

func InitRedis() error {
	//  Initialize Redis
	RDB = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	_, err := RDB.Ping(context.Background()).Result()
	return err
}
