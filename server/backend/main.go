package main

import (
	"log"

	"geo/server/backend/config"
	"geo/server/backend/cron"
	"geo/server/backend/dao"
	"geo/server/backend/handler"

	"github.com/gin-gonic/gin"
)

func main() {
	// 初始化配置
	log.Println("Initializing configuration...")
	config.InitConfig()
	log.Println("Configuration initialized successfully")

	// 初始化数据库
	dao.InitDB()
	log.Println("Database connection established successfully")

	// 启动定时任务
	cron.InitCron()
	log.Println("Cron jobs initialized successfully")

	r := gin.Default()

	// 基础路由
	r.GET("/health", handler.HealthCheck)

	// API 路由组
	api := r.Group("/api")
	{
		tasks := api.Group("/tasks")
		{
			tasks.GET("/pending", handler.GetPendingTasks)
			tasks.POST("/:id/result", handler.UpdateTaskResult)
		}

		prompts := api.Group("/prompts")
		{
			prompts.GET("", handler.GetAllPrompts)
			prompts.GET("/:id", handler.GetPromptByID)
			prompts.POST("", handler.CreatePrompt)
			prompts.PUT("/:id", handler.UpdatePrompt)
			prompts.DELETE("/:id", handler.DeletePrompt)
		}
	}

	// 启动服务
	port := config.AppConfig.Server.Port
	log.Printf("Server is starting on port %s...", port)
	if err := r.Run(port); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
