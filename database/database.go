package database

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"beepbop/helper"
	"beepbop/models"
	"beepbop/seed"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func ConnectDb() {

	dsn := "root:root@tcp(127.0.0.1:3306)/tiktok?charset=latin1&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println(err)
	}

	post_folder_path := "./assets/posts"

	if !helper.FolderExists(post_folder_path) {
		err = helper.MakeDir("./assets/posts/")

		if err != nil {
			fmt.Println("Error -> ", err)
		}
	}

	if err = db.AutoMigrate(&models.User{}, &models.Account{}, &models.Post{}, &models.Proxy{}, &models.Notification{}, &models.Name{}, &models.UserStatusType{}, &models.Device{}, &models.Role{}, &models.Country{}, &models.Otp{}, &models.Group{}, &models.GroupDetail{}, &models.CharacterSheet{}, &models.UserCharacter{}, &models.AccountPost{}); err == nil {

		if (db.Migrator().HasTable(&models.Role{})) {

			if err = db.First(&models.Role{}).Error; err != nil {
				role := []models.Role{
					{Name: "admin"},
					{Name: "user"},
				}
				_ = models.CreateRole(&role, db)
			}
		}

		if (db.Migrator().HasTable(&models.UserStatusType{})) {

			if err = db.First(&models.UserStatusType{}).Error; err != nil {
				UserStatusType := []models.UserStatusType{
					{Name: "active"},
					{Name: "blocked"},
				}
				_ = models.CreateUserStatusType(&UserStatusType, db)
			}
		}

		if (db.Migrator().HasTable(&models.Country{})) {
			if err = db.First(&models.Country{}).Error; err != nil {
				resp, err := http.Get("https://raw.githubusercontent.com/boythatcodes/countries_data/main/countries.json")
				if err != nil {
					log.Fatal(err)
				}
				var countries = []models.Country{}
				err = json.NewDecoder(resp.Body).Decode(&countries)
				if err != nil {
					log.Fatal(err)
				}
				models.CreateCountry(&countries, db)
			}
		}

		if (db.Migrator().HasTable(&models.User{})) {
			if err = db.First(&models.User{}).Error; err != nil {
				hashedPassword, err := bcrypt.GenerateFromPassword([]byte("beepbop"), bcrypt.DefaultCost)
				if err != nil {
					log.Fatal(err)
				}

				// head_accessory, body_hand_left, body_hand_right, head_eye_left, head_eye_right, head_mouth, body_accessory, profile_image_command := helper.RandomCharacterGen()

				user := []models.User{
					{
						Username: helper.EnvVariable("ADMIN_USERNAME"),
						Token:    helper.RandomString(50),
						// Character:        profile_image_command,
						CountryId:        160,
						RoleID:           1,
						Verified:         1,
						UserStatusTypeId: 1,
						Password:         string(hashedPassword),
					},
				}

				models.CreateUsers(&user, db)

				// for _, user_loop := range user {
				// 	if _, err = helper.CreateAndSendOtp(user_loop.Id, user_loop.Phone, db); err != nil {
				// 		log.Fatal(err)
				// 	}
				// 	models.UserCharacterSaveUser(user_loop.Id, head_accessory, body_hand_left, body_hand_right, head_eye_left, head_eye_right, head_mouth, body_accessory, db)
				// }

			}
		}

		if (db.Migrator().HasTable(&models.CharacterSheet{})) {
			if err = db.First(&models.CharacterSheet{}).Error; err != nil {
				models.CreateCharacterSheets(seed.HeadAccessory(), db)
				models.CreateCharacterSheets(seed.BodyHandLeft(), db)
				models.CreateCharacterSheets(seed.BodyHandRight(), db)
				models.CreateCharacterSheets(seed.HeadEyeLeft(), db)
				models.CreateCharacterSheets(seed.HeadEyeRight(), db)
				models.CreateCharacterSheets(seed.HeadMouth(), db)
				models.CreateCharacterSheets(seed.BodyAccessory(), db)
			}
		}
	}

	helper.Database = helper.DbInstance{
		Db: db,
	}
}
