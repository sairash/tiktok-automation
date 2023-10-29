package account

import (
	"beepbop/automation"
	"beepbop/helper"
	"beepbop/models"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

type AccountControllerRequest struct {
	UserIDs       string `form:"user_ids" query:"user_ids"`
	TypeOfAccount string `form:"type_of_account" query:"type_of_account"`
}

func RefreshAccount(c echo.Context) error {
	user_claims, err := helper.JWT(c)

	if err != nil || user_claims.Role != 2 {
		c.Redirect(http.StatusSeeOther, "/404")
	}

	request := AccountControllerRequest{}
	if err := c.Bind(&request); err != nil {
		return helper.ErrorResponse(c, "Error Validation", nil)
	}

	user := models.User{}

	helper.Database.Db.First(&user, "token = ?", user_claims.Token)

	automation.GenerateAutomationLog([]string{"use_account", "refresh_account"}, []string{request.UserIDs, "0"}, user.Id)

	return c.Redirect(http.StatusSeeOther, "/home/automations")
}

func ClearAccount(c echo.Context) error {
	user_claims, err := helper.JWT(c)

	if err != nil || user_claims.Role != 2 {
		c.Redirect(http.StatusSeeOther, "/404")
	}

	id := c.Param("id")

	helper.Database.Db.Model(&models.Account{}).Where("id = ?", id).Update("cleared", 1)

	return c.Redirect(http.StatusSeeOther, "/home")
}

func ContainedAccount(c echo.Context) error {
	user_claims, err := helper.JWT(c)

	if err != nil || user_claims.Role != 2 {
		c.Redirect(http.StatusSeeOther, "/404")
	}

	request := AccountControllerRequest{}
	if err := c.Bind(&request); err != nil {
		return helper.ErrorResponse(c, "Error Validation", nil)
	}

	helper.Database.Db.Model(&models.Account{}).Where("id In ?", strings.Split(request.UserIDs, ",")).Update("type_of_account", "contained")

	return c.Redirect(http.StatusSeeOther, "/home")
}

func AutomateDisplay(c echo.Context) error {
	user_claims, err := helper.JWT(c)

	if err != nil || user_claims.Role != 2 {
		c.Redirect(http.StatusSeeOther, "/404")
	}

	request := AccountControllerRequest{}
	if err := c.Bind(&request); err != nil {
		return helper.ErrorResponse(c, "Error Validation", nil)
	}

	user := models.User{}

	helper.Database.Db.Preload("Posts").First(&user, "token = ?", user_claims.Token)

	extra_value := helper.PageDataCreator(c, "Create Post", "", "", "", "", true, 2)
	extra_value["body"] = helper.UserSidebar("")
	extra_value["user_ids"] = request.UserIDs

	extra_value["posts"] = user.Posts

	return c.Render(http.StatusOK, "user/automate", extra_value)
}
