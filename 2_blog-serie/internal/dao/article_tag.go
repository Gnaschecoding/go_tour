package dao

import (
	"Golang_Programming_Journey/2_blog-serie/internal/model"
	"gorm.io/gorm"
)

func (d *Dao) GetArticleTagsByAID(articleID uint32, tx *gorm.DB) ([]*model.ArticleTag, error) {
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

func (d *Dao) CreateArticleTags(articleID uint32, tagIDs []uint32, createdBy string, tx *gorm.DB) error {

	for _, tagID := range tagIDs {
		articleTag := model.ArticleTag{
			Model:     &model.Model{CreatedBy: createdBy},
			ArticleID: articleID,
			TagID:     tagID,
		}
		if err := articleTag.Create(tx); err != nil {
			return err
		}
	}

	return nil
}

// 更新标签 应该先删除后添加
func (d *Dao) UpdateArticleTags(articleID uint32, tagIDs []uint32, modifiedBy string, tx *gorm.DB) error {
	//articleTag := model.ArticleTag{ArticleID: articleID}

	if len(tagIDs) <= 0 {
		return nil
	}

	//先删除再插入
	if err := d.DeleteArticleTag(articleID, tx); err != nil {
		return err
	}

	if err := d.CreateArticleTags(articleID, tagIDs, modifiedBy, tx); err != nil {
		return err
	}

	//插入
	//for _, tagID := range tagIDs {
	//	if tagID <= 0 {
	//		continue
	//	}
	//	values := map[string]interface{}{
	//		"article_id":  articleID,
	//		"tag_id":      tagID,
	//		"modified_by": modifiedBy,
	//	}
	//
	//	if err := articleTag.UpdateOne(tx, values); err != nil {
	//		return err
	//	}
	//
	//}

	return nil
}

func (d *Dao) DeleteArticleTag(articleID uint32, tx *gorm.DB) error {
	articleTag := model.ArticleTag{ArticleID: articleID}
	return articleTag.DeleteByArticleId(tx)
}
