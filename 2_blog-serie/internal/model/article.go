package model

import (
	"Golang_Programming_Journey/2_blog-serie/pkg/app"
	"gorm.io/gorm"
)

type Article struct {
	*Model
	Title         string `json:"title"`
	Desc          string `json:"desc"`
	Content       string `json:"content"`
	CoverImageUrl string `json:"cover_image_url"`
	State         uint8  `json:"state"`
}

type ArticleSwagger struct {
	List  []*Article
	Pager *app.Pager
}

func (a Article) TableName() string {
	return "blog_article"
}

// Session(&gorm.Session{})
func (a Article) Get(db *gorm.DB) (Article, error) {
	article := &Article{}

	err := db.Where("id = ? AND state = ? AND is_del = ?", a.ID, a.State, 0).First(article).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return *article, err
	}
	return *article, nil
}

func (a Article) Create(db *gorm.DB) (*Article, error) {
	// 添加日志信息，确认方法调用
	//log.Printf("Calling Article.Create with Title: %s\n", a.Title)

	if err := db.Create(&a).Error; err != nil {
		return nil, err
	}
	//log.Printf("Calling Article.Create with id: %d\n", a.ID)
	return &a, nil
}

func (a Article) Update(db *gorm.DB, values interface{}) error {
	db = db.Model(&a).Where("id = ? AND is_del = ?", a.ID, 0)
	err := db.Updates(values).Error
	if err != nil {
		return err
	}
	return nil
}

func (a Article) Delete(db *gorm.DB) error {
	db = db.Where("id = ? AND is_del = ?", a.Model.ID, 0)
	err := db.Delete(&a).Error
	if err != nil {
		return err
	}
	return nil
}

// 如果想要获取这个表中的数据 这里需要并表查询了
type ArticleRow struct {
	ArticleID     uint32
	TagID         uint32
	TagName       string
	ArticleTitle  string
	ArticleDesc   string
	CoverImageUrl string
	Content       string
}

func (a Article) ListByTagID(db *gorm.DB, tagID uint32, pageOffset, pageSize int) ([]*ArticleRow, error) {
	fields := []string{"ar.id AS article_id", "ar.title AS article_title", "ar.desc AS article_desc", "ar.cover_image_url", "ar.content"}
	fields = append(fields, []string{"t.id AS tag_id", "t.name AS tag_name"}...)
	if pageOffset >= 0 && pageSize > 0 {
		db = db.Offset(pageOffset).Limit(pageSize)
	}

	rows, err := db.Session(&gorm.Session{}).Select(fields).Table(ArticleTag{}.TableName()+" AS at").
		Joins("LEFT JOIN `"+Tag{}.TableName()+"` AS t ON at.tag_id = t.id").
		Joins("LEFT JOIN `"+Article{}.TableName()+"` AS ar ON at.article_id = ar.id").
		Where("at.`tag_id` = ? AND ar.state = ? AND ar.is_del = ?", tagID, a.State, 0).
		Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var articles []*ArticleRow
	for rows.Next() {
		r := &ArticleRow{}
		if err := rows.Scan(&r.ArticleID, &r.ArticleTitle, &r.ArticleDesc, &r.CoverImageUrl, &r.Content, &r.TagID, &r.TagName); err != nil {
			return nil, err
		}

		articles = append(articles, r)
	}

	return articles, nil
}

// 统计tagId 有多少文章
func (a Article) CountByTagID(db *gorm.DB, tagID uint32) (int, error) {
	var count int64

	err := db.Session(&gorm.Session{}).Table(ArticleTag{}.TableName()+" AS at").
		Joins("LEFT JOIN `"+Tag{}.TableName()+"` AS t ON at.tag_id = t.id").
		Joins("LEFT JOIN `"+Article{}.TableName()+"` AS ar ON at.article_id = ar.id").
		Where("at.`tag_id` = ? AND ar.state = ? AND ar.is_del = ?", tagID, a.State, 0).
		Count(&count).Error

	return int(count), err
}
