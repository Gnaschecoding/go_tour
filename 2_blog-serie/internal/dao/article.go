package dao

import (
	"Golang_Programming_Journey/2_blog-serie/internal/model"
	"Golang_Programming_Journey/2_blog-serie/pkg/app"
	"gorm.io/gorm"
)

type Article struct {
	ID            uint32 `json:"id"`
	TagID         uint32 `json:"tag_id"`
	Title         string `json:"title"`
	Desc          string `json:"desc"`
	Content       string `json:"content"`
	CoverImageUrl string `json:"cover_image_url"`
	CreatedBy     string `json:"created_by"`
	ModifiedBy    string `json:"modified_by"`
	State         uint8  `json:"state"`
}

//	func (d *Dao) GetArticle(id uint32, state uint8) (model.Article, error) {
//		article := model.Article{Model: &model.Model{ID: id}, State: state}
//		return article.Get(d.Engine)
//	}
func (d *Dao) GetArticle(id uint32, state uint8, tx *gorm.DB) (model.Article, error) {
	article := model.Article{Model: &model.Model{ID: id}, State: state}
	return article.Get(tx)
}

func (d *Dao) CreateArticle(param *Article, tx *gorm.DB) (*model.Article, error) {
	article := model.Article{
		Title:         param.Title,
		Desc:          param.Desc,
		Content:       param.Content,
		CoverImageUrl: param.CoverImageUrl,
		State:         param.State,
		Model:         &model.Model{CreatedBy: param.CreatedBy},
	}
	return article.Create(tx)

}

func (d *Dao) DeleteArticle(id uint32, tx *gorm.DB) error {
	article := model.Article{Model: &model.Model{ID: id}}
	return article.Delete(tx)
}

func (d *Dao) UpdateArticle(param *Article, tx *gorm.DB) error {
	article := &model.Article{Model: &model.Model{ID: param.ID}}
	values := map[string]interface{}{
		"modified_by": param.ModifiedBy,
		"state":       param.State,
	}

	// 仅当字段不为空时，才添加到更新字段集合，避免误将未初始化的字段更新为零值
	if param.Title != "" {
		values["title"] = param.Title
	}
	if param.Desc != "" {
		values["desc"] = param.Desc
	}
	if param.Content != "" {
		values["content"] = param.Content
	}
	if param.CoverImageUrl != "" {
		values["cover_image_url"] = param.CoverImageUrl
	}

	// 通过映射（map）动态更新字段，防止结构体字段的零值导致意外覆盖数据库中的已有数据
	return article.Update(tx, values)
}

func (d *Dao) CountArticleListByTagID(id uint32, state uint8) (int, error) {
	article := model.Article{State: state}
	return article.CountByTagID(d.Engine, id)
}

func (d *Dao) GetArticleListByTagID(id uint32, state uint8, page, pageSize int) ([]*model.ArticleRow, error) {
	article := model.Article{State: state}
	return article.ListByTagID(d.Engine, id, app.GetPageOffset(page, pageSize), pageSize)
}
