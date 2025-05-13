package main1

import (
	"flag"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"net/http"
	pb "tag_service/proto"
	"tag_service/server"
)

// 1. 不同协议用不同端口
var grpcPort string
var httpPort string

func init() {
	flag.StringVar(&grpcPort, "grpc_port", "8001", "gRPC 启动端口号")
	flag.StringVar(&httpPort, "http_port", "9001", "HTTP 启动端口号")
	flag.Parse()
}

func main1() {
	errs := make(chan error)

	go func() {
		err := RunHttpServer1(httpPort)
		if err != nil {
			errs <- err
		}
	}()

	go func() {
		err := RunGrpcServer1(grpcPort)
		if err != nil {
			errs <- err
		}
	}()

	select {
	case err := <-errs:
		log.Fatalf("Run Server err: %v", err)
	}

}

func RunHttpServer1(port string) error {
	serveMux := http.NewServeMux()
	serveMux.HandleFunc("/ping", func(writer http.ResponseWriter, request *http.Request) {
		_, _ = writer.Write([]byte(`pong`))
	})
	return http.ListenAndServe(":"+port, serveMux)
}

func RunGrpcServer1(port string) error {
	s := grpc.NewServer()
	pb.RegisterTagServiceServer(s, server.NewTagServer())
	reflection.Register(s)

	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return err
	}
	return s.Serve(lis)
}
