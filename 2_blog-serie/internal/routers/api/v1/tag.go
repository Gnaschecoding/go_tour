package v1

import (
	"Golang_Programming_Journey/2_blog-serie/global"
	"Golang_Programming_Journey/2_blog-serie/internal/service"
	"Golang_Programming_Journey/2_blog-serie/pkg/app"
	"Golang_Programming_Journey/2_blog-serie/pkg/convert"
	"Golang_Programming_Journey/2_blog-serie/pkg/errcode"
	"errors"
	"github.com/gin-gonic/gin"
)

type Tag struct {
}

func NewTag() Tag {
	return Tag{}
}

// List @Summary 获取多个标签
// @Produce json
// @Param name query string false "标签名称" maxlength(100)
// @Param state query int false "状态" Enums(0, 1) default(1)
// @Param page query int false "页码"
// @Param page_size query int false "每页数量"
// @Success 200 {object} model.TagSwagger "成功"
// @Failure 400 {object} errcode.Error "请求错误"
// @Failure 500 {object} errcode.Error "内部错误"
// @Router /api/v1/tags [get]
func (t Tag) List(c *gin.Context) {
	//time.Sleep(4 * time.Second)
	params := service.TagListRequest{}
	response := app.NewResponse(c)

	ok, errs := app.BindAndValid(c, &params)
	if !ok {
		global.Logger.Errorf(c, "app.BindAndValid err:%v", errs)
		errRsp := errcode.InvalidParams.WithDetails(errs.Errors()...)
		response.ToErrorResponse(errRsp)
		return
	}

	s := service.New(c.Request.Context())
	pager := app.Pager{
		Page:     app.GetPage(c),
		PageSize: app.GetPageSize(c),
	}

	totalRows, err := s.CountTag(&service.CountTagRequest{
		Name:  params.Name,
		State: params.State,
	})
	if err != nil {
		global.Logger.Errorf(c, "s.CountTag err:%v", err)
		response.ToErrorResponse(errcode.ErrorCountTagFail)
		return
	}

	tags, err := s.GetTagList(&params, &pager)
	if err != nil {
		global.Logger.Errorf(c, "s.GetListTag err:%v", err)
		response.ToErrorResponse(errcode.ErrorGetTagListFail)
		return
	}

	response.ToResponseList(tags, totalRows)
	return
}

// Create @Summary 新增标签
// @Produce json
// @Param name body string true "标签名称" minlength(3) maxlength(100)
// @Param state body int false "状态" Enums(0, 1) default(1)
// @Param created_by body string true "创建者" minlength(3) maxlength(100)
// @Success 200 {object} model.TagSwagger "成功"
// @Failure 400 {object} errcode.Error "请求错误"
// @Failure 500 {object} errcode.Error "内部错误"
// @Router /api/v1/tags [post]
func (t Tag) Create(c *gin.Context) {
	params := service.CreateTagRequest{}
	response := app.NewResponse(c)
	ok, errs := app.BindAndValid(c, &params)
	if !ok {
		global.Logger.Errorf(c, "app.BindAndValid err:%v", errs)
		errRsp := errcode.InvalidParams.WithDetails(errs.Errors()...)
		response.ToErrorResponse(errRsp)
		return
	}
	ser := service.New(c.Request.Context())

	err := ser.CreateTag(&params)
	if err != nil {
		var customErr *errcode.Error
		if errors.As(err, &customErr) && errors.Is(customErr, errcode.ErrorCreateTagRepeatFail) {
			global.Logger.Errorf(c, "service.CreateTag repeat err:%v", err)
			response.ToErrorResponse(errcode.ErrorCreateTagRepeatFail)
			return
		}

		global.Logger.Errorf(c, "service.CreateTag err:%v", err)
		response.ToErrorResponse(errcode.ErrorCreateTagFail)
		return
	}

	response.ToResponse(gin.H{})
	return

}

// Update @Summary 更新标签
// @Produce json
// @Param id path int true "标签ID"
// @Param name body string false "标签名称" minlength(3) maxlength(100)
// @Param state body int false "状态" Enums(0, 1) default(1)
// @Param modified_by body string true "修改者" minlength(3) maxlength(100)
// @Success 200 {array} model.TagSwagger "成功"
// @Failure 400 {object} errcode.Error "请求错误"
// @Failure 500 {object} errcode.Error "内部错误"
// @Router /api/v1/tags/{id} [put]
func (t Tag) Update(c *gin.Context) {
	params := service.UpdateTagRequest{
		ID: convert.StrTo(c.Param("id")).MustUInt32(), //先强制类型转为 StrTo，然后调用 其MustUInt32 方法
	}
	response := app.NewResponse(c)
	ok, errs := app.BindAndValid(c, &params)
	if !ok {
		global.Logger.Errorf(c, "app.BindAndValid err:%v", errs)
		errRsp := errcode.InvalidParams.WithDetails(errs.Errors()...)
		response.ToErrorResponse(errRsp)
		return
	}

	ser := service.New(c.Request.Context())
	err := ser.UpdateTag(&params)
	if err != nil {
		global.Logger.Errorf(c, "service.UpdateTag err:%v", err)
		response.ToErrorResponse(errcode.ErrorUpdateTagFail)
		return
	}
	response.ToResponse(gin.H{})
	return

}

// Delete @Summary 删除标签
// @Produce json
// @Param id path int true "标签ID"
// @Success 200 {string} string "成功"
// @Failure 400 {object} errcode.Error "请求错误"
// @Failure 500 {object} errcode.Error "内部错误"
// @Router /api/v1/tags/{id} [delete]
func (t Tag) Delete(c *gin.Context) {
	params := service.DeleteTagRequest{
		ID: convert.StrTo(c.Param("id")).MustUInt32(),
	}
	response := app.NewResponse(c)
	ok, errs := app.BindAndValid(c, &params)
	if !ok {
		global.Logger.Errorf(c, "app.BindAndValid err:%v", errs)
		errRsp := errcode.InvalidParams.WithDetails(errs.Errors()...)
		response.ToErrorResponse(errRsp)
		return
	}

	ser := service.New(c.Request.Context())
	err := ser.DeleteTag(&params)
	if err != nil {
		global.Logger.Errorf(c, "service.DeleteTag err:%v", err)
		response.ToErrorResponse(errcode.ErrorDeleteTagFail)
		return
	}

	response.ToResponse(gin.H{})

}
