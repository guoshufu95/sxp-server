package client

import (
	"fmt"
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
	// 失败重试和链路追踪middleware
	conn, err := grpc.Dial(config.Conf.Grpc.Addr, grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(grpc_middleware.ChainUnaryClient(tracer.ClientUnaryInterceptor(trace),
			grpc_retry.UnaryClientInterceptor(retryOpts...))),
		grpc.WithStreamInterceptor(tracer.ClientStreamInterceptor(trace)),
	)
	//grpc.WithUnaryInterceptor(grpc_retry.UnaryClientInterceptor(retryOpts...)))
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
//	@Description: grpc停止
func Stop() {
	fmt.Println("grpc服务停止")
	err := grpcConn.Close()
	if err != nil {
		log.Fatalf("关闭grpcConn失败: %s", err.Error())
		return
	}
}
