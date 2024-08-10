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
		r.POST("user/add", v1.AddUser)

	}
	_ = router.Run(utils.HttpPort)
}
