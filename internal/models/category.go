package models

import "time"

type Category struct {
	CategoryID  uint      `gorm:"column:category_id;primaryKey;autoIncrement" json:"category_id"`
	Name        string    `gorm:"column:name" json:"name"`
	Slug        string    `gorm:"column:slug" json:"slug"`
	Color       string    `gorm:"column:color" json:"color"`
	Description string    `gorm:"column:description" json:"description"`
	IsDefault   bool      `gorm:"column:is_default" json:"is_default"`
	CreatedAt   time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at" json:"updated_at"`
	Tasks       []Task    `gorm:"foreignKey:CategoryID;references:CategoryID" json:"tasks"`
}

func (Category) TableName() string {
	return "categories"
}
