package main

import (
	"context"
	"fmt"
	"github.com/HelloMrShu/grpc_demo/componets"
	. "github.com/HelloMrShu/grpc_demo/global"
	pb "github.com/HelloMrShu/grpc_demo/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"strconv"
)

type greeterServer struct {
	pb.UnimplementedGreeterServer
}

// SayHello 简单实现一下.proto文件中定义的 SayHello 方法
func (g *greeterServer) SayHello(ctx context.Context, in *pb.HelloReq) (*pb.HelloResp, error) {
	log.Printf("Received: %v", in.GetName())
	return &pb.HelloResp{Message: "Hello " + in.GetName()}, nil
}

func main() {
	// 初始化
	Initialize()

	ServerConfig.Ip = componets.GetLocalIp().String()
	ServerConfig.Port, _ = componets.GetFreePort()
	//注册服务
	RegisterConsul()

	addr := ServerConfig.Ip + ":" + fmt.Sprintf("%d", ServerConfig.Port)
	fmt.Println(addr)
	listen, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	srv := grpc.NewServer()
	healthpb.RegisterHealthServer(srv, health.NewServer())

	pb.RegisterGreeterServer(srv, &greeterServer{})
	reflection.Register(srv)
	fmt.Println("Server start...")
	if err = srv.Serve(listen); err != nil {
		log.Println("Serving gRPC on 0.0.0.0" + strconv.Itoa(ServerConfig.Port))
		log.Fatalf("failed to serve: %v", err)
	}
	log.Println("Serving gRPC on " + addr)
}
