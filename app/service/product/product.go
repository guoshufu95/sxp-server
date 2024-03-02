package product

import (
	"sxp-server/app/model"
	"sxp-server/app/service"
	g "sxp-server/common/grpc/service"
)

type ProductService struct {
	service.Service
}

// GetProduct
//
//	@Description: 通过id查询
//	@receiver s
//	@param id
//	@return err
//	@return res
func (s *ProductService) GetProduct(id string) (err error, res model.GetProductRes) {
	rpc := g.NewProductGrpcService()
	err, val := rpc.GetProductById(id)
	if err != nil {
		s.Logger.Error("远程调用gprc服务失败: %s", err.Error())
		return
	}
	v := val.GetProduct()
	res.Product = v
	return
}

// UpdateProduct
//
//	@Description: 新建产品
//	@receiver s
//	@param req
//	@return err
//	@return res
func (s *ProductService) UpdateProduct(req model.UpdateProductReq) (err error, res model.UpdateProductRes) {
	rpc := g.NewProductGrpcService()
	err, val := rpc.UpdateModel(req)
	if err != nil {
		s.Logger.Error("远程调用gprc服务失败: %s", err.Error())
		return
	}
	res.Message = val.Message
	return
}
