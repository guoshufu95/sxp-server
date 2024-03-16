package api

import (
	"github.com/gin-gonic/gin"
	serv "sxp-server/app/service"
	"sxp-server/app/service/dto"
)

type ProductApi struct {
	Api
}

var ps = serv.ProductService{}

// GetProduct
//
//	@Description: grpc调用获取产品
//	@receiver a
//	@param c
func (a *ProductApi) GetProduct(c *gin.Context) {
	a.MakeApi(c)
	var req = dto.GetProductReq{}
	err := c.ShouldBindJSON(&req)
	if err != nil {
		a.Logger.Error(err.Error())
		a.ResponseError(err)
		return
	}
	serv.MakeService(&ts.Service, c)
	token, _ := c.Get("sxp-token")
	err, res := ps.GetProduct(req.Id, token.(string))
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
	var req = dto.UpdateProductReq{}
	err := c.ShouldBindJSON(&req)
	if err != nil {
		a.Logger.Error(err.Error())
		a.ResponseError(err)
		return
	}
	serv.MakeService(&ts.Service, c)
	token, _ := c.Get("sxp-token")
	err, res := ps.UpdateProduct(req, token.(string))
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
	var req = dto.GetByStatusReq{}
	err := c.ShouldBindJSON(&req)
	if err != nil {
		a.Logger.Error(err.Error())
		a.ResponseError(err)
		return
	}
	serv.MakeService(&ts.Service, c)
	token, _ := c.Get("sxp-token")
	err, res := ps.GetByStatus(req.Status, token.(string))
	if err != nil {
		a.Logger.Error("根据status获取产品失败")
		a.ResponseError(err)
		return
	}
	a.Response("success", res)
}
