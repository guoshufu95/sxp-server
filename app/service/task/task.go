package task

import (
	"gorm.io/gorm"
	"strconv"
	"sxp-server/app/dao/task"
	"sxp-server/app/model"
	"sxp-server/app/service"
	"sxp-server/common/queue"
	"time"
)

type TaskService struct {
	service.Service
}

// SetTask
//
//	@Description: 定时任务入库启动
//	@receiver s
//	@param req
//	@return err
func (s *TaskService) SetTask(req model.StartTaskReq) (err error) {
	err = task.CreateTask(s.Db, req)
	if err != nil {
		s.Logger.Error(err.Error())
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
	err, flag = task.GetTaskFromName(s.Db, name)
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
func (s *TaskService) GetTasks(req model.GetTasksReq) (err error, tasks []model.Task) {
	err, tasks = task.GetAllTask(s.Db)
	db := s.buildQuery(req.Name, req.Status)
	err, tasks = task.GetAllTask(db)
	if err != nil {
		s.Logger.Error(err.Error())
	}
	return
}

func (s *TaskService) buildQuery(name string, status int) *gorm.DB {
	db := s.Db.Or("task_name like ?", "%"+name+"%").Or("task_name like ?", name+"%").Or("task_name like ?", "%"+name).
		Or("task_name like ?", name).Where("status = ?", status)
	return db
}
