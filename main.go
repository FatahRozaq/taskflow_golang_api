package main

import (
	"log"
	"os"

	"github.com/FatahRozaq/taskflow_golang_api/config"
	"github.com/FatahRozaq/taskflow_golang_api/internal/middleware"
	"github.com/FatahRozaq/taskflow_golang_api/internal/routes"
	"github.com/FatahRozaq/taskflow_golang_api/internal/services"
	"github.com/FatahRozaq/taskflow_golang_api/internal/workers"
	"github.com/gin-gonic/gin"
)

func main() {
	config.LoadConfig()
	config.ConnectDB()
	services.InitFirebase()
	workers.Start()

	router := gin.Default()

	router.Use(middleware.CorsMiddleware())

	routes.SetupRoutes(router)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("🚀 Server starting on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("❌ Failed to run server: %v", err)
	}
}
