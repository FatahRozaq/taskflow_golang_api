package controllers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/FatahRozaq/taskflow_golang_api/config"
	"github.com/FatahRozaq/taskflow_golang_api/internal/models"
	"github.com/gin-gonic/gin"
)

type TaskInput struct {
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Status      string     `json:"status"`
	Priority    string     `json:"priority"`
	UserID      uint       `json:"user_id"`
	CategoryID  uint       `json:"category_id"`
	DueDate     *time.Time `json:"due_date"`
	CompletedAt *time.Time `json:"completed_at"`
}

// GetUserTasks by user id
func GetUserTasks(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":     "error",
			"statusCode": http.StatusBadRequest,
			"message":    "Invalid user ID format. User ID must be a valid number.",
		})
		return
	}

	var tasks []models.Task

	result := config.DB.Preload("Category").Where("user_id = ?", userID).Find(&tasks)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":     "error",
			"statusCode": http.StatusInternalServerError,
			"message":    "Failed to retrieve tasks from database.",
			"error":      result.Error.Error(),
		})
		return
	}

	if len(tasks) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"status":     "success",
			"statusCode": http.StatusOK,
			"message":    "No tasks found for this user.",
			"data":       []models.Task{},
		})
		return
	}

	type TaskResponse struct {
		TaskID      uint       `json:"task_id"`
		UserID      uint       `json:"user_id"`
		CategoryID  uint       `json:"category_id"`
		Title       string     `json:"title"`
		Description string     `json:"description"`
		Status      string     `json:"status"`
		Priority    string     `json:"priority"`
		DueDate     *time.Time `json:"due_date"`
		CompletedAt *time.Time `json:"completed_at"`
		CreatedAt   time.Time  `json:"created_at"`
		UpdatedAt   time.Time  `json:"updated_at"`
		Category    struct {
			CategoryID  uint   `json:"category_id"`
			Name        string `json:"name"`
			Color       string `json:"color"`
			Description string `json:"description"`
			IsDefault   bool   `json:"is_default"`
		} `json:"category"`
	}

	var response []TaskResponse
	for _, t := range tasks {
		tr := TaskResponse{
			TaskID:      t.TaskID,
			UserID:      t.UserID,
			CategoryID:  t.CategoryID,
			Title:       t.Title,
			Description: t.Description,
			Status:      t.Status,
			Priority:    t.Priority,
			DueDate:     t.DueDate,
			CompletedAt: t.CompletedAt,
			CreatedAt:   t.CreatedAt,
			UpdatedAt:   t.UpdatedAt,
		}
		if t.Category != nil {
			tr.Category.CategoryID = t.Category.CategoryID
			tr.Category.Name = t.Category.Name
			tr.Category.Color = t.Category.Color
			tr.Category.Description = t.Category.Description
			tr.Category.IsDefault = t.Category.IsDefault
		}
		response = append(response, tr)
	}

	c.JSON(http.StatusOK, gin.H{
		"status":     "success",
		"statusCode": http.StatusOK,
		"message":    "Tasks retrieved successfully.",
		"data":       response,
	})
}

func validateTaskInput(input *TaskInput) map[string]string {
	errors := make(map[string]string)

	if input.Title == "" {
		errors["title"] = "Title is required and cannot be empty."
	} else if len(input.Title) > 255 {
		errors["title"] = "Title must be at most 255 characters long."
	}

	validStatuses := map[string]bool{
		"Todo":        true,
		"In Progress": true,
		"Done":        true,
		"pending":     true,
	}
	if input.Status == "" {
		errors["status"] = "Status is required and cannot be empty."
	} else if !validStatuses[input.Status] {
		errors["status"] = "Status must be one of: 'Todo', 'In Progress', 'Done', or 'pending'."
	}

	validPriorities := map[string]bool{
		"Low":    true,
		"Medium": true,
		"High":   true,
	}
	if input.Priority == "" {
		errors["priority"] = "Priority is required and cannot be empty."
	} else if !validPriorities[input.Priority] {
		errors["priority"] = "Priority must be one of: 'Low', 'Medium', or 'High'."
	}

	if input.UserID == 0 {
		errors["user_id"] = "User ID is required and must be a valid positive number."
	}

	if input.CategoryID == 0 {
		errors["category_id"] = "Category ID is required and must be a valid positive number."
	}

	if input.DueDate != nil && input.DueDate.Before(time.Now().Truncate(24*time.Hour)) {
		errors["due_date"] = "Due date cannot be in the past."
	}

	if input.CompletedAt != nil && input.CompletedAt.After(time.Now()) {
		errors["completed_at"] = "Completion date cannot be in the future."
	}

	if input.Status == "Done" && input.CompletedAt == nil {
		errors["completed_at"] = "Completion date is required when status is 'Done'."
	} else if input.Status != "Done" && input.CompletedAt != nil {
		errors["completed_at"] = "Completion date should only be set when status is 'Done'."
	}

	return errors
}

