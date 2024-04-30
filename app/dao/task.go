package dao

import (
	"gorm.io/gorm"
	"sxp-server/app/model"
)

// GetAllTask
//
//	@Description: 查询所有task列表
//	@receiver s
//	@return err
func GetAllTask(db *gorm.DB) (err error, tasks []model.Task) {
	err = db.Debug().Find(&tasks).Error
	return
}

// QueryTasksByParam
//
//	@Description: 条件查询返回
//	@param db
//	@param tasks
//	@return err
func QueryTasksByParam(db *gorm.DB, tasks *[]model.Task) (err error) {
	err = db.Debug().Find(&tasks).Error
	return
}

// QueryTaskById
//
//	@Description: 通过id返回详情
//	@param db
//	@param id
//	@param task
//	@return err
func QueryTaskById(db *gorm.DB, id int, task *model.Task) (err error) {
	err = db.Debug().Where("id = ?", id).Find(&task).Error
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
func CreateTask(db *gorm.DB, task model.Task) (err error) {
	err = db.Debug().Create(&task).Error
	return
}

// UpdateTask
//
//	@Description:更新
//	@param db
//	@param task
//	@return err
func UpdateTask(db *gorm.DB, task map[string]interface{}) (err error) {
	err = db.Debug().Where("id = ?", task["task_name"]).Updates(&task).Error
	return
}

// UpdateTaskStatus
//
//	@Description: 更新任务执行状态
//	@param db
//	@param name
//	@param status
//	@return err
func UpdateTaskStatus(db *gorm.DB, name string, status int) (err error) {
	err = db.Model(&model.Task{}).Debug().Where("task_name = ?", name).Update("status", status).Error
	return
}
