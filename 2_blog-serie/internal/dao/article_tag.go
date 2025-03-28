package dao

import (
	"Golang_Programming_Journey/2_blog-serie/internal/model"
	"gorm.io/gorm"
)

func (d *Dao) GetArticleTagByAID(articleID uint32, tx *gorm.DB) (model.ArticleTag, error) {
	articleTag := model.ArticleTag{ArticleID: articleID}
	return articleTag.GetByAID(tx)
}

func (d *Dao) GetArticleTagListByTID(tagID uint32) ([]*model.ArticleTag, error) {
	articleTag := model.ArticleTag{TagID: tagID}
	return articleTag.ListByTID(d.Engine)
}

func (d *Dao) GetArticleTagListByAIDs(articleIDs []uint32) ([]*model.ArticleTag, error) {
	articleTag := model.ArticleTag{}
	return articleTag.ListByAIDs(d.Engine, articleIDs)

}

func (d *Dao) CreateArticleTag(articleID, tagID uint32, createdBy string, tx *gorm.DB) error {
	articleTag := model.ArticleTag{
		Model:     &model.Model{CreatedBy: createdBy},
		ArticleID: articleID,
		TagID:     tagID,
	}

	return articleTag.Create(tx)
}
func (d *Dao) UpdateArticleTag(articleID, tagID uint32, modifiedBy string, tx *gorm.DB) error {
	articleTag := model.ArticleTag{ArticleID: articleID}

	if tagID <= 0 {
		return nil
	}
	values := map[string]interface{}{
		"article_id":  articleID,
		"tag_id":      tagID,
		"modified_by": modifiedBy,
	}
	return articleTag.UpdateOne(tx, values)
}

func (d *Dao) DeleteArticleTag(articleID uint32, tx *gorm.DB) error {
	articleTag := model.ArticleTag{ArticleID: articleID}
	return articleTag.DeleteOne(tx)
}
