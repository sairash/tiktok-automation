package models

import (
	"time"

	"gorm.io/gorm"
)

type Proxy struct {
	Id         uint      `json:"id" gorm:"primaryKey"`
	Url        string    `json:"url"`
	NotWorking bool      `json:"not_working"`
	UserId     uint      `gorm:"required"`
	User       User      `gorm:"constraint:OnDelete:CASCADE"`
	CreatedAt  time.Time `gorm:"autoCreateTime"`
	UpdatedAt  time.Time `gorm:"autoUpdateTime"`
}

func CreateProxy(proxies *[]Proxy, db *gorm.DB) *gorm.DB {
	return db.Create(proxies)
}
