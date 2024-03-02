package cal

import (
	"fmt"
	"math/rand"
	"sxp-server/app/model"
	"time"
)

// GrabIntegral
//
//	@Description: 随机算法生成积分红包
//	@param cal
//	@return int
func GrabIntegral(integral *model.IntegralReq) int {
	if integral.RemainCount <= 0 {
		panic("RemainCount <= 0")
	}
	//最后一个
	if integral.RemainCount-1 == 0 {
		amount := integral.RemainIntegral
		integral.RemainCount = 0
		integral.RemainIntegral = 0
		return amount
	}
	//是否可以直接0.01
	if (integral.RemainIntegral / integral.RemainCount) == 1 {
		fmt.Println(integral.RemainIntegral / integral.RemainCount)
		amount := 1
		integral.RemainIntegral -= amount
		integral.RemainCount--
		return amount
	}

	//最大可领积分 = 剩余积分的平均值x2 = (剩余积分 / 剩余数量) * 2
	//领取积分范围 = 0.01 ~ 最大可领积分
	maxAmount := (integral.RemainIntegral / integral.RemainCount) * 2
	rand.Seed(time.Now().UnixNano())
	amount := rand.Intn(maxAmount)
	for amount == 0 {
		//防止零
		amount = rand.Intn(maxAmount)
	}
	integral.RemainIntegral -= amount
	//防止剩余积分负数
	if integral.RemainIntegral < 0 {
		amount += integral.RemainIntegral
		integral.RemainIntegral = 0
		integral.RemainCount = 0
	} else {
		integral.RemainCount--
	}
	return amount
}
