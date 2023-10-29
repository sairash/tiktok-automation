package models

import (
	"time"

	"gorm.io/gorm"
)

type UserCharacter struct {
	Id               uint `gorm:"primaryKey"`
	UserId           uint `json:"user_id"`
	User             User `gorm:"constraint:OnDelete:CASCADE"`
	CharacterSheetId uint
	CharacterSheet   CharacterSheet
	CreatedAt        time.Time `gorm:"autoCreateTime"`
	UpdatedAt        time.Time `gorm:"autoUpdateTime"`
}

func UserCharacterSaveUser(user_id uint, head_accessory, body_hand_left, body_hand_right, head_eye_left, head_eye_right, head_mouth, body_accessory uint, db *gorm.DB) {
	_ = db.Create(&[]UserCharacter{
		{
			UserId:           user_id,
			CharacterSheetId: head_accessory,
		},
		{
			UserId:           user_id,
			CharacterSheetId: body_hand_left,
		},
		{
			UserId:           user_id,
			CharacterSheetId: body_hand_right,
		},
		{
			UserId:           user_id,
			CharacterSheetId: head_eye_left,
		},
		{
			UserId:           user_id,
			CharacterSheetId: head_eye_right,
		},
		{
			UserId:           user_id,
			CharacterSheetId: head_mouth,
		},
		{
			UserId:           user_id,
			CharacterSheetId: body_accessory,
		},
	})

}
