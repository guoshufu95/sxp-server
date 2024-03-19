package service

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"io"
	"strconv"
	"sxp-server/app/service/dto"
	"sxp-server/common/grpc/client"
	"sxp-server/common/grpc/helper"
	"sxp-server/common/grpc/pb"
	"sxp-server/common/logger"
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
func (ps *ProductGrpcService) GetProductById(ctx context.Context, id, token string) (err error, res *pb.ModelResponse) {
	log := logger.GetLogger()
	c := client.GetModelClient()
	var header, trailer metadata.MD
	response, err := c.GetModel(ctx,
		&pb.ModelRequest{
			ProductId: id,
		},
		grpc.Header(&header),   // 接收服务端发来的header
		grpc.Trailer(&trailer), // 接收服务端发来的trailer
	)
	if err != nil {
		log.Errorf("grpc服务调用失败: %s", err.Error())
		return
	}
	err, ok := helper.CheckTokenRes(header)
	if err != nil || !ok {
		log.Errorf("grpc服务端校验token失败:%s", err.Error())
		return
	}
	err, ok = helper.CheckSign(trailer)
	if err != nil || !ok {
		return
	}
	res = response
	return
}

// UpdateModel
//
//	@Description: 新建产品
//	@receiver ps
func (ps *ProductGrpcService) UpdateModel(req dto.UpdateProductReq, token string) (err error, res pb.UpdateResponse) {
	log := logger.GetLogger()
	c := client.GetModelClient()
	ctx := helper.BuildTokenCtx(token)
	response, err := c.UpdateModel(ctx, &pb.UpdateRequest{
		ProductId: strconv.Itoa(req.ProductId),
		Product:   req.Product,
	})
	if err != nil {
		log.Errorf("grpc服务调用失败: %s", err.Error())
		return
	}
	fmt.Println(response)
	return
}

// GetByStatus
//
//	@Description: 根据status状态获取产品信息
//	@receiver ps
//	@param status
//	@param token
//	@return err
//	@return response
func (ps *ProductGrpcService) GetByStatus(status, token string) (err error, response []*pb.StatusResponse) {
	ctx := helper.BuildTokenCtx(token)
	log := logger.GetLogger()
	c := client.GetModelClient()
	stream, err := c.GetByStatus(ctx)
	wch := make(chan struct{})
	if err != nil {
		log.Errorf("grpc服务调用失败: %s", err.Error())
		return
	}
	var ok bool
	header, _ := stream.Header()
	err, ok = helper.CheckTokenRes(header)
	if err != nil || !ok {
		log.Errorf("grpc服务端校验token失败:%s", err.Error())
		return
	}
	// 校验通过才会走到发送逻辑发送数据
	go func() {
		// 通过流发送消息
		err = stream.Send(&pb.StatusRequest{
			Status: status,
		})
		if err != nil {
			log.Info("发送流消息错误: %s", err.Error())
			return
		}
		_ = stream.CloseSend()
		return
	}()
	// 读取返回数据
	go func() {
		for {
			res, er := stream.Recv()
			if er != nil && er != io.EOF {
				log.Errorf("receive error:", er.Error())
			} else if er == io.EOF {
				log.Info("receive EOF")
				wch <- struct{}{}
				break
			}
			response = append(response, res)
		}
		return
	}()
	<-wch
	// 当RPC结束时读取trailer
	err, ok = helper.CheckSign(stream.Trailer())
	if err != nil || !ok {
		return
	}
	// todo 根据自己的业务逻辑处理
	return
}
