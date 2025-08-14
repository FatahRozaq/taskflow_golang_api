package controllers

import (
	"net/http"

	"github.com/FatahRozaq/taskflow_golang_api/config"
	"github.com/FatahRozaq/taskflow_golang_api/internal/models"

	"github.com/gin-gonic/gin"
)

func GetAllCategories(c *gin.Context) {
	var categories []models.Category

	result := config.DB.Find(&categories)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":     "error",
			"statusCode": http.StatusInternalServerError,
			"message":    "Failed to retrieve categories from database.",
			"error":      result.Error.Error(),
		})
		return
	}

	if len(categories) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"status":     "success",
			"statusCode": http.StatusOK,
			"message":    "No categories found.",
			"data":       []models.Category{},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":     "success",
		"statusCode": http.StatusOK,
		"message":    "Categories retrieved successfully.",
		"data":       categories,
	})
}
