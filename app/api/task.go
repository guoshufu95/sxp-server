package api

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	serv "sxp-server/app/service"
	"sxp-server/app/service/dto"
)

type TaskApi struct {
	Api
}

var ts = serv.TaskService{}

// List
//
//	@Description: 返回任务列表
//	@receiver a
//	@param c
func (a TaskApi) List(c *gin.Context) {
	a.BuildApi(c).BuildService(&ts.Service)
	err, tasks := ts.GetTaskList()
	if err != nil {
		a.ResponseError(http.StatusInternalServerError, err)
		return
	}
	a.Response("查询task列表成功", tasks)
}

// GetByParam
//
//	@Description: 条件查询
//	@receiver a
//	@param c
func (a TaskApi) GetByParam(c *gin.Context) {
	a.BuildApi(c).BuildService(&ts.Service)
	var req = dto.GetTasksByParamReq{}
	err := c.ShouldBindJSON(&req)
	if err != nil {
		a.ResponseError(http.StatusBadRequest, err)
		return
	}
	err, tasks := ts.QueryByParam(req)
	if err != nil {
		a.ResponseError(http.StatusInternalServerError, err)
		return
	}
	a.Response("查询成功!", tasks)
}

// GetById
//
//	@Description: 详情
//	@receiver a
//	@param c
func (a TaskApi) GetById(c *gin.Context) {
	a.BuildApi(c).BuildService(&ts.Service)
	var req = dto.GetTaskByIdParam{}
	err := c.ShouldBindJSON(&req)
	if err != nil {
		a.ResponseError(http.StatusBadRequest, err)
		return
	}
	err, task := ts.GetTaskById(req.Id)
	if err != nil {
		a.ResponseError(http.StatusInternalServerError, err)
		return
	}
	a.Response("查询成功!", task)
}

// CreateTask
//
//	@Description: 启动一个延时队列
//	@receiver a
//	@param c
func (a TaskApi) CreateTask(c *gin.Context) {
	a.BuildApi(c).BuildService(&ts.Service)
	var req = dto.StartTaskReq{}
	err := c.ShouldBindJSON(&req)
	if err != nil {
		a.ResponseError(http.StatusBadRequest, err)
		return
	}
	err, flag := ts.GetTaskByName(req.TaskName)
	if err != nil {
		a.ResponseError(http.StatusInternalServerError, err)
		return
	}
	if flag {
		err = errors.New("任务名重复")
		a.ResponseError(http.StatusInternalServerError, err)
		return
	}
	err = ts.SetTask(req)
	if err != nil {
		a.ResponseError(http.StatusInternalServerError, err)
		return
	}
	a.Response("设置定时任务成功!", nil)
}

// Update
//
//	@Description: 编辑更新
//	@receiver a
//	@param c
func (a TaskApi) Update(c *gin.Context) {
	a.BuildApi(c).BuildService(&ts.Service)
	var req dto.UpdateTaskReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		a.ResponseError(http.StatusBadRequest, err)
		return
	}
	err = ts.UpdateTask(req)
	if err != nil {
		a.ResponseError(http.StatusInternalServerError, err)
		return
	}
	a.Response("更新成功", nil)
}

// GetTasks
//
//	@Description: 获取定时任务队列
//	@receiver a
//	@param c
func (a TaskApi) GetTasks(c *gin.Context) {
	a.BuildApi(c).BuildService(&ts.Service)
	var req dto.GetTasksReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		a.ResponseError(http.StatusBadRequest, err)
		return
	}
	err, tasks := ts.GetTasks(req)
	if err != nil {
		a.ResponseError(http.StatusInternalServerError, err)
	}
	a.Response("success", tasks)
}

// DeleteTask
//
//	@Description: 删除任务
//	@receiver a
//	@param c
func (a TaskApi) DeleteTask(c *gin.Context) {
	a.BuildApi(c).BuildService(&ts.Service)
	var req dto.DelTaskReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		a.ResponseError(http.StatusBadRequest, err)
		return
	}
	err = ts.DeleteTask(req.Id)
	if err != nil {
		a.ResponseError(http.StatusInternalServerError, err)
	}
	a.Response("success", nil)
}
