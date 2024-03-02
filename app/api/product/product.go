package product

import (
	"github.com/gin-gonic/gin"
	"sxp-server/app/api"
	"sxp-server/app/model"
	serv "sxp-server/app/service"
	"sxp-server/app/service/product"
)

type ProductApi struct {
	api.Api
}

var ts = product.ProductService{}

func init() {
	serv.MakeService(&ts.Service)
}

// GetProduct
//
//	@Description: grpc调用获取产品
//	@receiver a
//	@param c
func (a *ProductApi) GetProduct(c *gin.Context) {
	a.MakeApi(c)
	var req = model.GetProductReq{}
	err := c.ShouldBindJSON(&req)
	if err != nil {
		a.Logger.Error(err.Error())
		a.ResponseError(err)
		return
	}
	err, res := ts.GetProduct(req.Id)
	if err != nil {
		a.Logger.Error("获取产品失败")
		a.ResponseError(err)
	}
	a.Response("success", res)
}

// UpdateProduct
//
//	@Description: 新建product
//	@receiver a
//	@param c
func (a *ProductApi) UpdateProduct(c *gin.Context) {
	a.MakeApi(c)
	var req = model.UpdateProductReq{}
	err := c.ShouldBindJSON(&req)
	if err != nil {
		a.Logger.Error(err.Error())
		a.ResponseError(err)
		return
	}
	err, res := ts.UpdateProduct(req)
	if err != nil {
		a.Logger.Error("获取产品失败")
		a.ResponseError(err)
	}
	a.Response("success", res)
}