func CreateTask(c *gin.Context) {
	var input TaskInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":     "error",
			"statusCode": http.StatusBadRequest,
			"message":    "Invalid JSON format. Please check your request body.",
			"error":      err.Error(),
		})
		return
	}

	jakartaLoc, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":     "error",
			"statusCode": http.StatusInternalServerError,
			"message":    "Failed to load timezone",
			"error":      err.Error(),
		})
		return
	}

	if input.DueDate != nil {
		*input.DueDate = input.DueDate.In(jakartaLoc)
	} else {
		today := time.Now().In(jakartaLoc)
		defaultDue := time.Date(today.Year(), today.Month(), today.Day(), 23, 59, 0, 0, jakartaLoc)
		input.DueDate = &defaultDue
	}

	if input.CompletedAt != nil {
		*input.CompletedAt = input.CompletedAt.In(jakartaLoc)
	}

	if validationErrors := validateTaskInput(&input); len(validationErrors) > 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":     "error",
			"statusCode": http.StatusBadRequest,
			"message":    "Validation failed. Please check the errors below.",
			"errors":     validationErrors,
		})
		return
	}

	task := models.Task{
		Title:       input.Title,
		Description: input.Description,
		Status:      input.Status,
		Priority:    input.Priority,
		UserID:      input.UserID,
		CategoryID:  input.CategoryID,
		DueDate:     input.DueDate,
		CompletedAt: input.CompletedAt,
	}

	if err := config.DB.Create(&task).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":     "error",
			"statusCode": http.StatusInternalServerError,
			"message":    "Failed to create task. Please try again later.",
			"error":      err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status":     "success",
		"statusCode": http.StatusCreated,
		"message":    "Task created successfully.",
		"data":       task,
	})
}

func UpdateTask(c *gin.Context) {
	id := c.Param("id")
	var task models.Task

	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":     "error",
			"statusCode": http.StatusBadRequest,
			"message":    "Task ID is required in the URL parameter.",
		})
		return
	}

	if err := config.DB.First(&task, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":     "error",
			"statusCode": http.StatusNotFound,
			"message":    "Task not found. Please check the task ID.",
		})
		return
	}

	var input TaskInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":     "error",
			"statusCode": http.StatusBadRequest,
			"message":    "Invalid JSON format. Please check your request body.",
			"error":      err.Error(),
		})
		return
	}

	jakartaLoc, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":     "error",
			"statusCode": http.StatusInternalServerError,
			"message":    "Failed to load timezone",
			"error":      err.Error(),
		})
		return
	}

	var dueDate *time.Time
	var completedAt *time.Time

	if input.DueDate != nil {
		jakartaTime := input.DueDate.In(jakartaLoc)
		dueDate = &jakartaTime
	}

	if input.CompletedAt != nil {
		jakartaTime := input.CompletedAt.In(jakartaLoc)
		completedAt = &jakartaTime
	}

	if input.Status == "Done" && completedAt == nil {
		now := time.Now().In(jakartaLoc)
		completedAt = &now
	}

	if input.Status != "Done" {
		completedAt = nil
	}

	if validationErrors := validateTaskInput(&input); len(validationErrors) > 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":     "error",
			"statusCode": http.StatusBadRequest,
			"message":    "Validation failed. Please check the errors below.",
			"errors":     validationErrors,
		})
		return
	}

	task.Title = input.Title
	task.Description = input.Description
	task.Status = input.Status
	task.Priority = input.Priority
	task.UserID = input.UserID
	task.CategoryID = input.CategoryID
	task.DueDate = dueDate
	task.CompletedAt = completedAt

	if err := config.DB.Save(&task).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":     "error",
			"statusCode": http.StatusInternalServerError,
			"message":    "Failed to update task. Please try again later.",
			"error":      err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":     "success",
		"statusCode": http.StatusOK,
		"message":    "Task updated successfully.",
		"data":       task,
	})
}

func DeleteTask(c *gin.Context) {
	id := c.Param("id")
	var task models.Task

	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":     "error",
			"statusCode": http.StatusBadRequest,
			"message":    "Task ID is required in the URL parameter.",
		})
		return
	}

	if err := config.DB.First(&task, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":     "error",
			"statusCode": http.StatusNotFound,
			"message":    "Task not found. Please check the task ID.",
		})
		return
	}

	if err := config.DB.Delete(&task).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":     "error",
			"statusCode": http.StatusInternalServerError,
			"message":    "Failed to delete task. Please try again later.",
			"error":      err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":     "success",
		"statusCode": http.StatusOK,
		"message":    "Task deleted successfully.",
	})
}
