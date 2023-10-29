package models

import (
	"time"

	"gorm.io/gorm"
)

type Account struct {
	Id            uint `gorm:"primaryKey"`
	UserId        uint
	User          User   `gorm:"constraint:OnDelete:CASCADE"`
	TikUserId     string `json:"user_id"`
	Email         string `json:"email"`
	Password      string `json:"password"`
	Session       string `json:"session"`
	ScreenName    string `json:"screen_name"`
	Name          string `json:"name"`
	AccountPosts  []AccountPost
	TypeOfAccount string    `gorm:"default:'disposable'"`
	Cleared       bool      `gorm:"default:false"`
	Followers     int       `gorm:"default:0"`
	TotalViews    int       `gorm:"default:0"`
	TotalLikes    int       `gorm:"default:0"`
	IsBanned      bool      `gorm:"default:0"`
	CreatedAt     time.Time `gorm:"autoCreateTime"`
	UpdatedAt     time.Time `gorm:"autoUpdateTime"`
}

func CreateAccount(account *Account, db *gorm.DB) *gorm.DB {
	return db.Create(account)
}
