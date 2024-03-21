package service

import (
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"sxp-server/app/dao"
	"sxp-server/app/service/dto"
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
		randomId := rand.Intn(100) // 随机数
		id := randomId + v
		//把积分初始化信息放入mysql
		err = dao.InsertIntegralInfo(s.Db, id, v)
		if err != nil {
			s.Logger.Error(err.Error())
			return
		}
		//放入未消费的积分队列(list),存放的是id
		err = dao.NConsumption(s.Cache, id)
		if err != nil {
			s.Logger.Error("存放未消费积分id到redis失败")
			return
		}
		//把积分相关字段放入redis，用hash表储存(id+money)
		err = dao.IntegralResolution(s.Cache, id, v)
		if err != nil {
			s.Logger.Error("存放未消费的id积分字典到redis失败")
			return
		}
	}
	return
}

// Do
//
//	@Description: 执行的业务逻辑
//	@receiver s
//	@param name
//	@return err
func (s *IntegralService) Do(name string) (err error, msg string, val int) {
	var (
		i, id int64
		res   string
	)
	//令牌桶
	err, i = dao.EvalLimit(s.Cache, []string{_const.Key1})
	if err != nil {
		s.Logger.Error("限流功能错误！")
		return
	}
	if i == 0 {
		s.Logger.Info("令牌已用完")
		return
	}
	if i > 0 { //令牌桶中有令牌
		// 抢积分的业务逻辑
		err, id = dao.DoIntegral(s.Cache, []string{_const.UserKey, name, _const.NIntegralList, _const.YIntegralList})
		if err != nil {
			s.Logger.Error("积分功能错误")
			return
		}
		switch {
		case id == 0:
			// 用户已经抢过积分，不能再抢
			val = int(id)
			msg = "您已经抢过积分了，请下次再来"
			s.Logger.Info("您已经抢过积分了，请下次再来")
			return
		case id == -1:
			// 积分已经抢完
			val = int(id)
			msg = "积分已经抢完了，请下次再来"
			s.Logger.Info("积分已经抢完了，请下次再来")
		case id > 0:
			// 抢到了积分
			err, res = dao.GetIntegralByIdFromCache(s.Cache, int(id))
			if err != nil {
				err = errors.New("获取积分值失败")
				return
			}
			val, _ = strconv.Atoi(res)
			msg = "抢到了积分"
			//todo 其他业务处理
			fmt.Println("抢到了积分")
		default:
			err = errors.New("预期之外的错误！")
			return
		}
	}
	return
}
