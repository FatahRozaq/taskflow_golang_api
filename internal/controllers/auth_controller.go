package controllers

import (
	"context"
	"net/http"
	"os"
	"path/filepath"
	"time"

	firebase "firebase.google.com/go/v4"
	"github.com/FatahRozaq/taskflow_golang_api/config"
	"github.com/FatahRozaq/taskflow_golang_api/internal/models"
	"github.com/gin-gonic/gin"
	"google.golang.org/api/option"
)

var firebaseApp *firebase.App

func init() {
	// Load service account credentials
	dir, _ := os.Getwd()
	opt := option.WithCredentialsFile(filepath.Join(dir, "serviceAccountKey.json"))
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		panic("Gagal init Firebase: " + err.Error())
	}
	firebaseApp = app
}

type RegisterRequest struct {
	Token string `json:"token" binding:"required"`
	Name  string `json:"name"`
	Email string `json:"email"`
	UID   string `json:"uid" binding:"required"`
}

func RegisterUser(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":     "error",
			"statusCode": http.StatusBadRequest,
			"message":    "Invalid request body",
			"error":      err.Error(),
		})
		return
	}

	// Verifikasi token Firebase
	authClient, err := firebaseApp.Auth(context.Background())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":     "error",
			"statusCode": http.StatusInternalServerError,
			"message":    "Gagal inisialisasi Firebase Auth",
			"error":      err.Error(),
		})
		return
	}

	token, err := authClient.VerifyIDToken(context.Background(), req.Token)
	if err != nil || token.UID != req.UID {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":     "error",
			"statusCode": http.StatusUnauthorized,
			"message":    "Token tidak valid",
		})
		return
	}

	// Simpan user ke database
	user := models.User{
		FirebaseUID: req.UID,
		Name:        req.Name,
		Email:       req.Email,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := config.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":     "error",
			"statusCode": http.StatusInternalServerError,
			"message":    "Gagal menyimpan user",
			"error":      err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":     "success",
		"statusCode": http.StatusOK,
		"message":    "User berhasil diregister",
		"data":       user,
	})
}

func GetUserByUID(c *gin.Context) {
	uid := c.Param("uid")
	if uid == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":     "error",
			"statusCode": http.StatusBadRequest,
			"message":    "UID wajib diisi",
		})
		return
	}

	var user models.User
	if err := config.DB.Where("firebase_uid = ?", uid).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":     "error",
			"statusCode": http.StatusNotFound,
			"message":    "User tidak ditemukan",
			"error":      err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":     "success",
		"statusCode": http.StatusOK,
		"data": gin.H{
			"userId": user.UserID,
			"uid":    user.FirebaseUID,
			"name":   user.Name,
			"email":  user.Email,
		},
	})
}
