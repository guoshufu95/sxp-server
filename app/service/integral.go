package service

import (
	"math/rand"
	"sxp-server/app/dao"
	"sxp-server/app/service/dto"
	"sxp-server/common/cal"
	"time"
)

type IntegralService struct {
	Service
}

func (s *IntegralService) InitIntegral(integral dto.IntegralReq) (err error) {
	for i := 0; integral.RemainCount > 0; i++ {
		v := cal.GrabIntegral(&integral)
		rand.Seed(time.Now().UnixNano())
		randomId := rand.Intn(100) // 生成0~9的随机数
		id := randomId + v
		//把积分初始化信息放入mysql
		err = dao.InertIntegralInfo(s.Db, id, v)
		if err != nil {
			s.Logger.Error(err.Error())
			return
		}
	}
	return
}
