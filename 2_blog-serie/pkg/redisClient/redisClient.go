package redisClient

import (
	"Golang_Programming_Journey/2_blog-serie/pkg/setting"
	"github.com/redis/go-redis/v9"
)

//// 初始化 Redis 客户端  //可以把他修改到 配置文件中
//var redisClient = redis.NewClient(&redis.Options{
//	Addr:     "localhost:6379", // 根据实际情况修改
//	Password: "",               // no password set
//	DB:       0,                // use default DB
//})

func NewRedisClient(redisSetting *setting.RedisSettings) *redis.Client {
	return redis.NewClient(&redis.Options{
		Network:      redisSetting.Network,
		Addr:         redisSetting.Addr,
		Username:     redisSetting.Username,
		Password:     redisSetting.Password,
		DialTimeout:  redisSetting.DialTimeout,
		ReadTimeout:  redisSetting.ReadTimeout,
		WriteTimeout: redisSetting.WriteTimeout,
		PoolSize:     redisSetting.PoolSize,
		MinIdleConns: redisSetting.MinIdleConns,
		MaxIdleConns: redisSetting.MaxIdleConns,
		MaxRetries:   redisSetting.MaxRetries,
		DB:           redisSetting.DB,
	})
}
