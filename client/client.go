package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/resolver"
	"grpc_pro/etc"
	"grpc_pro/log"
	//"log"
	"strconv"
	"time"

	"grpc_pro/etcdv3"
	pb "grpc_pro/proto"
)

var (
	// EtcdEndpoints etcd地址
	EtcdEndpoints = []string{"localhost:2379"}
	// SerName 服务名称
	SerName    = "sky_grpc"
	grpcClient pb.SimpleClient
)

func main() {
	err := etc.InitConfig("./conf/app-template.yml")
	if err != nil {
		panic(err)
	}
	log.NewLogger(etc.Conf.LogInfo.LogPath, etc.Conf.LogInfo.LogLevel, etc.Conf.LogInfo.LogAdapter)

	r := etcdv3.NewServiceDiscovery(EtcdEndpoints)
	resolver.Register(r)
	// 连接服务器
	conn, err := grpc.Dial(
		fmt.Sprintf("%s:///%s", r.Scheme(), SerName),
		grpc.WithBalancerName("weight"),
		grpc.WithInsecure(),
	)
	if err != nil {
		log.Logger.Error("net.Connect err: %v", err)
		panic(err)
	}
	defer conn.Close()

	// 建立gRPC连接
	grpcClient = pb.NewSimpleClient(conn)
	start := time.Now().Unix()
	for i := 0; i < 100; i++ {
		route(i)
		//time.Sleep(1 * time.Second)
	}
	end := time.Now().Unix()
	fmt.Println(start, ",", end)
}

// route 调用服务端Route方法
func route(i int) {
	// 创建发送结构体
	req := pb.SimpleRequest{
		Data: "grpc " + strconv.Itoa(i),
	}
	// 调用我们的服务(Route方法)
	// 同时传入了一个 context.Context ，在有需要时可以让我们改变RPC的行为，比如超时/取消一个正在运行的RPC
	res, err := grpcClient.Route(context.Background(), &req)
	if err != nil {
		log.Logger.Error("Call Route err: %v", err)
		panic(err.Error())
	}
	// 打印返回值
	fmt.Println("code:",res.Code,",value:",res.Value)
}
