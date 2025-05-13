package main2

import (
	"flag"
	"github.com/soheilhy/cmux"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"net/http"
	pb "tag_service/proto"
	"tag_service/server"
)

// 2 不同协议用相同端口  cmux
var port2 string

func init() {
	flag.StringVar(&port2, "port2", "8003", "启动端口号")
	flag.Parse()
}

func main2() {
	lis, err := RunTCPServer2(port2)
	if err != nil {
		log.Fatalf("Run TCP Server err: %v", err)
	}
	mux := cmux.New(lis)

	grpcL := mux.MatchWithWriters(cmux.HTTP2MatchHeaderFieldPrefixSendSettings("content-type", "application/grpc"))
	httpL := mux.Match(cmux.HTTP1Fast())

	grpcS := RunGrpcServer2()
	httpS := RunHttpServer2(port2)

	go grpcS.Serve(grpcL)
	go httpS.Serve(httpL)

	if err := mux.Serve(); err != nil {
		log.Fatalf("Run TCP Server err: %v", err)
	}

}

func RunTCPServer2(port string) (net.Listener, error) {
	return net.Listen("tcp", ":"+port)
}

func RunHttpServer2(port string) *http.Server {
	serveMux := http.NewServeMux()
	serveMux.HandleFunc("/ping", func(writer http.ResponseWriter, request *http.Request) {
		_, _ = writer.Write([]byte(`pong`))
	})
	return &http.Server{
		Addr:    ":" + port,
		Handler: serveMux,
	}
}

func RunGrpcServer2() *grpc.Server {
	s := grpc.NewServer()
	pb.RegisterTagServiceServer(s, server.NewTagServer())
	reflection.Register(s)

	return s
}
