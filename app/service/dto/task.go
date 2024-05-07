package dto

// StartTaskReq 设置延时任务队列的请求参数
type StartTaskReq struct {
	TaskName  string `json:"name"  binding:"required" msg:"task名不能为空"`
	Value     string `json:"value"  binding:"required" msg:"value不能为空"`
	Time      int64  `json:"execTime"`
	RetryTime string `json:"retryTime"`
}

// GetTasksReq
// @Description: 通过name查询入参
type GetTasksReq struct {
	Name   string `json:"name"`
	Status int    `json:"status"`
}

// DelTaskReq
// @Description: 刪除task入参
type DelTaskReq struct {
	Id int `json:"id"`
}

// GetTasksByParamReq
// @Description: 条件查询参数
type GetTasksByParamReq struct {
	Name   string `json:"name"`
	Status string `json:"status"`
}

// GetTaskByIdParam
// @Description: 详情入参
type GetTaskByIdParam struct {
	Id int `json:"id"`
}

// UpdateTaskReq
// @Description: 更新
type UpdateTaskReq struct {
	TaskName  string `json:"name"  binding:"required" msg:"task名不能为空"`
	Value     string `json:"value"  binding:"required" msg:"value不能为空"`
	Time      int64  `json:"execTime"`
	RetryTime string `json:"retryTime"`
}
