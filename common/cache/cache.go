package cache

import (
	"context"
	"github.com/redis/go-redis/v9"
	"sxp-server/common/logger"
)

var CacheClient *redis.Client

func IniCache() *redis.Client {
	l := logger.GetLogger()
	rdb := redis.NewClient(&redis.Options{
		Addr:     "192.168.111.143:6379",
		Password: "",
		DB:       0,
	})
	_, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		l.Panicf("redis启动失败: %s", err.Error())
	}
	return rdb
}
