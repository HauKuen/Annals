package routes

import (
	"net/http"

	v1 "github.com/HauKuen/Annals/internal/api/v1"
	"github.com/HauKuen/Annals/internal/middleware"
	"github.com/HauKuen/Annals/internal/utils"
	"github.com/gin-gonic/gin"
)

func InitRouter() {
	gin.SetMode(utils.AppMode)
	router := gin.New()

	// 添加 CORS 中间件
	router.Use(cors())
	router.Use(gin.Recovery())
	router.Use(utils.LoggerMiddleware())

	r := router.Group("/api/v1")
	{
		// 公开接口
		auth := r.Group("/auth")
		{
			auth.POST("login", v1.Login)
			auth.GET("validate", v1.ValidateToken)
		}

		// 需要认证的接口
		auth = r.Group("/")
		auth.Use(middleware.JWTAuth())
		{
			// 用户相关接口
			auth.GET("user/:id", v1.GetUserInfo)
			auth.GET("users", v1.GetUsers)
			auth.POST("user/add", v1.AddUser)
			auth.DELETE("user/delete/:id", v1.DeleteUser)
			auth.PUT("user/edit/:id", v1.EditUser)

			// 分类相关接口
			auth.POST("category/add", v1.AddCategory)
			auth.GET("category/:id", v1.GetCategory)
			auth.DELETE("category/delete/:id", v1.DeleteCategory)
			auth.GET("categories", v1.GetCategories)

			// 文章相关接口
			auth.GET("articles", v1.GetArticles)
			auth.GET("article/:id", v1.GetArticle)
			auth.POST("article/add", v1.AddArticle)
			auth.PUT("article/edit/:id", v1.EditArticle)
			auth.DELETE("article/delete/:id", v1.DeleteArticle)
			auth.GET("category/:id/articles", v1.GetCategoryArticles)
			auth.GET("user/:id/articles", v1.GetUserArticles)
			auth.GET("articles/search", v1.SearchArticles)
		}
	}

	if err := router.Run(utils.HttpPort); err != nil {
		utils.Log.Fatal("服务器启动失败:", err)
	}
}

// CORS 中间件
func cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}
