package controllers

import (
	"net/http"
	"time"

	"github.com/FatahRozaq/taskflow_golang_api/config"
	"github.com/gin-gonic/gin"
)

// PriorityStats represents task count per priority
type PriorityStats struct {
	Priority string `json:"priority"`
	Count    int64  `json:"count"`
}

// CategoryStats represents task count per category
type CategoryStats struct {
	CategoryName string `json:"category_name"`
	Count        int64  `json:"count"`
}

// CompletionStats represents completion rate statistics
type CompletionStats struct {
	Total     int64   `json:"total"`
	Completed int64   `json:"completed"`
	Rate      float64 `json:"completion_rate"`
}

// DashboardStats represents all dashboard statistics
type DashboardStats struct {
	TotalTasks      int64           `json:"total_tasks"`
	PriorityStats   []PriorityStats `json:"tasks_by_priority"`
	CategoryStats   []CategoryStats `json:"tasks_by_category"`
	CompletionStats CompletionStats `json:"completion_stats"`
	TasksDueToday   int64           `json:"tasks_due_today"`
}

func GetStats(c *gin.Context) {
	// Get user ID from Firebase UID parameter
	uid := c.Param("uid")
	if uid == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":     "error",
			"statusCode": http.StatusBadRequest,
			"message":    "UID wajib diisi",
		})
		return
	}

	// Get user by Firebase UID to get internal user ID
	var userID uint
	if err := config.DB.Table("users").Select("user_id").Where("firebase_uid = ?", uid).Scan(&userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":     "error",
			"statusCode": http.StatusNotFound,
			"message":    "User tidak ditemukan",
			"error":      err.Error(),
		})
		return
	}

	var stats DashboardStats

	if err := config.DB.Table("tasks").Where("user_id = ?", userID).Count(&stats.TotalTasks).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":     "error",
			"statusCode": http.StatusInternalServerError,
			"message":    "Gagal mengambil total tasks",
			"error":      err.Error(),
		})
		return
	}

	var priorityStats []PriorityStats
	if err := config.DB.Table("tasks").
		Select("priority, COUNT(*) as count").
		Where("user_id = ?", userID).
		Group("priority").
		Scan(&priorityStats).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":     "error",
			"statusCode": http.StatusInternalServerError,
			"message":    "Gagal mengambil statistik priority",
			"error":      err.Error(),
		})
		return
	}
	stats.PriorityStats = priorityStats

	var categoryStats []CategoryStats
	if err := config.DB.Table("tasks t").
		Select("c.name as category_name, COUNT(t.task_id) as count").
		Joins("LEFT JOIN categories c ON t.category_id = c.category_id").
		Where("t.user_id = ?", userID).
		Group("c.category_id, c.name").
		Scan(&categoryStats).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":     "error",
			"statusCode": http.StatusInternalServerError,
			"message":    "Gagal mengambil statistik category",
			"error":      err.Error(),
		})
		return
	}
	stats.CategoryStats = categoryStats

	var completedTasks int64
	if err := config.DB.Table("tasks").
		Where("user_id = ? AND status = ?", userID, "Done").
		Count(&completedTasks).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":     "error",
			"statusCode": http.StatusInternalServerError,
			"message":    "Gagal mengambil tasks yang completed",
			"error":      err.Error(),
		})
		return
	}

	stats.CompletionStats.Total = stats.TotalTasks
	stats.CompletionStats.Completed = completedTasks
	if stats.TotalTasks > 0 {
		stats.CompletionStats.Rate = float64(completedTasks) / float64(stats.TotalTasks) * 100
	} else {
		stats.CompletionStats.Rate = 0
	}

	today := time.Now().Format("2006-01-02")
	if err := config.DB.Table("tasks").
		Where("user_id = ? AND DATE(due_date) = ?", userID, today).
		Count(&stats.TasksDueToday).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":     "error",
			"statusCode": http.StatusInternalServerError,
			"message":    "Gagal mengambil tasks due today",
			"error":      err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":     "success",
		"statusCode": http.StatusOK,
		"message":    "Statistik dashboard berhasil diambil",
		"data":       stats,
	})
}
