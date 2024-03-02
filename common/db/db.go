package db

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	_ "gorm.io/gorm" // gorm
	"gorm.io/gorm/logger"
	zaplog "sxp-server/common/logger"
	"time"
)

func IniDb() *gorm.DB {
	l := zaplog.GetLogger()
	dsn := "root:123456@tcp(192.168.111.143:3306)/sxp-server?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: zaplog.New(
			logger.Config{
				SlowThreshold: time.Second,
				Colorful:      true,
				LogLevel: logger.LogLevel(
					4),
			},
		),
	})
	if err != nil {
		l.Panicf("连接mysql数据库失败:%s", err.Error())
	}
	return db
}
