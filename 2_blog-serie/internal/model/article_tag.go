package model

import (
	"errors"
	"gorm.io/gorm"
)

type ArticleTag struct {
	*Model
	TagID     uint32 `json:"tag_id"`
	ArticleID uint32 `json:"article_id"`
}

func (a ArticleTag) TableName() string {
	return "blog_article_tag"
}

// /Session(&gorm.Session{}).
func (a ArticleTag) GetByAID(db *gorm.DB) ([]*ArticleTag, error) {
	var articleTags []*ArticleTag
	err := db.Where("article_id = ? AND is_del = ?", a.ArticleID, 0).Find(&articleTags).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return articleTags, err
	}
	return articleTags, nil
}

func (a ArticleTag) ListByTID(db *gorm.DB) ([]*ArticleTag, error) {
	var articleTags []*ArticleTag
	if err := db.Where("article_id = ? AND id_del = ?", a.ArticleID, 0).Find(&articleTags).Error; err != nil && err != gorm.ErrRecordNotFound {
		return articleTags, err
	}
	return articleTags, nil
}

func (a ArticleTag) ListByAIDs(db *gorm.DB, articleIDs []uint32) ([]*ArticleTag, error) {
	var articleTags []*ArticleTag
	if err := db.Where("article_id IN (?) AND is_del", articleIDs, 0).Find(&articleTags).Error; err != nil && err != gorm.ErrRecordNotFound {
		return articleTags, err
	}
	return articleTags, nil
}

func (a ArticleTag) Create(db *gorm.DB) error {

	if err := db.Create(&a).Error; err != nil {
		return err
	}
	return nil
}

func (a ArticleTag) UpdateOne(db *gorm.DB, values interface{}) error {
	db = db.Model(&a).Where("article_id = ? AND is_del = ?", a.ArticleID, 0)
	if err := db.Limit(1).Updates(values).Error; err != nil {
		return err
	}
	return nil
}

func (a ArticleTag) Delete(db *gorm.DB) error {
	if err := db.Where("id = ? AND is_del = ?", a.Model.ID, 0).Delete(&a).Error; err != nil {
		return err
	}
	return nil
}

func (a ArticleTag) DeleteByArticleId(db *gorm.DB) error {
	if err := db.Where("article_id = ? AND is_del = ?", a.ArticleID, 0).Delete(&a).Error; err != nil {
		return err
	}
	return nil
}
