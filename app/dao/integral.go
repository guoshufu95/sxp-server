package dao

import (
	"context"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"strconv"
	"sxp-server/app/model"
	"sxp-server/common/cache"
	_const "sxp-server/common/const"
	"time"
)

// InsertIntegralInfo
//
//	@Description: 积分初始化字段入库
//	@param db
//	@param id
//	@param integral
//	@return err
func InsertIntegralInfo(db *gorm.DB, id, integral int) (err error) {
	var param = model.IntegralInfo{
		IntegralsId: id,
		Integrals:   integral,
	}
	return db.Debug().Create(&param).Error
}

// NConsumption
//
//	@Description: 未消费的积分id
//	@param rdb
//	@param id
//	@return err
func NConsumption(rdb *redis.Client, id int) (err error) {
	_, err = rdb.LPush(context.Background(), _const.NIntegralList, id).Result()
	return
}

// IntegralResolution
//
//	@Description: 拆分后的积分id+积分值
//	@param rdb
//	@param id
//	@param val
//	@return err
func IntegralResolution(rdb *redis.Client, id, val int) (err error) {
	_, err = rdb.HSet(context.Background(), _const.IntegralInfo, id, val).Result()
	return
}

// EvalLimit
//
//	@Description: 限流
//	@param rdb
//	@param keys
//	@return err
func EvalLimit(rdb *redis.Client, keys []string) (error, int64) {
	res, err := rdb.Eval(context.Background(), cache.FilterScript, keys, time.Now().UnixMilli()).Result()
	i := res.(int64)
	return err, i
}

// DoIntegral
//
//	@Description: 抢积分的逻辑执行
//	@param rdb
//	@param params
//	@return err
//	@return i
func DoIntegral(rdb *redis.Client, params []string) (error, int64) {
	res, err := rdb.Eval(context.Background(), cache.IntegralScript, params).Result()
	i := res.(int64)
	return err, i
}

// GetIntegralByIdFromCache
//
//	@Description: 获取积分值
//	@param rdb
//	@param id
//	@return err
//	@return val
func GetIntegralByIdFromCache(rdb *redis.Client, id int) (err error, val string) {
	val, err = rdb.HGet(context.Background(), _const.IntegralInfo, strconv.Itoa(id)).Result()
	return
}
