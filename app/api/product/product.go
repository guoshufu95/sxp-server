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
	serv.MakeService(&ts.Service, c)
	token, _ := c.Get("sxp-token")
	err, res := ts.GetProduct(req.Id, token.(string))
	if err != nil {
		a.Logger.Error("获取产品失败")
		a.ResponseError(err)
		return
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
	serv.MakeService(&ts.Service, c)
	token, _ := c.Get("sxp-token")
	err, res := ts.UpdateProduct(req, token.(string))
	if err != nil {
		a.Logger.Error("获取产品失败")
		a.ResponseError(err)
	}
	a.Response("success", res)
}

// GetByStatus
//
//	@Description: 根据在线状态获取产品信息
//	@receiver a
//	@param c
func (a *ProductApi) GetByStatus(c *gin.Context) {
	a.MakeApi(c)
	var req = model.GetByStatusReq{}
	err := c.ShouldBindJSON(&req)
	if err != nil {
		a.Logger.Error(err.Error())
		a.ResponseError(err)
		return
	}
	serv.MakeService(&ts.Service, c)
	token, _ := c.Get("sxp-token")
	err, res := ts.GetByStatus(req.Status, token.(string))
	if err != nil {
		a.Logger.Error("根据status获取产品失败")
		a.ResponseError(err)
		return
	}
	a.Response("success", res)
}
