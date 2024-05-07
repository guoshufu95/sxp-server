package service

import (
	"fmt"
	"gorm.io/gorm"
	"strconv"
	"sxp-server/app/dao"
	"sxp-server/app/model"
	"sxp-server/app/service/dto"
	"sxp-server/common/timingWheel"
	"time"
)

type TaskService struct {
	Service
}

// GetTaskList
//
//	@Description: 获取所有的task
//	@receiver s
//	@return err
func (s *TaskService) GetTaskList() (err error, tasks []model.Task) {
	err, tasks = dao.GetAllTask(s.Db)
	if err != nil {
		s.Logger.Error("")
		return
	}
	for i, _ := range tasks {
		tasks[i].CreateTime = tasks[i].CreatedAt.Format(time.DateTime)
		tasks[i].ExecTime = time.Unix(tasks[i].Time, 0).Format(time.DateTime)
	}
	return
}

// QueryByParam
//
//	@Description: 条件查询
//	@receiver s
//	@param req
//	@return err
//	@return tasks
func (s *TaskService) QueryByParam(req dto.GetTasksByParamReq) (err error, tasks []model.Task) {
	db := buildCondition(s.Db, req)
	err = dao.QueryTasksByParam(db, &tasks)
	if err != nil {
		s.Logger.Error("查询tasks返回失败")
	}
	for i, _ := range tasks {
		tasks[i].CreateTime = tasks[i].CreatedAt.Format(time.DateTime)
		tasks[i].ExecTime = time.Unix(tasks[i].Time, 0).Format(time.DateTime)
	}
	return
}

// GetTaskById
//
//	@Description: 通过
//	@receiver s
//	@param id
//	@return err
//	@return task
func (s *TaskService) GetTaskById(id int) (err error, task model.Task) {
	err = dao.QueryTaskById(s.Db, id, &task)
	if err != nil {
		s.Logger.Error("通过id查询task详情失败")
		return
	}
	task.CreateTime = task.CreatedAt.Format(time.DateTime)
	task.ExecTime = time.Unix(task.Time, 0).Format(time.DateTime)
	return
}

// buildCondition
//
//	@Description: 构造查询条件
//	@param db
//	@param req
//	@return *gorm.DB
func buildCondition(db *gorm.DB, req dto.GetTasksByParamReq) *gorm.DB {
	if req.Name != "" {
		db = db.Where(fmt.Sprintf("task_name like \"%s\" or task_name like \"%s\" or task_name like \"%s\" or task_name =\"%s\"",
			"%"+req.Name+"%",
			"%"+req.Name,
			req.Name+"%",
			req.Name))
	}
	if req.Status == "成功" {
		db = db.Where("status = ?", 1)
	}
	if req.Status == "失败" {
		db = db.Where("status = ?", 2)
	}
	if req.Status == "未执行" {
		db = db.Where("status = ?", 0)
	}
	if req.Status == "执行中" {
		db = db.Where("status = ?", 3)
	}
	return db
}

