package dao

import (
	"gorm.io/gorm"
)

type Dao struct {
	Engine *gorm.DB
}

func NewDao(engine *gorm.DB) *Dao {
	return &Dao{Engine: engine}
}
