package dao

import "Golang_Programming_Journey/2_blog-serie/internal/model"

func (d *Dao) GetAuth(appKey, appSecret string) (model.Auth, error) {
	auth := model.Auth{
		AppKey:    appKey,
		AppSecret: appSecret,
	}
	return auth.Get(d.Engine)
}
