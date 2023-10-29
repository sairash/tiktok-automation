package models

import (
	"gorm.io/gorm"
)

type Country struct {
	Id                 uint    `json:"id" gorm:"primaryKey"`
	Name               string  `json:"name" sql:"type: CHARACTER SET utf8 COLLATE utf8_general_ci"`
	TimezoneOffset     float32 `json:"timezone_offset" sql:"type: CHARACTER SET utf8 COLLATE utf8_general_ci"`
	LatLng             string  `json:"latlong" sql:"type: CHARACTER SET utf8 COLLATE utf8_general_ci"`
	Code               string  `json:"code" sql:"type: CHARACTER SET utf8 COLLATE utf8_general_ci"`
	Region             string  `json:"region" sql:"type: CHARACTER SET utf8 COLLATE utf8_general_ci"`
	Capital            string  `json:"capital" sql:"type: CHARACTER SET utf8 COLLATE utf8_general_ci"`
	Isocode            string  `json:"isocode" sql:"type: CHARACTER SET utf8 COLLATE utf8_general_ci"`
	DiallingCode       string  `json:"dialling_code" sql:"type: CHARACTER SET utf8 COLLATE utf8_general_ci"`
	CurrencyCode       string  `json:"currency_code" sql:"type: CHARACTER SET utf8 COLLATE utf8_general_ci"`
	CurrencyName       string  `json:"currency_name" sql:"type: CHARACTER SET utf8 COLLATE utf8_general_ci"`
	CurrencySymbol     string  `json:"currency_symbol" sql:"type: CHARACTER SET utf8 COLLATE utf8_general_ci"`
	LanguageCode       string  `json:"language_code" sql:"type: CHARACTER SET utf8 COLLATE utf8_general_ci"`
	LanguageName       string  `json:"language_name" sql:"type: CHARACTER SET utf8 COLLATE utf8_general_ci"`
	LanguageNativeName string  `json:"language_native_name" sql:"type: CHARACTER SET utf8 COLLATE utf8_general_ci"`
	PhoneLengthMin     int     `json:"phone_length_min"`
	PhoneLengthMax     int     `json:"phone_length_max"`
}

func CreateCountry(c *[]Country, db *gorm.DB) *gorm.DB {
	result := db.Create(c)
	return result
}
