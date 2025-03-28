package service

import (
	"Golang_Programming_Journey/2_blog-serie/internal/model"
	"Golang_Programming_Journey/2_blog-serie/pkg/app"
	"Golang_Programming_Journey/2_blog-serie/pkg/errcode"
)

//检验入参校验的结构体和绑定的姐结构体

type CountTagRequest struct {
	Name  string `form:"name" binding:"max=100"`
	State uint8  `form:"state,default=1" binding:"oneof=0 1"`
}

type TagListRequest struct {
	Name  string `form:"name" binding:"max=100"`
	State uint8  `form:"state,default=1" binding:"oneof=0 1"`
}

type CreateTagRequest struct {
	Name      string `form:"name" binding:"required,min=2,max=100"`
	CreatedBy string `form:"created_by" binding:"required,min=2,max=100"`
	State     uint8  `form:"state,default=1" binding:"oneof=0 1"`
}

type UpdateTagRequest struct {
	ID         uint32 `form:"id" binding:"required,gte=1"`
	Name       string `form:"name" binding:"max=100"`
	State      uint8  `form:"state" binding:"oneof=0 1"`
	ModifiedBy string `form:"modified_by" binding:"required,min=2,max=100"`
}

type DeleteTagRequest struct {
	ID uint32 `form:"id" binding:"required,gte=1"`
}

func (s *Service) CountTag(params *CountTagRequest) (int, error) {

	return s.dao.CountTag(params.Name, params.State)
}

func (s *Service) GetTagList(params *TagListRequest, pager *app.Pager) ([]*model.Tag, error) {

	return s.dao.GetTagList(params.Name, params.State, pager.Page, pager.PageSize)

}

func (s *Service) CreateTag(params *CreateTagRequest) error {
	// 检查标签是否已存在
	exists, err := s.dao.CheckTagExists(params.Name, params.State)
	if err != nil {
		return err
	}
	if exists {
		return errcode.ErrorCreateTagRepeatFail
	}

	return s.dao.CreateTag(params.Name, params.State, params.CreatedBy)
}

func (s *Service) UpdateTag(params *UpdateTagRequest) error {

	return s.dao.UpdateTag(params.ID, params.Name, params.State, params.ModifiedBy)
}

func (s *Service) DeleteTag(params *DeleteTagRequest) error {

	return s.dao.DeleteTag(params.ID)
}
