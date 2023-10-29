package models

import (
	"time"

	"gorm.io/gorm"
)

type CharacterSheet struct {
	Id          uint `gorm:"primaryKey"`
	DisplayName string
	Command     string
	SvgPath     string
	Type        string
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`
}

func CreateCharacterSheet(charactersheet *CharacterSheet, db *gorm.DB) {
	_ = db.Create(&charactersheet)
}
func CreateCharacterSheets(charactersheet []CharacterSheet, db *gorm.DB) {
	_ = db.Create(&charactersheet)
}
