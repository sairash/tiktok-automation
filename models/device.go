package models

import "time"

type Device struct {
	Id          uint `gorm:"primaryKey"`
	DeviceInfo  string
	DeviceToken string
	Cookie      string
	UserAgent   string
	DId         string
	IId         string
	Blocked     bool      `gorm:"default:false"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`
}
