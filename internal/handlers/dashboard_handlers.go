package handlers

import (
	"net/http"

	"github.com/FatahRozaq/taskflow_golang_api/internal/services"
	"github.com/gin-gonic/gin"
)

func GetWeather(c *gin.Context) {
	weatherData, lastSync := services.GetCachedWeatherData()

	if weatherData == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Weather data is not available yet. Please try again in a minute.",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"last_sync": lastSync,
		"data":      weatherData,
	})
}
