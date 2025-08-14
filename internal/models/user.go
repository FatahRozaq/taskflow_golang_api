package models

import "time"

type User struct {
	UserID          uint       `gorm:"column:user_id;primaryKey;autoIncrement" json:"user_id"`
	FirebaseUID     string     `gorm:"column:firebase_uid;unique" json:"firebase_uid"`
	Name            string     `gorm:"column:name" json:"name"`
	Email           string     `gorm:"column:email" json:"email"`
	EmailVerifiedAt *time.Time `gorm:"column:email_verified_at" json:"email_verified_at"`
	Password        string     `gorm:"column:password" json:"-"`
	RememberToken   string     `gorm:"column:remember_token" json:"-"`
	CreatedAt       time.Time  `gorm:"column:created_at" json:"created_at"`
	UpdatedAt       time.Time  `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt       *time.Time `gorm:"column:deleted_at" json:"deleted_at"`
	Tasks           []Task     `gorm:"foreignKey:UserID;references:UserID" json:"tasks"`
	FCMToken        string     `gorm:"column:fcm_token" json:"-"`
}

func (User) TableName() string {
	return "users"
}
