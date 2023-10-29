package models

import "time"

type Group struct {
	Id          uint `gorm:"primaryKey"`
	Name        string
	CreatedById int
	CreatedBy   User      `gorm:"constraint:OnDelete:CASCADE"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`
}
