package dao

import (
	"gorm.io/gorm"
	"sxp-server/app/model"
	"sxp-server/app/service/dto"
	"time"
)

// GetAllTask
//
//	@Description: 查询所有task列表
//	@receiver s
//	@return err
func GetAllTask(db *gorm.DB) (err error, tasks []model.Task) {
	err = db.Find(&tasks).Error
	return
}

// GetTaskFromName
//
//	@Description:
//	@param db
//	@param taskName
//	@return err
//	@return flag
func GetTaskFromName(db *gorm.DB, taskName string) (err error, flag bool) {
	var count int64
	if err = db.Table("task").Where("task_name = ?", taskName).Count(&count).Error; err != nil {
		return
	}
	if count != 0 {
		flag = true
		return
	}
	return
}

// CreateTask
//
//	@Description: 创建task任务
//	@param db
//	@param req
//	@return err
func CreateTask(db *gorm.DB, req dto.StartTaskReq) (err error) {
	var task model.Task
	if req.Time == 0 {
		task.Time = time.Now().Unix()
	} else {
		task.Time = req.Time
	}
	task.TaskName = req.TaskName
	task.RetryTime = req.RetryTime
	task.Value = req.Value
	err = db.Create(&task).Error
	return
}
