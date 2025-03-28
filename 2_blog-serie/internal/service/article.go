package service

import (
	"Golang_Programming_Journey/2_blog-serie/internal/dao"
	"Golang_Programming_Journey/2_blog-serie/internal/model"
	"Golang_Programming_Journey/2_blog-serie/pkg/app"
	"log"
)

//检验入参校验的结构体和绑定的姐结构体

type ArticleRequest struct {
	ID    uint32 `form:"id" binding:"required,gte=1"`
	State uint8  `form:"state,default=1" binding:"oneof=0 1"`
}

type ArticleListRequest struct {
	TagID uint32 `form:"tag_id" binding:"gte=1"`
	State uint8  `form:"state,default=1" binding:"oneof=0 1"`
}

type CreateArticleRequest struct {
	TagIDs        []uint32 `form:"tag_ids" binding:"required,gte=1"`
	Title         string   `form:"title" binding:"required,min=2,max=100"`
	Desc          string   `form:"desc" binding:"required,min=2,max=255"`
	Content       string   `form:"content" binding:"required,min=2,max=4294967295"`
	CoverImageUrl string   `form:"cover_image_url" binding:"required,url"`
	CreatedBy     string   `form:"created_by" binding:"required,min=2,max=100"`
	State         uint8    `form:"state,default=1" binding:"oneof=0 1"`
}

type UpdateArticleRequest struct {
	ID            uint32   `form:"id" binding:"required,gte=1"`
	TagIDs        []uint32 `form:"tag_ids" `
	Title         string   `form:"title" binding:"min=0,max=100"`
	Desc          string   `form:"desc" binding:"min=0,max=255"`
	Content       string   `form:"content" binding:"min=0,max=4294967295"`
	CoverImageUrl string   `form:"cover_image_url" `
	ModifiedBy    string   `form:"modified_by" binding:"required,min=2,max=100"`
	State         uint8    `form:"state,default=1" binding:"oneof=0 1"`
}

type DeleteArticleRequest struct {
	ID uint32 `form:"id" binding:"required,gte=1"`
}

type Article struct {
	ID            uint32       `json:"id"`
	Title         string       `json:"title"`
	Desc          string       `json:"desc"`
	Content       string       `json:"content"`
	CoverImageUrl string       `json:"cover_image_url"`
	State         uint8        `json:"state"`
	Tags          []*model.Tag `json:"tag"`
}

func (svc *Service) GetArticle(param *ArticleRequest) (*Article, error) {

	// 开始一个新的事务
	tx := svc.dao.Engine.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}

	// 确保在方法结束时处理事务
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback() // 发生错误时回滚事务
		}
	}()

	article, err := svc.dao.GetArticle(param.ID, param.State, tx)
	if err != nil {
		tx.Rollback() // 出现错误时回滚事务
		return nil, err
	}

	articleTags, err := svc.dao.GetArticleTagsByAID(article.ID, tx)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	var tags []*model.Tag
	for _, articleTag := range articleTags {
		tag, err := svc.dao.GetTag(articleTag.TagID, model.STATE_OPEN, tx)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
		tags = append(tags, &tag)
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	//// 返回处理后的结果
	//article.Tag = &tag

	return &Article{
		ID:            article.ID,
		Title:         article.Title,
		Desc:          article.Desc,
		Content:       article.Content,
		CoverImageUrl: article.CoverImageUrl,
		State:         article.State,
		Tags:          tags,
	}, nil

}

func (svc *Service) GetListArticle(param *ArticleListRequest, pager *app.Pager) ([]*Article, int, error) {

	//统计TagID 的数量
	articleCount, err := svc.dao.CountArticleListByTagID(param.TagID, param.State)
	log.Println(articleCount)
	if err != nil {
		return nil, 0, err
	}

	//这个地方 底层肯定需要并表查询 ，因为 返回的结果 既有AID 又有TID
	//通过 TagID把文章找出来
	articles, err := svc.dao.GetArticleListByTagID(param.TagID, param.State, pager.Page, pager.PageSize)
	if err != nil {
		return nil, 0, err
	}

	var articleList []*Article

	//通过list展示的articles list 的tag 只展示和他索引相关的tag信息
	for _, article := range articles {

		articleList = append(articleList, &Article{
			ID:            article.ArticleID,
			Title:         article.ArticleTitle,
			Desc:          article.ArticleDesc,
			Content:       article.Content,
			CoverImageUrl: article.CoverImageUrl,
			Tags: []*model.Tag{
				&model.Tag{
					Model: &model.Model{
						ID: article.TagID,
					},
					Name: article.TagName,
				},
			},
		})
	}
	return articleList, articleCount, nil
}

func (svc *Service) CreateArticle(param *CreateArticleRequest) error {
	tx := svc.dao.Engine.Begin() // 开始一个新的事务
	// 开始一个新的事务
	if tx.Error != nil {
		return tx.Error
	}

	article, err := svc.dao.CreateArticle(&dao.Article{
		Title:         param.Title,
		Desc:          param.Desc,
		Content:       param.Content,
		CoverImageUrl: param.CoverImageUrl,
		CreatedBy:     param.CreatedBy,
		State:         param.State,
	}, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	//新建了 一个 article，那么应该也要新建一下对应的 article 的id 和  tag 的tid对应表
	err = svc.dao.CreateArticleTags(article.ID, param.TagIDs, param.CreatedBy, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (svc *Service) UpdateArticle(param *UpdateArticleRequest) error {
	tx := svc.dao.Engine.Begin()
	if tx.Error != nil {
		return tx.Error
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			return
		}
	}()

	err := svc.dao.UpdateArticle(&dao.Article{
		ID:            param.ID,
		Title:         param.Title,
		Desc:          param.Desc,
		Content:       param.Content,
		CoverImageUrl: param.CoverImageUrl,
		State:         param.State,
		ModifiedBy:    param.ModifiedBy,
	}, tx)

	if err != nil {
		tx.Rollback()
		return err
	}
	//这个表 可能也许要修，因为他可能调整他所在的标签号
	err = svc.dao.UpdateArticleTags(param.ID, param.TagIDs, param.ModifiedBy, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (svc *Service) DeleteArticle(param *DeleteArticleRequest) error {
	tx := svc.dao.Engine.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			return
		}
	}()

	err := svc.dao.DeleteArticle(param.ID, tx)
	if err != nil {
		tx.Rollback()
		return err
	}
	err = svc.dao.DeleteArticleTag(param.ID, tx)
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}
