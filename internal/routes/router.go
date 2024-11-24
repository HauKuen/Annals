package routes

import (
	v1 "github.com/HauKuen/Annals/internal/api/v1"
	"github.com/HauKuen/Annals/internal/middleware"
	"github.com/HauKuen/Annals/internal/utils"
	"github.com/gin-gonic/gin"
)

func InitRouter() {
	gin.SetMode(utils.AppMode)
	router := gin.New()

	router.Use(gin.Recovery())
	router.Use(utils.LoggerMiddleware())

	r := router.Group("/api/v1")
	{
		// 公开接口
		r.POST("auth/login", v1.Login)

		// 需要认证的接口
		auth := r.Group("/")
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
		}
	}

	if err := router.Run(utils.HttpPort); err != nil {
		utils.Log.Fatal("服务器启动失败:", err)
	}
}
