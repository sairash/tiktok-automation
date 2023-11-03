package models

import (
	"time"

	"gorm.io/gorm"
)

type Name struct {
	Id        uint      `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name"`
	UserId    uint      `gorm:"required"`
	User      User      `gorm:"constraint:OnDelete:CASCADE"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

func CreateName(names *[]Name, db *gorm.DB) *gorm.DB {
	return db.Create(names)
}
