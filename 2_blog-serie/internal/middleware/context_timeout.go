package middleware

import (
	"context"
	"github.com/gin-gonic/gin"
	"time"
)

func ContextTimeout(t time.Duration) func(c *gin.Context) {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), t)
		defer cancel()

		//将原请求对象替换为一个关联了新上下文的请求对象，
		//这样后续对该请求的处理过程中，就可以通过请求对象获取到这个新的上下文信息。
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}
