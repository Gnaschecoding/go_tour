package routers

import (
	_ "Golang_Programming_Journey/2_blog-serie/docs"
	"Golang_Programming_Journey/2_blog-serie/global"
	"Golang_Programming_Journey/2_blog-serie/internal/middleware"
	"Golang_Programming_Journey/2_blog-serie/internal/routers/api"
	v1 "Golang_Programming_Journey/2_blog-serie/internal/routers/api/v1"
	"Golang_Programming_Journey/2_blog-serie/pkg/limiter"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"net/http"
)

// 基于本地计数的 限流器
var methodLimiter limiter.LimiterIface

// 创建 Redis 限流器
var redisLimiter limiter.LimiterIface

func returnLimiter() limiter.LimiterIface {
	redisLimiter = limiter.NewRedisLimiter(global.RedisClient, "rate_limit:").AddBuckets(
		limiter.LimiterBucketRule{
			Key:          "/auth",
			FillInterval: global.LimiterSetting.FillInterval,
			Capacity:     global.LimiterSetting.Capacity,
			Quantum:      global.LimiterSetting.Quantum,
		},
	)

	methodLimiter = limiter.NewMethodLimiter().AddBuckets(
		limiter.LimiterBucketRule{
			Key:          "/auth",
			FillInterval: global.LimiterSetting.FillInterval,
			Capacity:     global.LimiterSetting.Capacity,
			Quantum:      global.LimiterSetting.Quantum,
		},
	)

	return redisLimiter
}

func NewRouter() *gin.Engine {
	r := gin.New()

	if global.ServerSetting.RunMode == "debug" {
		r.Use(gin.Logger())
		r.Use(gin.Recovery())
	} else {
		r.Use(middleware.AccessLog())
		r.Use(middleware.Recovery())
	}

	r.Use(middleware.Tracing())

	r.Use(middleware.RateLimiter(returnLimiter()))
	r.Use(middleware.ContextTimeout(global.AppSetting.DefaultContextTimeout))
	r.Use(middleware.AppInfo())
	r.Use(middleware.Translations())

	//这个是加载 swagger文档的
	r.GET("swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.StaticFS("/static", http.Dir(global.AppSetting.UploadSavePath))

	article := v1.NewArticle()
	tag := v1.NewTag()

	upload := api.NewUpload()
	r.POST("/upload/file", upload.UploadFile)

	r.POST("/auth", api.GetAuth)

	apiv1 := r.Group("/api/v1")
	apiv1.Use(middleware.JWT())
	{
		// 创建标签
		apiv1.POST("/tags", tag.Create)
		// 删除指定标签
		apiv1.DELETE("/tags/:id", tag.Delete)
		// 更新指定标签
		//apiv1.PATCH("/tags/:id/state", tag.Update)
		apiv1.PUT("/tags/:id", tag.Update)
		// 获取标签列表
		apiv1.GET("/tags", tag.List)

		// 创建文章
		apiv1.POST("/articles", article.Create)
		// 删除指定文章
		apiv1.DELETE("/articles/:id", article.Delete)
		// 更新指定文章
		apiv1.PUT("/articles/:id", article.Update)

		//apiv1.PATCH("/articles/:id/state", article.Update)
		// 获取指定文章
		apiv1.GET("/articles/:id", article.Get)
		// 获取文章列表
		apiv1.GET("/articles", article.List)

	}

	return r

}
