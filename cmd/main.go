package main

import (
	"time"

	"github.com/HauKuen/Annals/internal/model"
	"github.com/HauKuen/Annals/internal/routes"
	"github.com/HauKuen/Annals/internal/utils"
)

func main() {
	// 初始化数据库连接
	if err := model.InitDb(); err != nil {
		utils.Log.Fatal("数据库初始化失败:", err)
	}

	// 定期检查数据库健康状况
	go monitorDatabaseHealth()

	// 初始化路由并启动服务器
	routes.InitRouter()
}

func monitorDatabaseHealth() {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		if err := model.CheckDatabaseHealth(); err != nil {
			utils.Log.Error("数据库健康检查失败:", err)
		}
	}
}
