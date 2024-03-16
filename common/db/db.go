package db

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	_ "gorm.io/gorm" // gorm
	"gorm.io/gorm/logger"
	zaplog "sxp-server/common/logger"
	"sxp-server/config"
	"time"
)

func IniDb() *gorm.DB {
	l := zaplog.GetLogger()
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.Conf.Mysql.UserName,
		config.Conf.Mysql.Password,
		config.Conf.Mysql.Host,
		config.Conf.Mysql.Port,
		config.Conf.Mysql.Db)

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
	fmt.Println("mysql连接成功")
	// 表迁移，不用可注释
	//err = db.AutoMigrate(&model.User{}, &model.Role{}, &model.Menu{})
	//if err != nil {
	//	fmt.Println(err)
	//	return nil
	//}
	return db
}
