package middleware

import (
	"Golang_Programming_Journey/2_blog-serie/pkg/app"
	"Golang_Programming_Journey/2_blog-serie/pkg/errcode"
	"Golang_Programming_Journey/2_blog-serie/pkg/limiter"
	"github.com/gin-gonic/gin"
)

func RateLimiter(l limiter.LimiterIface) gin.HandlerFunc {
	return func(c *gin.Context) {
		key := l.Key(c)

		if bucket, ok := l.GetBucket(key); ok {
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
