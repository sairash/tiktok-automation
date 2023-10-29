package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Id               uint   `json:"id" gorm:"primaryKey"`
	Username         string `json:"username" sql:"type:VARCHAR(5) CHARACTER SET utf8 COLLATE utf8_general_ci"`
	Password         string `json:"password"`
	LocHash          string
	Character        string
	Token            string
	CountryId        uint
	Country          Country
	RoleID           uint `gorm:"default:2"`
	Role             Role
	UserStatusTypeId int `gorm:"default:1"`
	UserStatusType   UserStatusType
	Verified         int `gorm:"default:0"`
	Status_until     time.Time
	Posts            []Post
	Accounts         []Account
	Dob              time.Time `json:"dob"`
	CreatedAt        time.Time `gorm:"autoCreateTime"`
	UpdatedAt        time.Time `gorm:"autoUpdateTime"`
}

func CreateUser(user *User, db *gorm.DB) {
	db.Create(&user)

}

func CreateUsers(user *[]User, db *gorm.DB) {
	db.Create(&user)
}
