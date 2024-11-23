package main

import (
	"github.com/HauKuen/Annals/internal/model"
	"github.com/HauKuen/Annals/internal/routes"
)

func main() {
	// 初始化数据库连接
	model.InitDb()

	// 初始化路由并启动服务器
	routes.InitRouter()
}
