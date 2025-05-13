package main3

import (
	"context"
	"flag"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net/http"
	"strings"
	"tag_service"
	pb "tag_service/proto"
	"tag_service/server"
)

// 不同协议用相同端口  grpc_gate
var port3 string

func init() {
	flag.StringVar(&port3, "port3", "8004", "启动端口号")
	flag.Parse()
}

func main3() {
	err := main.RunServer(port3)
	if err != nil {
		log.Fatalf("Run Serve err: %v", err)
	}

}

func grpcHandlerFunc3(grpcServer *grpc.Server, otherHandler http.Handler) http.Handler {
	return h2c.NewHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.ProtoMajor == 2 && strings.Contains(r.Header.Get("Content-Type"), "application/grpc") {
			grpcServer.ServeHTTP(w, r)
		} else {
			otherHandler.ServeHTTP(w, r)
		}
	}), &http2.Server{})
}

func RunServer3(port string) error {
	httpMux := RunHttpServer3()
	grpcS := RunGrpcServer3()
	gatewayMux := runGrpcGatewayServer3()

	httpMux.Handle("/", gatewayMux)
	return http.ListenAndServe(":"+port, grpcHandlerFunc3(grpcS, httpMux))
}

func RunHttpServer3() *http.ServeMux {
	serveMux := http.NewServeMux()
	serveMux.HandleFunc("/ping", func(writer http.ResponseWriter, request *http.Request) {
		_, _ = writer.Write([]byte(`pong`))
	})
	return serveMux
}

func RunGrpcServer3() *grpc.Server {
	s := grpc.NewServer()
	pb.RegisterTagServiceServer(s, server.NewTagServer())
	reflection.Register(s)

	return s
}

func runGrpcGatewayServer3() *runtime.ServeMux {
	endpoint := "0.0.0.0:" + main.port
	gwmux := runtime.NewServeMux()
	dopts := []grpc.DialOption{grpc.WithInsecure()}
	_ = pb.RegisterTagServiceHandlerFromEndpoint(context.Background(), gwmux, endpoint, dopts)
	return gwmux
}
