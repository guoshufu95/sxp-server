package api

import (
	"errors"
	"github.com/gin-gonic/gin"
	serv "sxp-server/app/service"
	"sxp-server/app/service/dto"
)

type TaskApi struct {
	Api
}

var ts = serv.TaskService{}

// StartTask
//
//	@Description: 启动一个延时队列
//	@receiver a
//	@param c
func (a *TaskApi) StartTask(c *gin.Context) {
	a.MakeApi(c)
	var req = dto.StartTaskReq{}
	err := c.ShouldBindJSON(&req)
	if err != nil {
		a.Logger.Error(err.Error())
		a.ResponseError(err)
		return
	}
	serv.MakeService(&ts.Service, c)
	err, flag := ts.GetTaskByName(req.TaskName)
	if err != nil {
		err = errors.New("获取定时任务失败")
		a.ResponseError(err)
		return
	}
	if flag {
		err = errors.New("任务名重复")
		a.ResponseError(err)
		return
	}
	err = ts.SetTask(req)

	if err != nil {
		a.Logger.Error("设置定时任务失败!")
		a.ResponseError(err)
		return
	}
	a.Response("设置定时任务成功!", nil)
}

// GetTasks
//
//	@Description: 获取定时任务队列
//	@receiver a
//	@param c
func (a *TaskApi) GetTasks(c *gin.Context) {
	a.MakeApi(c)
	var req dto.GetTasksReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		a.Logger.Error(err.Error())
		a.ResponseError(err)
		return
	}
	serv.MakeService(&ts.Service, c)
	err, tasks := ts.GetTasks(req)
	if err != nil {
		err = errors.New("获取定时任务列表失败")
		a.ResponseError(err)
	}
	a.Response("success", tasks)
}
