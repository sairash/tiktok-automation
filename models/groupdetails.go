package models

import "time"

type GroupDetail struct {
	Id        uint `gorm:"primaryKey"`
	GroupId   uint
	Group     Group
	UserId    uint
	User      User      `gorm:"constraint:OnDelete:CASCADE"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}
