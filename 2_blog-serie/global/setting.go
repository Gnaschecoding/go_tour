package global

import (
	"Golang_Programming_Journey/2_blog-serie/pkg/logger"
	"Golang_Programming_Journey/2_blog-serie/pkg/setting"
)

var (
	ServerSetting   *setting.ServerSettings
	AppSetting      *setting.AppSettings
	DatabaseSetting *setting.DatabaseSettings
	JWTSetting      *setting.JWTSettings
	Logger          *logger.Logger
	EmailSetting    *setting.EmailSettings
	RedisSetting    *setting.RedisSettings
	LimiterSetting  *setting.LimiterSettings
)
