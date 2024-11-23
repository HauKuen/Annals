package routes

import (
	v1 "github.com/HauKuen/Annals/internal/api/v1"
	"github.com/HauKuen/Annals/internal/utils"
	"github.com/gin-gonic/gin"
)

func InitRouter() {
	gin.SetMode(utils.AppMode)
	router := gin.Default()

	router.Use(utils.LoggerMiddleware())

	r := router.Group("/api/v1")
	{
		r.GET("user/:id", v1.GetUserInfo)
		r.GET("users", v1.GetUsers)
		r.POST("user/add", v1.AddUser)
		r.DELETE("user/delete/:id", v1.DeleteUser)
		r.PUT("user/edit/:id", v1.EditUser)

		r.POST("category/add", v1.AddCategory)
		r.GET("category/:id", v1.GetCategory)
		r.DELETE("category/delete/:id", v1.DeleteCategory)
	}
	_ = router.Run(utils.HttpPort)
}
