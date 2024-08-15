package routes

import (
	"github.com/HauKuen/Annals/api/v1"
	"github.com/HauKuen/Annals/utils"
	"github.com/gin-gonic/gin"
)

func InitRouter() {
	gin.SetMode(utils.AppMode)
	router := gin.Default()

	r := router.Group("/api/v1")
	{
		r.GET("user/:id", v1.GetUserInfo)
		r.GET("users", v1.GetUsers)
		r.POST("user/add", v1.AddUser)
		r.DELETE("user/delete/:id", v1.DeleteUser)
		r.PUT("user/edit/:id", v1.EditUser)

		r.POST("category/add", v1.AddCategory)

	}
	_ = router.Run(utils.HttpPort)
}
