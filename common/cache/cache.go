package cache

import (
	"context"
	"github.com/redis/go-redis/v9"
	"sxp-server/common/logger"
	"sxp-server/config"
)

var CClient *redis.Client

func IniCache() *redis.Client {
	l := logger.GetLogger()
	rdb := redis.NewClient(&redis.Options{
		Addr:     config.Conf.Redis.Addr,
		Password: config.Conf.Redis.Password,
		DB:       0,
	})
	_, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		l.Panicf("redis启动失败: %s", err.Error())
	}
	CClient = rdb
	return rdb
}
