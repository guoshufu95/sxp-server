package model

import (
	"gorm.io/gorm"
)

// StartTaskReq 设置延时任务队列的请求参数
type StartTaskReq struct {
	TaskName  string `json:"taskName"  binding:"required" msg:"task名不能为空"`
	Value     int    `json:"value"  binding:"required" msg:"value不能为空"`
	Time      int64  `json:"time"`
	RetryTime int    `json:"retryTime"`
}

// GetTasksReq
// @Description: 通过name查询入参
type GetTasksReq struct {
	Name   string `json:"name"`
	Status int    `json:"status"`
}

// Task
// @Description: 数据库字段
type Task struct {
	gorm.Model
	TaskName  string `json:"taskName" ` //任务名称
	Value     int    `json:"value" `    //值
	Time      int64  `json:"time" `     //执行时间
	RetryTime int    `json:"retryTime"` //重试次数
	Status    int    `json:"status"`    // 执行状态
}

func (Task) TableName() string {
	return "task"
}

func (t *Task) BeforeCreate(_ *gorm.DB) error {
	t.TaskName = t.TaskName + "???"
	return nil
}

// TaskField
// @Description: redis储存字段
type TaskField struct {
	TaskName string `json:"taskName" `
	Value    string `json:"value" `
}
