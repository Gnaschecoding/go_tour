package service

import (
	"Golang_Programming_Journey/2_blog-serie/global"
	"Golang_Programming_Journey/2_blog-serie/internal/dao"
	"Golang_Programming_Journey/2_blog-serie/pkg/otgorm"
	"context"
)

type Service struct {
	ctx context.Context
	dao *dao.Dao
}

func New(ctx context.Context) *Service {
	svc := &Service{ctx: ctx}
	svc.dao = dao.NewDao(otgorm.WithContext(svc.ctx, global.DBEngine))
	return svc
}
