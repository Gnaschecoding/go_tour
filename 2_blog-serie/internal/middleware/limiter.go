package middleware

import (
	"Golang_Programming_Journey/2_blog-serie/global"
	"Golang_Programming_Journey/2_blog-serie/pkg/app"
	"Golang_Programming_Journey/2_blog-serie/pkg/errcode"
	"Golang_Programming_Journey/2_blog-serie/pkg/limiter"
	"github.com/gin-gonic/gin"
)

func RateLimiter(l limiter.LimiterIface) gin.HandlerFunc {

	return func(c *gin.Context) {
		if rl, ok := l.(*limiter.RedisLimiter); ok {
			key := rl.Key(c) //根据url拼接ip地址加上前缀
			_, ok := rl.GetBucket(key)
			//如果 redis中没有存这个key，则根据Rule新建一个Buckets
			if !ok {
				rl.AddBuckets(
					limiter.LimiterBucketRule{
						Key:          key,
						FillInterval: global.LimiterSetting.FillInterval,
						Capacity:     global.LimiterSetting.Capacity,
						Quantum:      global.LimiterSetting.Quantum,
					},
				)
			}

			//取一个令牌，如果能去掉不执行if里面的内容
			if !rl.TakeToken(key) {
				app.NewResponse(c).ToErrorResponse(errcode.TooManyRequests)
				c.Abort()
				return
			}
		} else {
			// 兼容旧的限流器
			key := l.Key(c)
			bucket, ok := l.GetBucket(key)
			if !ok {
				l.AddBuckets(limiter.LimiterBucketRule{
					Key:          key,
					FillInterval: global.LimiterSetting.FillInterval,
					Capacity:     global.LimiterSetting.Capacity,
					Quantum:      global.LimiterSetting.Quantum,
				})
				bucket, _ = l.GetBucket(key)
			}

			count := bucket.TakeAvailable(1)
			if count == 0 {
				app.NewResponse(c).ToErrorResponse(errcode.TooManyRequests)
				c.Abort()
				return
			}
		}
		c.Next()
	}
}
