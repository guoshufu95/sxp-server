package dto

// GetProductReq
// @Description: id获取product入参
type GetProductReq struct {
	Id string `json:"id"`
}

// GetProductRes
// @Description: id获取response
type GetProductRes struct {
	Product string `json:"product"`
}

// UpdateProductReq
// @Description: 更新product入参
type UpdateProductReq struct {
	ProductId int    `json:"productId"`
	Product   string `json:"product"`
}

// UpdateProductRes
// @Description: 更新response
type UpdateProductRes struct {
	Message string `json:"message"`
}

// GetByStatusReq
// @Description: 通过status查询的入参
type GetByStatusReq struct {
	Status string `json:"status"`
}

// GetByStatusRes
// @Description: status查询的response
type GetByStatusRes struct {
	ProductId string `json:"productId"`
	Product   string `json:"product"`
	Status    string `json:"status"`
}
