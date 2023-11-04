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
	request := AccountControllerRequest{}
	if err := c.Bind(&request); err != nil {
		return helper.ErrorResponse(c, "Error Validation", nil)
	}

	extra_value, user_id, err := helper.PageDataCreator(c, "Create Post", "", "", "", "", true, 2)
	if err != nil {
		c.Redirect(http.StatusSeeOther, "/logout")
	}

	posts := []models.Post{}

	helper.Database.Db.Where("user_id = ?", user_id).Find(&posts)

	extra_value["body"] = helper.UserSidebar("")
	extra_value["user_ids"] = request.UserIDs

	extra_value["posts"] = posts

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

	user := models.User{}
	helper.Database.Db.Select("id").Where("token = ?", user_claims.Token).First(&user)
	account := models.Account{}
	helper.Database.Db.Select("id", "session", "total_views", "total_likes").Where("id = ?", id).First(&account)
	total_likes := account.TotalLikes
	total_views := account.TotalViews

	post := models.AccountPost{}
	helper.Database.Db.Select("id", "tik_id", "total_views", "total_likes").Where("account_id = ? AND tik_id = ?", account.Id, post_id).Find(&post)

	total_views, total_likes = automation.DeletePost(account.Session, post.TikId, 0, total_likes, total_views, post.TotalViews, post.TotalLikes, user.Id)

	if account.TotalLikes != total_likes || account.TotalViews != total_views {
		helper.Database.Db.Where("id = ?", account.Id).Updates(&models.Account{
			TotalViews: total_views,
			TotalLikes: total_likes,
		})
	}

	return c.Redirect(http.StatusSeeOther, "/home/tiktok/posts/"+tiktok_id+"/"+id)
}

func NameAdd(c echo.Context) error {
	user_claims, err := helper.JWT(c)

	if err != nil || user_claims.Role != 2 {
		c.Redirect(http.StatusSeeOther, "/404")
	}

	user := models.User{}

	helper.Database.Db.First(&user, "token = ?", user_claims.Token)

	if user.Id == 0 {
		return c.Redirect(http.StatusSeeOther, "/logout")
	}

	name := c.FormValue("name")

	if name == "" {
		return c.Redirect(http.StatusSeeOther, "/home/names")
	}

	helper.Database.Db.Create(&models.Name{
		Name:   name,
		UserId: user.Id,
	})

	return c.Redirect(http.StatusSeeOther, "/home/names")
}

func ProxyAdd(c echo.Context) error {
	user_claims, err := helper.JWT(c)

	if err != nil || user_claims.Role != 2 {
		c.Redirect(http.StatusSeeOther, "/404")
	}

	user := models.User{}

	helper.Database.Db.First(&user, "token = ?", user_claims.Token)

	if user.Id == 0 {
		return c.Redirect(http.StatusSeeOther, "/logout")
	}

	proxy_url := c.FormValue("url")

	if proxy_url == "" {
		return c.Redirect(http.StatusSeeOther, "/home/proxies")
	}

	helper.Database.Db.Create(&models.Proxy{
		Url:    proxy_url,
		UserId: user.Id,
	})

	return c.Redirect(http.StatusSeeOther, "/home/proxies")
}

func ProxyRefresh(c echo.Context) error {
	user_claims, err := helper.JWT(c)

	if err != nil || user_claims.Role != 2 {
		c.Redirect(http.StatusSeeOther, "/404")
	}

	user := models.User{}

	helper.Database.Db.First(&user, "token = ?", user_claims.Token)

	if user.Id == 0 {
		return c.Redirect(http.StatusSeeOther, "/logout")
	}

	id := c.Param("id")

	proxy := models.Proxy{}
	helper.Database.Db.Where("id = ?", id).Find(&proxy)

	if proxy.UserId == user.Id {
		helper.CheckProxyWorking(proxy.Url, proxy.Id, user.Id, true)
	}

	return c.Redirect(http.StatusSeeOther, "/home/proxies")
}

func ProxyDelete(c echo.Context) error {
	user_claims, err := helper.JWT(c)

	if err != nil || user_claims.Role != 2 {
		c.Redirect(http.StatusSeeOther, "/404")
	}

	user := models.User{}

	helper.Database.Db.First(&user, "token = ?", user_claims.Token)

	if user.Id == 0 {
		return c.Redirect(http.StatusSeeOther, "/logout")
	}

	id := c.Param("id")

	helper.Database.Db.Where("id = ?", id).Delete(&models.Proxy{})

	return c.Redirect(http.StatusSeeOther, "/home/proxies")
}
