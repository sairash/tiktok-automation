package models

import (
	"time"

	"gorm.io/gorm"
)

type Notification struct {
	Id        uint      `json:"id" gorm:"primaryKey"`
	Message   string    `json:"message"`
	Seen      bool      `json:"seen"`
	UserId    uint      `gorm:"required"`
	User      User      `gorm:"constraint:OnDelete:CASCADE"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

func CreateNotification(notifications *[]Notification, db *gorm.DB) *gorm.DB {
	return db.Create(notifications)
}
