package service

import (
	"context"
	"fmt"
	"math/rand"
	"sxp-server/app/dao"
	"sxp-server/app/service/dto"
	"sxp-server/common/cache"
	"sxp-server/common/cal"
	_const "sxp-server/common/const"
	"time"
)

type IntegralService struct {
	Service
}

// InitIntegral
//
//	@Description: 初始化积分相关信息到mysql和redis
//	@receiver s
//	@param integral
//	@return err
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
		//放入未消费的积分队列(list),存放的是id
		_, err = s.Cache.LPush(context.Background(), _const.NIntegralList, id).Result()
		if err != nil {
			s.Logger.Error("存放未消费的积分到redis失败")
			return
		}
		//把积分具体信息放入redis，用hash表储存(id+money)
		_, err = s.Cache.HSet(context.Background(), _const.IntegralInfo, id, v).Result()
		if err != nil {
			s.Logger.Error("存放未消费的积分到redis失败")
			return
		}
	}
	return
}

func (s *IntegralService) Do(name string) (err error) {
	var (
		i, res interface{}
		flag   bool
	)
	//限流处理
	i, err = s.Cache.Eval(context.Background(), cache.FilterScript, []string{_const.Key1}).Result()
	if err != nil {
		s.Logger.Error("限流功能错误！")
		return
	}
	if i.(int64) == 0 {
		fmt.Println("令牌已用完")
		return
	}
	if i.(int64) >= 1 {
		flag = true
		return
	}
	if flag {
		// 抢积分的业务逻辑
		res, err = s.Cache.EvalSha(context.Background(),
			cache.IntegralScript, []string{
				_const.UserKey,
				name,
				_const.NIntegralList,
				_const.YIntegralList,
			},
		).Result()
		if err != nil {
			s.Logger.Error("积分功能错误")
			return
		}
		switch {
		case res.(int64) == 0:
			// 用户已经抢过积分，不能再抢
		case res.(int64) == -1:
			// 积分已经抢完
		case res.(int64) > 0:
			// 抢到了积分
		default:
			// 发生了预期之外的错误
		}
	}
	return
}
