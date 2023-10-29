package models

import (
	"time"

	"gorm.io/gorm"
)

type Otp struct {
	Id          uint      `gorm:"primaryKey"`
	Otp         int       `gorm:"required"`
	AccessToken string    `gorm:"required"`
	UserId      uint      `gorm:"required"`
	User        User      `gorm:"constraint:OnDelete:CASCADE"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`
}

func CreateOtp(user_id uint, otp int, access_token string, db *gorm.DB) {
	db.Create(&Otp{
		Otp:         otp,
		AccessToken: access_token,
		UserId:      user_id,
	})

}
