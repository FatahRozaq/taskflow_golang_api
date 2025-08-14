package routes

import (
	"github.com/FatahRozaq/taskflow_golang_api/internal/controllers"
	"github.com/FatahRozaq/taskflow_golang_api/internal/handlers"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {

	api := r.Group("/api")
	{
		dashboardGroup := api.Group("/dashboard")
		{
			dashboardGroup.GET("/weather", handlers.GetWeather)
			dashboardGroup.GET("/stats/:uid", controllers.GetStats)
		}

		authGroup := api.Group("/auth")
		{
			authGroup.POST("/register", controllers.RegisterUser)
			authGroup.GET("/user/:uid", controllers.GetUserByUID)
		}

		api.GET("/users/:id/tasks", controllers.GetUserTasks)
		api.POST("/tasks", controllers.CreateTask)
		api.PUT("/tasks/:id", controllers.UpdateTask)
		api.DELETE("/tasks/:id", controllers.DeleteTask)
		api.GET("/categories", controllers.GetAllCategories)
		api.GET("/dashboard/stats", controllers.GetStats)
	}
}
