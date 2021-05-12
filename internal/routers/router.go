package routers

import (
	"github.com/gin-gonic/gin"
	_ "github.com/go-programming-tour-book/blog-service/docs"
	"github.com/go-programming-tour-book/blog-service/global"
	"github.com/go-programming-tour-book/blog-service/internal/middleware"
	"github.com/go-programming-tour-book/blog-service/internal/routers/api"
	v1 "github.com/go-programming-tour-book/blog-service/internal/routers/api/v1"
	"github.com/go-programming-tour-book/blog-service/pkg/limiter"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	"net/http"
	"time"
)

//添加一个鉴权存放令牌桶的map（多个令牌桶）
var methodLimiters = limiter.NewMethodLimiter().AddBuckets(limiter.LimiteBucketRule{
	Key:          "/auth",
	FillInterval: time.Second,
	Capacity:     10,
	Quantum:      10,
})

func NewRouter() *gin.Engine {
	r := gin.New()
	if global.ServerSetting.RunMode == "dubug" {
		r.Use(gin.Logger())
		r.Use(gin.Recovery())
	} else {
		r.Use(middleware.AccessLog()) //自定义日志记录
		r.Use(middleware.Recovery()) //自定义捕获异常（告警）
	}

	r.Use(middleware.RateLimiter(methodLimiters)) //限流（通过限制访问令牌桶中令牌数量实现）
	r.Use(middleware.ContextTimeout(global.AppSetting.DefultContextTimeout)) //超时控制
	r.Use(middleware.Translations())
	r.Use(middleware.Tracing()) //注册跟踪链路中间件
	//url := ginSwagger.URL("http://127.0.0.1:8080/swagger/doc.json")
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	article := v1.NewArticle()
	tag := v1.NewTag()

	//文件上传
	upload := api.NewUpload()
	r.POST("/upload/file", upload.UploadFile) //使用curl上传文件路径时请使用绝对路径
	r.StaticFS("/static", http.Dir(global.AppSetting.UploadSavePath))

	//鉴权（获取token）
	r.POST("/auth", api.GetAuth)

	apiv1 := r.Group("/api/v1")

	//只对apiv1路由分组里的路由方法启用鉴权
	//鉴权（验证token）
	apiv1.Use(middleware.JWT())
	{
		apiv1.POST("/tags", tag.Create)
		apiv1.DELETE("/tags/:id", tag.Delete)
		apiv1.PUT("/tags/:id", tag.Update)
		apiv1.PATCH("tags/:id/state", tag.Update)
		apiv1.GET("/tags", tag.List)

		apiv1.POST("/articles", article.Create)
		apiv1.DELETE("/articles/:id", article.Delete)
		apiv1.PUT("/articles/:id", article.Update)
		apiv1.PATCH("/articles/:id/state", article.Update)
		apiv1.GET("/articles/:id", article.Get)
		apiv1.GET("/articles", article.List)
	}

	return r
}
