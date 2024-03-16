package dto

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
