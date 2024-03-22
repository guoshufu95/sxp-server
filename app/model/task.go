package model

import (
	"gorm.io/gorm"
)

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
