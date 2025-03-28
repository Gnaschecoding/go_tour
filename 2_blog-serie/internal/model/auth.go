package model

import "gorm.io/gorm"

type Auth struct {
	*Model
	AppKey    string `json:"app_key"`
	AppSecret string `json:"app_secret"`
}

// GORM 在数据库操作时就会使用 "blog_auth" 作为表名，而不是 "auths"
func (a Auth) TableName() string {
	return "blog_auth"
}

func (a Auth) Get(db *gorm.DB) (Auth, error) {
	var auth Auth

	db = db.Where("app_key = ? AND app_secret = ? AND is_del = ?",
		a.AppKey,
		a.AppSecret,
		0,
	)
	err := db.First(&auth).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return auth, err
	}
	return auth, nil
}
