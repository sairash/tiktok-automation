package models

import (
	"time"

	"gorm.io/gorm"
)

type AccountPost struct {
	Id         uint `gorm:"primaryKey"`
	UserId     uint
	User       User `gorm:"constraint:OnDelete:CASCADE"`
	PostID     uint
	Post       Post `gorm:"constraint:OnDelete:CASCADE"`
	AccountId  uint
	Account    Account `gorm:"constraint:OnDelete:CASCADE"`
	TikId      string
	TotalViews int       `gorm:"default:0"`
	TotalLikes int       `gorm:"default:0"`
	CreatedAt  time.Time `gorm:"autoCreateTime"`
	UpdatedAt  time.Time `gorm:"autoUpdateTime"`
}

func CreateAccountPost(accountPost *AccountPost, db *gorm.DB) *gorm.DB {
	return db.Create(accountPost)
}
