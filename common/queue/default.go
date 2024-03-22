package queue

import (
	"fmt"
	"strconv"
	"sxp-server/common/initial"
)

var GlobalQueue *DelayQueue

// Operator
//
//	@Description: 模拟操作
//	@param i
//	@return flag
func Operator(i string) (flag bool) { //模拟业务场景
	fmt.Println("参数: ", i)
	num, _ := strconv.Atoi(i)
	if num/2 > 27 {
		flag = true
		return
	} else {
		return
	}
}

// StartQueue
//
//	@Description: 开启并消费队列
func StartQueue() {
	app := initial.App
	GlobalQueue = NewQueue(app.ProjectName, app.Cache, Operator)
	done := GlobalQueue.StartConsume()
	<-done
}
