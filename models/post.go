package models

import (
	"time"

	"gorm.io/gorm"
)

type Post struct {
	Id        uint      `json:"id" gorm:"primaryKey"`
	Title     string    `json:"title"`
	Desc      string    `json:"desc"`
	Music     string    `json:"music"`
	Type      string    `json:"type"`
	Path      string    `json:"path"`
	UserId    uint      `gorm:"required"`
	User      User      `gorm:"constraint:OnDelete:CASCADE"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

func CreatePost(post *[]Post, db *gorm.DB) *gorm.DB {
	return db.Create(post)
}
