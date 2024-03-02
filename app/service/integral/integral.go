package integral

import (
	"math/rand"
	"sxp-server/app/dao/integralDao"
	"sxp-server/app/model"
	"sxp-server/app/service"
	"sxp-server/common/cal"
	"time"
)

type IntegralService struct {
	service.Service
}

func (s *IntegralService) InitIntegral(integral model.IntegralReq) (err error) {
	for i := 0; integral.RemainCount > 0; i++ {
		v := cal.GrabIntegral(&integral)
		rand.Seed(time.Now().UnixNano())
		randomId := rand.Intn(100) // 生成0~9的随机数
		id := randomId + v
		//把积分初始化信息放入mysql
		err = integralDao.InertIntegralInfo(s.Db, id, v)
		if err != nil {
			s.Logger.Error(err.Error())
			return
		}
	}
	return
}
