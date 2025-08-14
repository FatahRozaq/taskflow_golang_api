package services

import (
	"context"
	"fmt"
	"log"
	"time"

	"firebase.google.com/go/v4/messaging"
	"github.com/FatahRozaq/taskflow_golang_api/config"
	"github.com/FatahRozaq/taskflow_golang_api/internal/models"
)

func CheckAndSendReminders() {
	log.Println("Notification worker: Checking for task deadlines...")

	loc, _ := time.LoadLocation("Asia/Jakarta")
	now := time.Now().In(loc)
	deadlineThreshold := now.Add(5 * time.Minute)

	var tasksDueSoon []models.Task

	result := config.DB.
		Joins("JOIN users ON users.user_id = tasks.user_id").
		Where("tasks.due_date BETWEEN ? AND ?", now, deadlineThreshold).
		Where("tasks.status != ?", "Done").
		Where("users.fcm_token IS NOT NULL AND users.fcm_token != ''").
		Preload("User").
		Find(&tasksDueSoon)

	if result.Error != nil {
		log.Printf("Notification worker: Error querying for tasks: %v", result.Error)
		return
	}

	if len(tasksDueSoon) == 0 {
		log.Println("Notification worker: No tasks are due soon. Work cycle complete.")
		log.Println("Go now():", now)
		log.Println("Go now UTC():", now.UTC())
		return
	}

	log.Printf("Notification worker: Found %d task(s) to send reminders for.", len(tasksDueSoon))
	for _, task := range tasksDueSoon {
		sendPushNotification(task)
	}
}

func sendPushNotification(task models.Task) {
	if task.User == nil {
		log.Printf("Notification worker: Cannot send notification for task %d because user data is missing.", task.TaskID)
		return
	}

	message := &messaging.Message{
		Notification: &messaging.Notification{
			Title: "Task Reminder: " + task.Title,
			Body:  fmt.Sprintf("Your task '%s' is due in about 5 minutes!", task.Title),
		},
		Token: task.User.FCMToken,
		Data: map[string]string{
			"taskId":      fmt.Sprintf("%d", task.TaskID),
			"clickAction": "FLUTTER_NOTIFICATION_CLICK",
		},
		Android: &messaging.AndroidConfig{
			Priority: "high",
		},
		APNS: &messaging.APNSConfig{
			Headers: map[string]string{
				"apns-priority": "10",
			},
		},
	}

	response, err := FcmClient.Send(context.Background(), message)
	if err != nil {
		log.Printf("Notification worker: FAILED to send notification for task %d to user %d. Error: %v", task.TaskID, task.UserID, err)
		return
	}

	log.Printf("Notification worker: SUCCESSFULLY sent message for task %d. Message ID: %s", task.TaskID, response)
}
