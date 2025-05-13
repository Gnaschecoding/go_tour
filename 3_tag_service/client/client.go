package main

import (
	"context"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"log"
	"tag_service/global"
	"tag_service/internal/middleware"
	"tag_service/pkg/tracer"
	pb "tag_service/proto"
)

type Auth struct {
	AppKey    string
	AppSecret string
}

func (a *Auth) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	return map[string]string{"app_key": a.AppKey, "app_secret": a.AppSecret}, nil
}

func (a *Auth) RequireTransportSecurity() bool {
	return false
}

func init() {
	err := setupTracer()
	if err != nil {
		log.Fatalf("init.setupTracer err: %v", err)
	}
}

func main() {
	auth := Auth{
		AppKey:    "go-programming-tour-book",
		AppSecret: "eddycjy",
	}

	ctx := context.Background()
	opts := []grpc.DialOption{grpc.WithPerRPCCredentials(&auth)}
	opts = append(opts, grpc.WithUnaryInterceptor(
		grpc_middleware.ChainUnaryClient(
			middleware.UnaryContextTimeout(),
			middleware.ClientTracing(),
		),
	))
	newCtx := metadata.AppendToOutgoingContext(ctx, "eddycjy", "Go语言编程之旅")

	clientConn, err := GetClientConn(newCtx, "localhost:8004", opts)

	if err != nil {
		log.Fatalf("GetClientConn err: %v", err)
	}
	defer clientConn.Close()
	tagServiceClient := pb.NewTagServiceClient(clientConn)

	resp, err := tagServiceClient.GetTagList(ctx, &pb.GetTagListRequest{Name: "Golang"})
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("resp: %v", resp)
}

func GetClientConn(ctx context.Context, target string, opts []grpc.DialOption) (*grpc.ClientConn, error) {
	opts = append(opts, grpc.WithInsecure())
	return grpc.DialContext(ctx, target, opts...)
}

func setupTracer() error {
	var err error
	jaegerTracer, _, err := tracer.NewJaegerTracer("article-service", "127.0.0.1:6831")
	if err != nil {
		return err
	}
	global.Tracer = jaegerTracer
	return nil
}
