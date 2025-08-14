package workers

import (
	"log"

	"github.com/FatahRozaq/taskflow_golang_api/internal/services"
	"github.com/robfig/cron/v3"
)

func Start() {
	log.Println("Initializing cron scheduler...")
	c := cron.New()

	go services.SyncWeatherData()

	_, err := c.AddFunc("@every 30m", services.SyncWeatherData)
	if err != nil {
		log.Fatalf("Could not add weather worker to cron: %v", err)
	}

	_, err = c.AddFunc("* * * * *", services.CheckAndSendReminders)
	if err != nil {
		log.Fatalf("Could not add notification worker to cron: %v", err)
	}

	go c.Start()
	log.Println("Cron scheduler started.")
}
