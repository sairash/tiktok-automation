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

func RemovePost(c echo.Context) error {
	user_claims, err := helper.JWT(c)

	if err != nil || user_claims.Role != 2 {
		c.Redirect(http.StatusSeeOther, "/404")
	}

	id := c.Param("id")
	post_id := c.Param("postid")
	tiktok_id := c.Param("tiktok_id")

	account := models.Account{}
	helper.Database.Db.Select("id", "session", "total_views", "total_likes").Where("id = ?", id).First(&account)
	total_likes := account.TotalLikes
	total_views := account.TotalViews

	post := models.AccountPost{}
	helper.Database.Db.Select("id", "tik_id", "total_views", "total_likes").Where("account_id = ? AND tik_id = ?", account.Id, post_id).Find(&post)

	total_views, total_likes = automation.DeletePost(account.Session, post.TikId, 0, total_likes, total_views, post.TotalViews, post.TotalLikes)

	if account.TotalLikes != total_likes || account.TotalViews != total_views {
		helper.Database.Db.Where("id = ?", account.Id).Updates(&models.Account{
			TotalViews: total_views,
			TotalLikes: total_likes,
		})
	}

	return c.Redirect(http.StatusSeeOther, "/home/tiktok/posts/"+tiktok_id+"/"+id)
}
