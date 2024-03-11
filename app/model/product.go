package model

type GetProductReq struct {
	Id string `json:"id"`
}

type GetProductRes struct {
	Product string `json:"product"`
}

type UpdateProductReq struct {
	ProductId int    `json:"productId"`
	Product   string `json:"product"`
}

type UpdateProductRes struct {
	Message string `json:"message"`
}

type GetByStatusReq struct {
	Status string `json:"status"`
}

type GetByStatusRes struct {
	ProductId string `json:"productId"`
	Product   string `json:"product"`
	Status    string `json:"status"`
}
