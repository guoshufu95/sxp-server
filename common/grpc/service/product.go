package service

import (
	"context"
	"fmt"
	"sxp-server/app/model"
	"sxp-server/common/grpc/client"
	"sxp-server/common/grpc/pb"
	"sxp-server/common/logger"
	"time"
)

type ProductGrpcService struct {
	User   string
	Token  string
	RoleId int
}

type mOption func(service *ProductGrpcService)

func WithToken(token string) mOption {
	return func(s *ProductGrpcService) {
		s.Token = token
	}
}

func WithUser(name string) mOption {
	return func(s *ProductGrpcService) {
		s.User = name
	}
}

func WithRole(roleId int) mOption {
	return func(s *ProductGrpcService) {
		s.RoleId = roleId
	}
}

// NewProductGrpcService
//
//	@Description: 通过service对外暴露grpc调用方法
//	@return *ProductGrpcService
func NewProductGrpcService(opts ...mOption) *ProductGrpcService {
	s := &ProductGrpcService{}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

// GetProductById
//
//	@Description: 通过id查询产品
//	@receiver ps
func (ps *ProductGrpcService) GetProductById(id string) (err error, res *pb.ModelResponse) {
	log := logger.GetLogger()
	c := client.GetModelClient()
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	response, err := c.GetModel(ctx, &pb.ModelRequest{
		ProductId: id,
	})
	if err != nil {
		log.Errorf("grpc服务调用失败: %s", err.Error())
		return
	}
	fmt.Println(response)
	res = response
	return
}

// UpdateModel
//
//	@Description: 新建产品
//	@receiver ps
func (ps *ProductGrpcService) UpdateModel(req model.UpdateProductReq) (err error, res pb.UpdateResponse) {
	log := logger.GetLogger()
	c := client.GetModelClient()
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	response, err := c.UpdateModel(ctx, &pb.UpdateRequest{
		ProductId: req.ProductId,
		Product:   req.Product,
	})
	if err != nil {
		log.Errorf("grpc服务调用失败: %s", err.Error())
		return
	}
	fmt.Println(response)
	return
}
