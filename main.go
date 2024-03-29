package main

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"sxp-server/app/router"
	"sxp-server/common/grpc/client"
	g "sxp-server/common/grpc/client"
	ini "sxp-server/common/initial"
	"sxp-server/common/logger"
	"sxp-server/common/queue"
	"sxp-server/config"
	"time"
)

func main() {
	fmt.Println("#############sxp项目启动中#############")
	go SetUp()
	go queue.StartQueue() //开启延时队列
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
	app := ini.App
	l := logger.GetLogger()
	r := app.Engine
	router.InitRouter(r)
	err := client.Init() //初始化grpc客户端
	if err != nil {
		l.Panicf("初始化grpc-client失败:%s", err.Error())
	}
	a := app.Engine.Routes()
	data := make([]string, 0)
	for _, v := range a {
		if strings.Contains(v.Path, "product") {
			data = append(data, v.Path)
		}
	}
	fmt.Println(data)
	port := fmt.Sprintf(":%s", config.Conf.Server.Port)
	err = r.Run(port)
	if err != nil {
		l.Panicf("程序启动失败:%s", err.Error())
		return
	}
}
