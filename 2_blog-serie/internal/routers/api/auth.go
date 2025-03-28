package api

import (
	"Golang_Programming_Journey/2_blog-serie/global"
	"Golang_Programming_Journey/2_blog-serie/internal/service"
	"Golang_Programming_Journey/2_blog-serie/pkg/app"
	"Golang_Programming_Journey/2_blog-serie/pkg/errcode"
	"github.com/gin-gonic/gin"
)

func GetAuth(c *gin.Context) {
	param := service.AuthRequest{}
	response := app.NewResponse(c)

	ok, errs := app.BindAndValid(c, &param)

	if !ok {
		global.Logger.Errorf(c, "BindAndValid err:%v", errs)
		errRsp := errcode.InvalidParams.WithDetails(errs.Errors()...)
		response.ToErrorResponse(errRsp)
		return
	}

	svc := service.New(c.Request.Context())
	err := svc.CheckAuth(&param)

	if err != nil {
		global.Logger.Errorf(c, "CheckAuth err:%v", err)
		errRsp := errcode.InvalidParams.WithDetails(errs.Errors()...)
		response.ToErrorResponse(errRsp)
		return
	}

	token, err := app.GenerateToken(param.AppKey, param.AppSecret)
	if err != nil {
		global.Logger.Errorf(c, "GenerateToken err:%v", err)
		errRsp := errcode.InvalidParams.WithDetails(errs.Errors()...)
		response.ToErrorResponse(errRsp)
		return
	}
	response.ToResponse(gin.H{"token": token})
}
