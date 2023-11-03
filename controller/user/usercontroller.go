package user

import (
	"beepbop/automation"
	"beepbop/models"
	"fmt"
	"log"
	"net/http"
	"time"

	"beepbop/helper"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

func Me(c echo.Context) error {
	claims, err := helper.JWT(c)

	if err != nil {
		return helper.ErrorResponse_er(c, err, nil)
	}

	var user models.User

	if err := helper.JWTAuthUser(claims.Token, &user); err != nil {
		return helper.ErrorResponse(c, "Error occoured while fetching data.", nil)
	}

	return helper.SuccessResponse(c, "Token Valid", echo.Map{"user_id": user.Id})
}

type loginUser struct {
	Username string `form:"username" validate:"required,lowercase,min=3"`
	Password string `form:"password" validate:"required,min=6"`
}

func Login(c echo.Context) error {

	request := loginUser{}
	if err := c.Bind(&request); err != nil {
		return helper.ErrorResponse(c, "Error Validation", nil)
	}

	if request.Username == "" || request.Password == "" {
		return helper.ErrorResponse(c, "Please fill in something!", nil)
	}
	user := models.User{}

	helper.Database.Db.Select("token", "password", "id", "verified", "RoleID").Where("username = ?", request.Username).Find(&user)

	if user.Id == 0 {
		return helper.ErrorResponse(c, "Username or Password error!", nil)
	}

	if user.Verified != 1 {
		site_data, _, _ := helper.PageDataCreator(c, "Chito Tiktok", "Signup Successful!", "What to do?", "Wait for the account to be approved by admin!", "/", false, 0)

		return c.Render(http.StatusOK, "404", site_data)
	}
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password))
	if err != nil {
		return helper.ErrorResponse(c, "Email or Password error!", nil)
	}

	claims := &helper.UserJwtClaims{
		Token: user.Token,
		Role:  int(user.RoleID),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 72)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	t, err := token.SignedString([]byte(helper.EnvVariable("SERECT")))
	if err != nil {
		return helper.ErrorResponse(c, "Token Not Correct", nil)
	}

	cookie := new(http.Cookie)
	cookie.Name = "toke"
	cookie.Value = t
	cookie.Path = "/"
	c.SetCookie(cookie)

	if user.RoleID == 1 {
		return c.Redirect(http.StatusSeeOther, "/admin")
	}
	return c.Redirect(http.StatusSeeOther, "/home")
}

type UserChecker struct {
	Username string `json:"username" validate:"required,lowercase,min=3"`
}

func CheckUsername(c echo.Context) error {

	request := UserChecker{}
	if err := c.Bind(&request); err != nil {
		return helper.ErrorResponse(c, "Error Validation", nil)
	}

	err := helper.Validator(request)
	if err != nil {
		return helper.ErrorResponse(c, fmt.Sprintf("%e", err), nil)
	}
	user := models.User{}
	helper.Database.Db.First(&user, "username = ?", request.Username)
	if user.Id == 0 {
		return helper.SuccessResponse(c, "Username Available", nil)
	}
	return helper.ErrorResponse(c, "Username Not Available", nil)
}

type SignupUser struct {
	Username string `form:"username" validate:"required,lowercase,min=3"`
	Password string `form:"password" validate:"required,min=6"`
}

func Signup(c echo.Context) error {

	request := SignupUser{}
	if err := c.Bind(&request); err != nil {
		return helper.ErrorResponse(c, "Error Validation", nil)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
	}

	new_user := models.User{
		Username:         request.Username,
		Token:            helper.RandomString(50),
		RoleID:           2,
		Verified:         0,
		CountryId:        1,
		UserStatusTypeId: 1,
		Password:         string(hashedPassword),
	}

	models.CreateUser(&new_user, helper.Database.Db)

	site_data, _, _ := helper.PageDataCreator(c, "Chito Tiktok", "Signup Successful!", "What to do?", "Wait for the account to be approved by admin!", "/", false, 0)

	return c.Render(http.StatusOK, "404", site_data)
}

func Accept(c echo.Context) error {
	userID := c.Param("id")
	claims, err := helper.JWT(c)

	if err != nil || claims.Role != 1 {
		return c.Redirect(http.StatusSeeOther, "/404")
	}

	helper.Database.Db.Model(&models.User{}).Where("id = ?", userID).Update("verified", 1)

	return c.Redirect(http.StatusSeeOther, "/admin")
}

func Delete(c echo.Context) error {
	userID := c.Param("id")

	claims, err := helper.JWT(c)

	if err != nil || claims.Role != 1 {
		return c.Redirect(http.StatusSeeOther, "/404")
	}
	helper.Database.Db.Model(&models.User{}).Delete("id = ?", userID)
	return c.Redirect(http.StatusSeeOther, "/admin")
}

func GetAnyVerifiedUsers(verified int) []models.User {
	users := []models.User{}
	helper.Database.Db.Select("id", "created_at", "username").Where("verified = ? AND role_id != 1", verified).Find(&users)

	return users
}

type automationStruct struct {
	TypeOfAutomation []string `form:"type_of_automation[]"`
	Amount           []string `form:"amount[]"`
}

func CreateAutomation(c echo.Context) error {
	user_claims, err := helper.JWT(c)

	if err != nil || user_claims.Role != 2 {
		c.Redirect(http.StatusSeeOther, "/404")
	}

	user := models.User{}
	helper.Database.Db.First(&user, "token = ?", user_claims.Token)

	request := automationStruct{}
	if err := c.Bind(&request); err != nil {
		return helper.ErrorResponse(c, "Error Validation", nil)
	}

	if err = automation.GenerateAutomationLog(request.TypeOfAutomation, request.Amount, user.Id); err != nil {

		return c.Redirect(http.StatusSeeOther, "/404")
	}

	return c.Redirect(http.StatusSeeOther, "/home/automations")
}

func StopAutomation(c echo.Context) error {
	user_claims, err := helper.JWT(c)

	if err != nil || user_claims.Role != 2 {
		c.Redirect(http.StatusSeeOther, "/404")
	}

	user := models.User{}

	hash_id := c.Param("id")
	helper.Database.Db.First(&user, "token = ?", user_claims.Token)
	if err = automation.RemoveAutomation(user.Id, hash_id); err != nil {
		return c.Redirect(http.StatusSeeOther, "/404")
	}

	return c.Redirect(http.StatusSeeOther, "/home/automations")
}
