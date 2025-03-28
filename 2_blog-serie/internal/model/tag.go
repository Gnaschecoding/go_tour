package model

import (
	"Golang_Programming_Journey/2_blog-serie/pkg/app"
	"gorm.io/gorm"
)

type Tag struct {
	*Model
	Name  string `json:"name"`
	State uint8  `json:"state"`
}

type TagSwagger struct {
	List  []*Tag
	Pager *app.Pager
}

func (t Tag) TableName() string {
	return "blog_tag"
}

// Session(&gorm.Session{})
func (t Tag) Get(db *gorm.DB) (Tag, error) {
	var tag Tag
	err := db.Where("id = ? AND is_del = ? AND state = ?", t.ID, 0, t.State).Find(&tag).Error

	if err != nil && err != gorm.ErrRecordNotFound {
		return tag, err
	}
	return tag, nil

}

func (t Tag) ListByIDs(db *gorm.DB, ids []uint32) ([]*Tag, error) {
	var tags []*Tag
	db = db.Where("is_del = ? AND state = ?", 0, t.State)
	err := db.Where("id IN (?)", ids).Find(&tags).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return tags, err
	}
	return tags, nil
}

func (t Tag) Count(db *gorm.DB) (int, error) {
	var count int64
	if t.Name != "" {
		db = db.Where("name = ?", t.Name)
	}

	db = db.Where("state = ?", t.State)
	err := db.Model(&t).Where("is_del = ?", 0).Count(&count).Error

	if err != nil {
		return 0, err
	}
	return int(count), nil
}

func (t Tag) List(db *gorm.DB, pageOffset, pageSize int) ([]*Tag, error) {
	var tags []*Tag
	var err error
	//Offset 方法：
	//Offset(pageOffset) 中的 pageOffset 是一个整数类型的变量，表示要跳过的记录数。
	//它用于指定从数据库查询结果的哪一条记录开始返回。
	//例如，如果 pageOffset 的值为 10，那么查询结果会跳过前面的 10 条记录。
	//常用于实现分页时，计算当前页之前的所有记录数并跳过，以便获取当前页的数据。

	//Limit 方法：
	//Limit(pageSize) 中的 pageSize 是一个整数类型的变量，表示返回的记录数量。
	//它用于限制查询结果返回的记录条数。例如，如果 pageSize 的值为 20，那么查询结果最多返回 20 条记录。
	//常用于指定每页显示的记录数，配合 Offset 方法实现分页功能。
	if pageOffset >= 0 && pageSize > 0 {
		db = db.Offset(pageOffset).Limit(pageSize)
	}

	if t.Name != "" {
		db = db.Where("name = ?", t.Name)
	}
	db = db.Where("state = ?", t.State)
	err = db.Where("is_del = ?", 0).Find(&tags).Error
	if err != nil {
		return nil, err
	}
	return tags, nil
}

func (t Tag) Create(db *gorm.DB) error {
	//a := 2
	//b := 0
	//c := a / b
	//fmt.Printf("c:%d", c)
	return db.Create(&t).Error
}

func (t Tag) Update(db *gorm.DB, values interface{}) error {

	//
	return db.Model(&t).Where("id = ? AND is_del = ?", t.ID, 0).Updates(values).Error
}

func (t Tag) Delete(db *gorm.DB) error {
	//return db.Unscoped().Where("id = ? AND is_del = ?", t.ID, 0).Delete(&t).Error //这个是硬删除
	return db.Where("id = ? AND is_del = ?", t.ID, 0).Delete(&t).Error
}

//func (t Tag) CheckTagExists(db *gorm.DB) (int, error) {
//	var count int64
//	err := db.Where("name = ? AND is_del = ?", t.Name, 0).Count(&count).Error
//	return int(count), err
//}
