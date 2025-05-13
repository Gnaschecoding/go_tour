package server

import (
	"context"
	"encoding/json"
	"google.golang.org/grpc/metadata"
	"tag_service/pkg/bapi"
	"tag_service/pkg/errcode"
	pb "tag_service/proto"
)

type Auth struct{}

func (a *Auth) GetAppKey() string {
	return "go-programming-tour-book"
}

func (a *Auth) GetAppSecret() string {
	return "eddycjy"
}

func (a *Auth) Check(ctx context.Context) error {
	md, _ := metadata.FromIncomingContext(ctx)
	var appKey, appSecret string
	if value, ok := md["app_key"]; ok {
		appKey = value[0]
	}
	if value, ok := md["app_secret"]; ok {
		appSecret = value[0]
	}
	if appKey != a.GetAppKey() || appSecret != a.GetAppSecret() {
		return errcode.TogRPCError(errcode.Unauthorized)
	}
	return nil
}

// 重写了这个pb里面的
type TagServer struct {
	pb.UnimplementedTagServiceServer
	auth *Auth
}

func NewTagServer() *TagServer {
	return &TagServer{}
}

func (t *TagServer) GetTagList(ctx context.Context, r *pb.GetTagListRequest) (*pb.GetTagListReply, error) {

	//	panic("测试抛出异常！")
	//if err := t.auth.Check(ctx); err != nil {
	//	return nil, err
	//}

	api := bapi.NewAPI("http://127.0.0.1:8000")
	body, err := api.GetTagList(ctx, r.GetName())

	if err != nil {
		return nil, err
	}

	tagList := pb.GetTagListReply{}
	err = json.Unmarshal([]byte(body), &tagList)
	if err != nil {
		return nil, err
	}
	return &tagList, nil
}
