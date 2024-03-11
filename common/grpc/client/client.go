package client

import (
	"context"
	"fmt"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"sxp-server/common/grpc/pb"
	"sxp-server/common/tracer"
	"time"
)

var (
	modelClient   pb.ModelClient
	modelClienttt pb.ModelClient
	addr          = "192.168.111.40:9001"
	grpcConn      *grpc.ClientConn
)

// Init
//
//	@Description: 初始化grpc-client
func Init() (err error) {
	retryOpts := []grpc_retry.CallOption{
		// 最大重试次数
		grpc_retry.WithMax(3),
		// 超时时间
		grpc_retry.WithPerRetryTimeout(30 * time.Second),
		// 只有返回对应的code才会执行重试
		grpc_retry.WithCodes(codes.Unknown, codes.DeadlineExceeded, codes.Unavailable),
	}
	trace, _, err := tracer.NewJaegerTracer("sxp-server", "192.168.111.143:6831")
	if err != nil {
		return
	}
	// 失败重试和链路追踪middleware
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()),
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

// Inittt
//
//	@Description: 初始化grpc-client
func Inittt(ctx context.Context) (err error) {
	retryOpts := []grpc_retry.CallOption{
		// 最大重试次数
		grpc_retry.WithMax(3),
		// 超时时间
		grpc_retry.WithPerRetryTimeout(30 * time.Second),
		// 只有返回对应的code才会执行重试
		grpc_retry.WithCodes(codes.Unknown, codes.DeadlineExceeded, codes.Unavailable),
	}
	trace, _, err := tracer.NewJaegerTracer("sxp-server", "192.168.111.143:6831")
	if err != nil {
		return
	}
	// 失败重试和链路追踪middleware
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(grpc_middleware.ChainUnaryClient(tracer.ClientUnaryInterceptor(trace),
			grpc_retry.UnaryClientInterceptor(retryOpts...))),
		grpc.WithStreamInterceptor(tracer.ClientStreamInterceptor(trace)),
	)
	//grpc.WithUnaryInterceptor(grpc_retry.UnaryClientInterceptor(retryOpts...)))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	grpcConn = conn
	modelClienttt = pb.NewModelClient(grpcConn)
	return
}

func GetModelClient() pb.ModelClient {
	return modelClient
}

func GetModelClienttt() pb.ModelClient {
	return modelClienttt
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
