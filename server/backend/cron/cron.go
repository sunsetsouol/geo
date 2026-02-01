package cron

import (
	"log"

	"geo/server/backend/service"

	"github.com/robfig/cron/v3"
)

func InitCron() {
	c := cron.New()

	// 服务启动时立即运行一次任务生成
	log.Println("Running initial task generation on startup...")
	go service.GenerateDailyTasks()

	// 每天 0 点执行一次
	_, err := c.AddFunc("0 0 * * *", service.GenerateDailyTasks)
	if err != nil {
		log.Fatalf("failed to setup cron: %v", err)
	}
	c.Start()
	log.Println("Cron service started successfully")
}
