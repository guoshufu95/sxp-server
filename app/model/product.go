package model

type GetProductReq struct {
	Id string `json:"id"`
}

type GetProductRes struct {
	Product string `json:"product"`
}

type UpdateProductReq struct {
	ProductId string `json:"productId"`
	Product   string `json:"product"`
}

type UpdateProductRes struct {
	Message string `json:"message"`
}
