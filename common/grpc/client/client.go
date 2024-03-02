package client

import (
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"sxp-server/common/grpc/pb"
)

var (
	modelClient pb.ModelClient
	addr        = "192.168.111.143:9001"
	grpcConn    *grpc.ClientConn
)

func Init() {
	// 连接到server端，此处禁用安全传输
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	grpcConn = conn
	//defer conn.Close()
	modelClient = pb.NewModelClient(grpcConn)
}

func GetModelClient() pb.ModelClient {
	return modelClient
}

func Stop() {
	fmt.Println("grpc服务停止")
	err := grpcConn.Close()
	if err != nil {
		log.Fatalf("关闭grpcConn失败: %s", err.Error())
		return
	}
}
