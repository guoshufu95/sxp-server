package main

import (
	"fmt"
	"os"
	"os/signal"
	"sxp-server/app/router"
	"sxp-server/common/grpc/client"
	g "sxp-server/common/grpc/client"
	ini "sxp-server/common/initial"
	"sxp-server/common/kafka"
	"sxp-server/common/logger"
	"sxp-server/common/timingWheel"
	"sxp-server/common/websocket"
	"sxp-server/config"
	"time"
)

func main() {
	l := logger.GetLogger()
	go SetUp()
	l.Info("########################### sxp项目启动中 #########################")
	go func() {
		timingWheel.StartTimingWheel()
		kafka.StartKafkaConsume() //开启kafka消费者
	}()
	//queue.StartQueue()      //开启延时队列
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit //优雅退出
	kafka.StopKafkaConsume()
	//queue.GlobalQueue.StopConsume()
	g.Stop()
	websocket.CloseSocket()
	ini.App.Cache.Close()
	l.Infof("%s sxp服务停止 ... \r\n", time.Now().Format("2006-01-02 15:04:05"))
}

// SetUp
//
//	@Description: 启动操作
func SetUp() {
	app := ini.App
	l := logger.GetLogger()
	r := app.Engine
	router.InitRouter(r)
	err := client.Init() //初始化grpc客户端
	if err != nil {
		l.Panicf("初始化grpc-client失败:%s", err.Error())
	}
	port := fmt.Sprintf(":%s", config.Conf.Server.Port)
	err = r.Run(port)
	if err != nil {
		l.Panicf("程序启动失败:%s", err.Error())
		return
	}
}
