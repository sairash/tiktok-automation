package models

import (
	"time"

	"gorm.io/gorm"
)

type UserStatusType struct {
	Id        uint `json:"id" gorm:"primaryKey"`
	Name      string
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

func CreateUserStatusType(user_status_types *[]UserStatusType, db *gorm.DB) *gorm.DB {
	result := db.Create(&user_status_types)
	return result
}
