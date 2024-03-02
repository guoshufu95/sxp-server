package service

import (
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	ini "sxp-server/common/initial"
	"sxp-server/common/logger"
)

type Service struct {
	Db     *gorm.DB
	Cache  *redis.Client
	Logger *logger.ZapLog
}

func MakeService(s *Service) {
	app := ini.MakeApp()
	s.Db = app.Db
	s.Cache = app.Cache
	s.Logger = app.Logger
}
