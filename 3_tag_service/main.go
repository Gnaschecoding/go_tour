package main

import (
	"context"
	"flag"
	assetfs "github.com/elazarl/go-bindata-assetfs"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net/http"
	"path"
	"strings"
	"tag_service/global"
	"tag_service/internal/middleware"
	"tag_service/pkg/swagger"
	"tag_service/pkg/tracer"
	pb "tag_service/proto"
	"tag_service/server"
)

// 不同协议用相同端口  cmux
var port string

func init() {
	flag.StringVar(&port, "port", "8004", "启动端口号")
	flag.Parse()

	err := setupTracer()
	if err != nil {
		log.Fatalf("init.setupTracer err: %v", err)
	}

}

func setupTracer() error {
	jaegerTracer, _, err := tracer.NewJaegerTracer("tour-service", "127.0.0.1:6831")
	if err != nil {
		return err
	}
	global.Tracer = jaegerTracer
	return nil
}

func main() {
	err := RunServer(port)
	if err != nil {
		log.Fatalf("Run Serve err: %v", err)
	}

}

func grpcHandlerFunc(grpcServer *grpc.Server, otherHandler http.Handler) http.Handler {
	return h2c.NewHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.ProtoMajor == 2 && strings.Contains(r.Header.Get("Content-Type"), "application/grpc") {
			grpcServer.ServeHTTP(w, r)
		} else {
			otherHandler.ServeHTTP(w, r)
		}
	}), &http2.Server{})
}

func RunServer(port string) error {
	httpMux := RunHttpServer()
	grpcS := RunGrpcServer()
	gatewayMux := runGrpcGatewayServer()

	httpMux.Handle("/", gatewayMux)
	return http.ListenAndServe(":"+port, grpcHandlerFunc(grpcS, httpMux))
}

func RunHttpServer() *http.ServeMux {

	serveMux := http.NewServeMux()
	serveMux.HandleFunc("/ping", func(writer http.ResponseWriter, request *http.Request) {
		_, _ = writer.Write([]byte(`pong`))
	})

	prefix := "/swagger-ui/"
	fileServer := http.FileServer(&assetfs.AssetFS{
		Asset:    swagger.Asset,
		AssetDir: swagger.AssetDir,
		Prefix:   "third_party/swagger-ui",
	})
	serveMux.Handle(prefix, http.StripPrefix(prefix, fileServer))

	serveMux.HandleFunc("/swagger/", func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasSuffix(r.URL.Path, "swagger.json") {
			http.NotFound(w, r)
			return
		}

		p := strings.TrimPrefix(r.URL.Path, "/swagger/")
		p = path.Join("proto", p)

		http.ServeFile(w, r, p)
	})
	return serveMux
}

func RunGrpcServer() *grpc.Server {

	//grpc需要拦截器，做一些非核心业务 比如计时啊，比如RPC 方法进行鉴权校验，
	//RPC 方法进行上下文的超时控制每个 RPC 方法的请求都做日志记录
	options := []grpc.ServerOption{
		//拦截器在前面的先执行
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			//HelloInterceptor,
			//WorldInterceptor,
			middleware.AccessLog,
			middleware.ErrorLog,
			middleware.Recovery,
			middleware.ServerTracing,
		)),
	}
	s := grpc.NewServer(options...)
	pb.RegisterTagServiceServer(s, server.NewTagServer())
	reflection.Register(s)

	return s
}

func runGrpcGatewayServer() *runtime.ServeMux {
	endpoint := "0.0.0.0:" + port
	gwmux := runtime.NewServeMux()
	dopts := []grpc.DialOption{grpc.WithInsecure()}
	_ = pb.RegisterTagServiceHandlerFromEndpoint(context.Background(), gwmux, endpoint, dopts)
	return gwmux
}

func HelloInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	log.Println("你好，红烧煎鱼")
	resp, err := handler(ctx, req)
	log.Println("再见，红烧煎鱼")
	return resp, err
}
func WorldInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	log.Println("你好，清蒸煎鱼")
	resp, err := handler(ctx, req)
	log.Println("再见，清蒸煎鱼")
	return resp, err
}
