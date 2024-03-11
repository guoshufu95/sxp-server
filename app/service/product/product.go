package product

import (
	"sxp-server/app/model"
	"sxp-server/app/service"
	"sxp-server/common/grpc/helper"
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
func (s *ProductService) GetProduct(id, token string) (err error, res model.GetProductRes) {
	rpc := g.NewProductGrpcService()
	ctx := helper.BuildTokenCtx(token)
	err, val := rpc.GetProductById(ctx, id, token)
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
func (s *ProductService) UpdateProduct(req model.UpdateProductReq, token string) (err error, res model.UpdateProductRes) {
	rpc := g.NewProductGrpcService()
	err, val := rpc.UpdateModel(req, token)
	if err != nil {
		s.Logger.Error("远程调用gprc服务失败: %s", err.Error())
		return
	}
	res.Message = val.Message
	return
}

func (s *ProductService) GetByStatus(status, token string) (err error, res []model.GetByStatusRes) {
	rpc := g.NewProductGrpcService()
	err, val := rpc.GetByStatus(status, token)
	if err != nil {
		s.Logger.Errorf("远程调用gprc服务失败: %s", err.Error())
		return
	}
	for _, v := range val {
		item := model.GetByStatusRes{
			ProductId: v.ProductId,
			Product:   v.Product,
			Status:    v.Status,
		}
		res = append(res, item)
	}
	return
}
