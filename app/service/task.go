package service

import (
	"gorm.io/gorm"
	"strconv"
	"sxp-server/app/dao"
	"sxp-server/app/model"
	"sxp-server/app/service/dto"
	"sxp-server/common/queue"
	"time"
)

type TaskService struct {
	Service
}

// SetTask
//
//	@Description: 定时任务入库启动
//	@receiver s
//	@param req
//	@return err
func (s *TaskService) SetTask(req dto.StartTaskReq) (err error) {
	err = dao.CreateTask(s.Db, req)
	if err != nil {
		s.Logger.Error("创建start定时任务入库失败")
		return
	}
	go func() {
		if req.Time == 0 { //立即执行
			err = queue.GlobalQueue.SendScheduleMsg(model.TaskField{TaskName: req.TaskName, Value: strconv.Itoa(req.Value)}, time.Now(), queue.RetryCountOpt(req.RetryTime))
		} else { // 延时
			t := time.Unix(req.Time, 0).Sub(time.Now())
			err = queue.GlobalQueue.SendDelayMsg(model.TaskField{TaskName: req.TaskName, Value: strconv.Itoa(req.Value)}, t, req.RetryTime)
		}
		if err != nil {
			s.Logger.Errorf("定时任务执行失败:%s", err.Error())
			return
		}
	}()
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
