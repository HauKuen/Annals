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
		r.GET("/user/:id", v1.GetUserInfo)
	}
	_ = router.Run(utils.HttpPort)
}
