package client

import (
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"sxp-server/common/grpc/pb"
	"sxp-server/common/tracer"
	"sxp-server/config"
	"time"
)

var (
	modelClient pb.ModelClient
	grpcConn    *grpc.ClientConn
)

// Init
//
//	@Description: 初始化grpc-client
func Init() (err error) {
	retryOpts := []grpc_retry.CallOption{
		// 最大重试次数
		grpc_retry.WithMax(uint(config.Conf.Grpc.Retry)),
		// 超时时间
		grpc_retry.WithPerRetryTimeout(time.Duration(config.Conf.Grpc.TimeOut) * time.Second),
		// 只有返回对应的code才会执行重试
		grpc_retry.WithCodes(codes.Unknown, codes.DeadlineExceeded, codes.Unavailable),
	}
	trace, _, err := tracer.NewJaegerTracer("sxp-server", config.Conf.Jaeger.Addr)
	if err != nil {
		return
	}
	conn, err := grpc.Dial(config.Conf.Grpc.Addr, grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(grpc_middleware.ChainUnaryClient(
			tracer.ClientUnaryInterceptor(trace),
			grpc_retry.UnaryClientInterceptor(retryOpts...))),
		grpc.WithStreamInterceptor(grpc_middleware.ChainStreamClient(
			tracer.ClientStreamInterceptor(trace),
		)))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	grpcConn = conn
	modelClient = pb.NewModelClient(grpcConn)
	return
}

func GetModelClient() pb.ModelClient {
	return modelClient
}

// Stop
//
//	@Description: 关闭grpc-client
func Stop() {
	err := grpcConn.Close()
	if err != nil {
		log.Fatalf("关闭grpcConn失败: %s", err.Error())
		return
	}
}
