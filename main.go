package main

import (
	"fmt"
	"os"
	"os/signal"
	"sxp-server/app/router"
	"sxp-server/common/grpc/client"
	g "sxp-server/common/grpc/client"
	ini "sxp-server/common/initial"
	"sxp-server/common/logger"
	"sxp-server/common/queue"
	"time"
)

func main() {
	fmt.Println("#############sxp项目启动中#############")
	//ini.Init() //初始化项目信息
	go SetUp()
	go queue.StartQueue() //开启延时队列
	// 等待中断信号以优雅地关闭服务器（设置 5 秒的超时时间）
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	queue.GlobalQueue.StopConsume()
	g.Stop()
	fmt.Printf("%s sxp服务停止 ... \r\n", time.Now().Format("2006-01-02 15:04:05"))
}

// SetUp
//
//	@Description: 启动操作
func SetUp() {
	l := logger.GetLogger()
	app := ini.App
	r := app.Engine
	router.InitRouter(r)
	client.Init() //初始化grpc客户端
	err := r.Run(":8000")
	if err != nil {
		l.Panicf("程序启动失败:%s", err.Error())
		return
	}
}
