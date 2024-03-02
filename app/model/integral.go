package model

// IntegralReq
// @Description: 积分
type IntegralReq struct {
	Count          int `json:"count"`          //个数
	Integral       int `json:"cal"`            //总金额
	RemainCount    int `json:"remainCount"`    //剩余个数
	RemainIntegral int `json:"remainIntegral"` //剩余积分
	//BestIntegral      int   `json:"bestIntegral"`      //手气最佳金额
	//BestIntegralIndex int   `json:"bestIntegralIndex"` //手气最佳序号
	//IntegralList      []int `json:"integralList"`      //拆分列表
}
