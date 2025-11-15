package main

import (
	"log"

	"learn-hub-backend/config"
	"learn-hub-backend/database"
	"learn-hub-backend/routes"
)

func main() {
	// 加载配置
	config.LoadConfig()

	// 初始化数据库
	database.InitDB()

	// 设置路由
	r := routes.SetupRoutes()

	// 启动服务器
	port := ":" + config.AppConfig.Port
	log.Printf("服务器启动在端口 %s", port)
	if err := r.Run(port); err != nil {
		log.Fatal("服务器启动失败:", err)
	}
}
