package dao

import (
	"Golang_Programming_Journey/2_blog-serie/internal/model"
	"Golang_Programming_Journey/2_blog-serie/pkg/app"
	"gorm.io/gorm"
)

func (d *Dao) GetTag(id uint32, state uint8, tx *gorm.DB) (model.Tag, error) {
	tag := model.Tag{Model: &model.Model{ID: id}, State: state}
	return tag.Get(tx)
}

func (d *Dao) GetTagListByIDs(ids []uint32, state uint8) ([]*model.Tag, error) {
	tag := model.Tag{State: state}
	return tag.ListByIDs(d.Engine, ids)
}

func (d *Dao) CountTag(name string, state uint8) (int, error) {
	tag := model.Tag{
		Name:  name,
		State: state,
	}
	return tag.Count(d.Engine)
}

func (d *Dao) GetTagList(name string, state uint8, page, pageSize int) ([]*model.Tag, error) {
	tag := &model.Tag{
		Name:  name,
		State: state,
	}

	pageOffset := app.GetPageOffset(page, pageSize)

	return tag.List(d.Engine, pageOffset, pageSize)

}

func (d *Dao) CreateTag(name string, state uint8, createBy string) error {
	tag := model.Tag{
		Name:  name,
		State: state,
		Model: &model.Model{
			CreatedBy: createBy,
		},
	}
	return tag.Create(d.Engine)

}

func (d *Dao) UpdateTag(id uint32, name string, state uint8, modifiedBy string) error {
	tag := model.Tag{
		Model: &model.Model{
			ID: id,
		},
	}

	values := map[string]interface{}{
		"state":       state,
		"modified_by": modifiedBy,
	}

	if name != "" {
		values["name"] = name
	}

	//这样要注意，如果使用tag去更新的话，gorm看到结构体里面是 0 值就不去更新了， 所以建议使用这种传递map的形式更新
	return tag.Update(d.Engine, values)
}

func (d *Dao) DeleteTag(id uint32) error {
	tag := model.Tag{
		Model: &model.Model{
			ID: id,
		},
	}
	return tag.Delete(d.Engine)
}

func (d *Dao) CheckTagExists(tagName string, state uint8) (bool, error) {
	tag := model.Tag{
		Name:  tagName,
		State: state,
	}
	result, err := tag.Count(d.Engine)
	return result > 0, err
}
