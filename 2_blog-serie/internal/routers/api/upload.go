package api

import (
	"Golang_Programming_Journey/2_blog-serie/global"
	"Golang_Programming_Journey/2_blog-serie/internal/service"
	"Golang_Programming_Journey/2_blog-serie/pkg/app"
	"Golang_Programming_Journey/2_blog-serie/pkg/convert"
	"Golang_Programming_Journey/2_blog-serie/pkg/errcode"
	"Golang_Programming_Journey/2_blog-serie/pkg/upload"
	"github.com/gin-gonic/gin"
)

type Upload struct {
}

func NewUpload() Upload {
	return Upload{}
}

func (u Upload) UploadFile(c *gin.Context) {
	response := app.NewResponse(c)
	file, fileHeader, err := c.Request.FormFile("file")
	fileType := convert.StrTo(c.PostForm("type")).MustInt() //获取type字段，并且转为整数

	if err != nil {
		errRsp := errcode.InvalidParams.WithDetails(err.Error())
		response.ToErrorResponse(errRsp)
		return
	}

	if fileHeader == nil || fileType <= 0 {
		response.ToResponse(errcode.InvalidParams)
		return
	}

	svc := service.New(c.Request.Context())

	fileInfo, err := svc.UploadFile(upload.FileType(fileType), file, fileHeader)

	if err != nil {
		global.Logger.Errorf(c, "svc.UploadFile err: %v", err)
		errRsp := errcode.InvalidParams.WithDetails(err.Error())
		response.ToErrorResponse(errRsp)
		return
	}

	response.ToResponse(gin.H{
		"file_access_url": fileInfo.AccessUrl,
	})

}