// SetTask
//
//	@Description: 定时任务入库启动
//	@receiver s
//	@param req
//	@return err
func (s *TaskService) SetTask(req dto.StartTaskReq) (err error) {
	ss := strconv.Itoa(int(req.Time))
	ss = ss[0 : len(ss)-3]
	execTime, _ := strconv.Atoi(ss)
	var task model.Task
	if execTime == 0 || int64(execTime) <= time.Now().Unix() { //立即执行
		task.Time = time.Now().Unix()
	} else {
		task.Time = int64(execTime)
	}
	task.TaskName = req.TaskName
	retry, _ := strconv.Atoi(req.RetryTime)
	task.RetryTime = retry
	v, _ := strconv.Atoi(req.Value)
	task.Value = v
	err = dao.CreateTask(s.Db, task)
	if err != nil {
		s.Logger.Error("创建start定时任务入库失败")
		return
	}
	// 使用时间轮实现的延时队列
	var tw = timingWheel.ReturnTimingWheel()
	_, err = tw.CreateTask(task.TaskName, "productFn", time.Duration(task.Time), task.Value, retry)
	if err != nil {
		s.Logger.Error("创建任务队列失败")
		return
	}
	//go func() {
	//	if req.Time == 0 { //立即执行
	//		err = queue.GlobalQueue.SendScheduleMsg(model.TaskField{TaskName: req.TaskName, Value: strconv.Itoa(req.Value)}, time.Now(), queue.RetryCountOpt(req.RetryTime))
	//	} else { // 延时
	//		t := time.Unix(req.Time, 0).Sub(time.Now())
	//		err = queue.GlobalQueue.SendDelayMsg(model.TaskField{TaskName: req.TaskName, Value: strconv.Itoa(req.Value)}, t, req.RetryTime)
	//	}
	//	if err != nil {
	//		s.Logger.Errorf("定时任务执行失败:%s", err.Error())
	//		return
	//	}
	//}()
	return
}

// UpdateTask
//
//	@Description: 更新task
//	@receiver s
//	@param req
//	@return err
func (s *TaskService) UpdateTask(req dto.UpdateTaskReq) (err error) {
	var (
		execTime, value int
		task            = make(map[string]interface{})
	)
	if req.Time == 0 || req.Time < time.Now().Unix() { //立即执行
		task["time"] = time.Now().Unix()
	} else {
		ss := strconv.Itoa(int(req.Time))
		ss = ss[0 : len(ss)-3]
		execTime, _ = strconv.Atoi(ss)
		task["time"] = int64(execTime)
	}
	task["task_name"] = req.TaskName
	retry, _ := strconv.Atoi(req.RetryTime)
	task["retry_time"] = retry
	v, _ := strconv.Atoi(req.Value)
	task["value"] = v
	err = dao.UpdateTask(s.Db, task)
	if err != nil {
		s.Logger.Error("更新task失败")
		return
	}
	var tw = timingWheel.ReturnTimingWheel()
	_, err = tw.CreateTask(req.TaskName, "productFn", time.Duration(execTime), value, retry)
	if err != nil {
		s.Logger.Error("创建任务队列失败")
		return
	}
	return
}

// GetTaskByName
//
//	@Description: 查询任务信息service
//	@receiver s
//	@param name
//	@return err
//	@return flag
func (s *TaskService) GetTaskByName(name string) (err error, flag bool) {
	err, flag = dao.GetTaskFromName(s.Db, name)
	if err != nil {
		s.Logger.Error(err.Error())
	}
	return
}

// GetTasks
//
//	@Description: 根据入参获取定时任务列表
//	@receiver s
//	@return err
//	@return tasks
func (s *TaskService) GetTasks(req dto.GetTasksReq) (err error, tasks []model.Task) {
	err, tasks = dao.GetAllTask(s.Db)
	db := s.buildQuery(req.Name, req.Status)
	err, tasks = dao.GetAllTask(db)
	if err != nil {
		s.Logger.Error("获取定时任务列表失败")
		return
	}
	return
}

// DeleteTask
//
//	@Description: 删除任务
//	@receiver s
//	@param id
//	@return err
func (s *TaskService) DeleteTask(id int) (err error) {
	var task model.Task
	err = dao.GetTaskById(s.Db, id, &task)
	err = dao.DeleteTaskById(s.Db, task)
	return
}

// buildQuery
//
//	@Description: 构造查询字段
//	@receiver s
//	@param name
//	@param status
//	@return *gorm.DB
func (s *TaskService) buildQuery(name string, status int) *gorm.DB {
	db := s.Db.Or("task_name like ?", "%"+name+"%").Or("task_name like ?", name+"%").Or("task_name like ?", "%"+name).
		Or("task_name like ?", name).Where("status = ?", status)
	return db
}
