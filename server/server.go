package main

import (
	"context"
	"fmt"
	"grpc_pro/etc"
	"grpc_pro/log"
	"net"
	"strings"

	"google.golang.org/grpc"

	"grpc_pro/etcdv3"
	pb "grpc_pro/proto"
)

// MyService 定义我们的服务
type MyService struct{}

const (
	// Address 监听地址
	Address string = "127.0.0.1:8001"
	// Network 网络通信协议
	Network string = "tcp"
	// SerName 服务名称
	SerName string = "sky_grpc"
)

// EtcdEndpoints etcd地址
//var EtcdEndpoints = []string{"127.0.0.1:2379"}

func main() {
	err := etc.InitConfig("../conf/app-template.yml")
	if err != nil {
		panic(err)
	}
	log.NewLogger(etc.Conf.LogInfo.LogPath, etc.Conf.LogInfo.LogLevel, etc.Conf.LogInfo.LogAdapter)

	// 监听本地端口
	listener, err := net.Listen(Network, Address)
	if err != nil {
		panic(fmt.Sprintf("net.Listen err: %v", err))
	}
	log.Logger.Info(Address + " net.Listing...")
	// 新建gRPC服务器实例
	grpcServer := grpc.NewServer()
	// 在gRPC服务器注册我们的服务
	pb.RegisterSimpleServer(grpcServer, &MyService{})

	//把服务注册到etcd
	etcdEndpoints := strings.Split(etc.Conf.Server.EtcdAddr, ",") //从配置文件中获取etcd集群地址
	fmt.Println("etcd集群地址",etcdEndpoints)
	ser, err := etcdv3.NewServiceRegister(etcdEndpoints, SerName+"/"+Address, 1, 5)
	if err != nil {
		log.Logger.Error("register service err: %v", err)
		panic(err.Error())
	}
	defer ser.Close()
	//用服务器 Serve() 方法以及我们的端口信息区实现阻塞等待，直到进程被杀死或者 Stop() 被调用
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Logger.Error("grpcServer.Serve err: %v", err)
		panic(err.Error())
	}
}

// Route 实现Route方法
func (s *MyService) Route(ctx context.Context, req *pb.SimpleRequest) (*pb.SimpleResponse, error) {
	log.Logger.Info("接收到数据: " + req.Data)
	res := pb.SimpleResponse{
		Code:  200,
		Value: "你好:" + req.Data,
	}
	return &res, nil
}
