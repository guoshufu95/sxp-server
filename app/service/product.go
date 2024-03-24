package service

import (
	"sxp-server/app/service/dto"
	g "sxp-server/common/grpc/service"
)

type ProductService struct {
	Service
}

// GetProduct
//
//	@Description: 通过id查询
//	@receiver s
//	@param id
//	@return err
//	@return res
func (s *ProductService) GetProduct(id, token string) (err error, res dto.GetProductRes) {
	rpc := g.NewProductGrpcService(g.WithLog(s.Logger))
	err, val := rpc.GetProductById(id, token)
	if err != nil {
		s.Logger.Errorf("远程调用gprc服务失败: %s", err.Error())
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
func (s *ProductService) UpdateProduct(req dto.UpdateProductReq, token string) (err error, res dto.UpdateProductRes) {
	rpc := g.NewProductGrpcService(g.WithLog(s.Logger))
	err, val := rpc.UpdateModel(req, token)
	if err != nil {
		s.Logger.Errorf("远程调用gprc服务失败: %s", err.Error())
		return
	}
	res.Message = val.Message
	return
}

func (s *ProductService) GetByStatus(status, token string) (err error, res []dto.GetByStatusRes) {
	rpc := g.NewProductGrpcService(g.WithLog(s.Logger))
	err, val := rpc.GetByStatus(status, token)
	if err != nil {
		s.Logger.Errorf("远程调用gprc服务失败: %s", err.Error())
		return
	}
	for _, v := range val {
		item := dto.GetByStatusRes{
			ProductId: v.ProductId,
			Product:   v.Product,
			Status:    v.Status,
		}
		res = append(res, item)
	}
	return
}
