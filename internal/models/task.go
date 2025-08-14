package models

import "time"

type Task struct {
	TaskID      uint       `gorm:"column:task_id;primaryKey;autoIncrement" json:"task_id"`
	UserID      uint       `gorm:"column:user_id" json:"user_id"`
	CategoryID  uint       `gorm:"column:category_id" json:"category_id"`
	Title       string     `gorm:"column:title" json:"title"`
	Description string     `gorm:"column:description" json:"description"`
	Status      string     `gorm:"column:status" json:"status"`
	Priority    string     `gorm:"column:priority" json:"priority"`
	DueDate     *time.Time `gorm:"column:due_date" json:"due_date"`
	CompletedAt *time.Time `gorm:"column:completed_at" json:"completed_at"`
	CreatedAt   time.Time  `gorm:"column:created_at" json:"created_at"`
	UpdatedAt   time.Time  `gorm:"column:updated_at" json:"updated_at"`
	User        *User      `gorm:"foreignKey:UserID;references:UserID" json:"user,omitempty"`
	Category    *Category  `gorm:"foreignKey:CategoryID;references:CategoryID" json:"category,omitempty"`
}

func (Task) TableName() string {
	return "tasks"
}
