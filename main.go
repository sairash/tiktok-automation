package main

import (
	"beepbop/automation"
	"beepbop/controller/user"
	"beepbop/database"
	"beepbop/helper"
	"beepbop/models"
	"beepbop/router"
	"time"

	"net/http"

	"github.com/foolin/goview/supports/echoview-v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {

	// go automation.InitiateAutomation(helper.RandomString(20), [][]automation.AutomationLog{
	// 	{
	// 		// {TypeOfAutomation: "wait", Amount: helper.RandomNumber(1), UserId: 2},
	// 		// {TypeOfAutomation: "create_account", Amount: 2, UserId: 2},
	// 		// {TypeOfAutomation: "create_user", Amount: 1, UserId: 2},
	// 		// {TypeOfAutomation: "wait", Amount: helper.RandomNumber(1), UserId: 2},
	// 		// {TypeOfAutomation: "post", Amount: 1, UserId: 2, NeededAccountIds: []uint{1, 2}},
	// 		// {TypeOfAutomation: "clear_account", Amount: 1, UserId: 2, NeededAccountIds: []uint{3, 4, 5}},
	// 		// {TypeOfAutomation: "repeat", UserId: 2},
	// 	},
	// }, true)

	// go automation.InitiateAutomation(helper.RandomString(20), [][]automation.AutomationLog{
	// 	{
	// 		{TypeOfAutomation: "wait", Amount: helper.RandomNumber(1), UserId: 2},
	// 		{TypeOfAutomation: "create_user", Amount: 1, UserId: 2},
	// 		{TypeOfAutomation: "wait", Amount: helper.RandomNumber(1), UserId: 2},
	// 		{TypeOfAutomation: "clear", Amount: 1, UserId: 2},
	// 		// {TypeOfAutomation: "repeat", UserId: 2},
	// 	},
	// }, true)

	e := echo.New()

	database.ConnectDb()

	// print(account.CheckBanned("nasadzhl46v1"))
	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	router.Setuprouter(e)

	e.Renderer = echoview.Default()
	e.Static("/static", "assets")
	e.File("/favicon.ico", "assets/favicon.ico")

	e.GET("/", func(c echo.Context) error {
		site_data, _, _ := helper.PageDataCreator(c, "Chito Tiktok", "", "", "", "", true, 0)

		return c.Render(http.StatusOK, "index", site_data)
	})

	e.GET("/signin", func(c echo.Context) error {
		site_data, _, _ := helper.PageDataCreator(c, "Chito Tiktok", "", "", "", "", true, 0)
		return c.Render(http.StatusOK, "sign_in", site_data)
	})

	e.GET("/signup", func(c echo.Context) error {
		site_data, _, _ := helper.PageDataCreator(c, "Chito Tiktok", "", "", "", "", true, 0)

		return c.Render(http.StatusOK, "sign_up", site_data)
	})

	e.GET("/404", func(c echo.Context) error {
		site_data, _, _ := helper.PageDataCreator(c, "Chito Tiktok", "404", "", "It looks like you are lost.. Want to go back?", "/", false, 0)

		return c.Render(http.StatusOK, "404",
			site_data)
	})

	e.GET("/home", func(c echo.Context) error {

		extra_value, user_id, err := helper.PageDataCreator(c, "Chito Tiktok", "", "", "", "", true, 2)
		if err != nil {
			c.Redirect(http.StatusSeeOther, "/logout")
		}
		var accounts []models.Account
		helper.Database.Db.Table("accounts").Joins("INNER JOIN users ON accounts.user_id = users.id").
			Where("users.id = ? And accounts.type_of_account = ? And accounts.cleared = 0", user_id, "disposable").
			Select("accounts.*").
			Find(&accounts)

		extra_value["body"] = helper.UserSidebar("/home")
		extra_value["accounts"] = accounts
		extra_value["total_accounts_count"] = len(accounts)
		extra_value["add"] = func(a int, b int) int {
			return a + b
		}

		return c.Render(http.StatusOK, "user/index", extra_value)
	})

	e.GET("/home/tiktok/posts/:screenName/:id", func(c echo.Context) error {
		extra_value, _, err := helper.PageDataCreator(c, "Chito Tiktok", "", "", "", "", true, 2)
		if err != nil {
			c.Redirect(http.StatusSeeOther, "/logout")
		}
		id := c.Param("id")
		screenName := c.Param("screenName")

		var posts []models.AccountPost

		helper.Database.Db.Preload("Post").Find(&posts, "account_id = ?", id)

		extra_value["body"] = helper.UserSidebar("")
		extra_value["posts"] = posts
		extra_value["ScreenName"] = screenName
		extra_value["id"] = id

		return c.Render(http.StatusOK, "user/post", extra_value)
	})

	e.GET("/home/posts", func(c echo.Context) error {

		extra_value, user_id, err := helper.PageDataCreator(c, "Chito Tiktok", "", "", "", "", true, 2)
		if err != nil {
			c.Redirect(http.StatusSeeOther, "/logout")
		}

		user := models.User{}
		helper.Database.Db.Preload("Posts").First(&user, "id = ?", user_id)

		extra_value["body"] = helper.UserSidebar("/home/posts")
		extra_value["posts"] = user.Posts
		return c.Render(http.StatusOK, "user/posts/all", extra_value)
	})

	e.GET("/home/posts/create", func(c echo.Context) error {
		extra_value, _, err := helper.PageDataCreator(c, "Create Post", "", "", "", "", true, 2)
		if err != nil {
			c.Redirect(http.StatusSeeOther, "/logout")
		}
		extra_value["body"] = helper.UserSidebar("/home/posts")

		// posts := post.GetStoredPost(1)
		return c.Render(http.StatusOK, "user/posts/create", extra_value)
	})

	e.GET("/home/names", func(c echo.Context) error {
		extra_value, user_id, err := helper.PageDataCreator(c, "Names", "", "", "", "", true, 2)
		if err != nil {
			c.Redirect(http.StatusSeeOther, "/logout")
		}
		extra_value["body"] = helper.UserSidebar("/home/names")

		names := []models.Name{}
		helper.Database.Db.Where("user_id = ?", user_id).Find(&names)
		extra_value["names"] = names

		return c.Render(http.StatusOK, "user/name", extra_value)
	})

	e.GET("/home/notification", func(c echo.Context) error {
		extra_value, user_id, err := helper.PageDataCreator(c, "Names", "", "", "", "", true, 2)
		if err != nil {
			c.Redirect(http.StatusSeeOther, "/logout")
		}
		extra_value["body"] = helper.UserSidebar("/home/notification")

		notifications := []models.Notification{}
		helper.Database.Db.Where("user_id = ?", user_id).Find(&notifications)
		helper.Database.Db.Model(&models.Notification{}).Where("seen = ? AND user_id = ?", false, user_id).Update("seen", true)
		extra_value["notifications"] = notifications

		return c.Render(http.StatusOK, "user/notification", extra_value)
	})

	e.GET("/home/proxies", func(c echo.Context) error {

		extra_value, user_id, err := helper.PageDataCreator(c, "Proxy", "", "", "", "", true, 2)
		if err != nil {
			c.Redirect(http.StatusSeeOther, "/logout")
		}
		extra_value["body"] = helper.UserSidebar("/home/proxies")

		names := []models.Proxy{}
		helper.Database.Db.Where("user_id = ?", user_id).Find(&names)
		extra_value["proxies"] = names

		return c.Render(http.StatusOK, "user/proxy", extra_value)
	})

	e.GET("/home/contained", func(c echo.Context) error {

		extra_value, user_id, err := helper.PageDataCreator(c, "Chito Tiktok", "", "", "", "", true, 2)
		if err != nil {
			c.Redirect(http.StatusSeeOther, "/logout")
		}
		accounts := []models.Account{}
		helper.Database.Db.Table("accounts").Joins("INNER JOIN users ON accounts.user_id = users.id").
			Where("users.id = ? And accounts.type_of_account = ? And accounts.cleared = 0", user_id, "contained").
			Select("accounts.*").
			Find(&accounts)
		extra_value["body"] = helper.UserSidebar("/home/contained")
		extra_value["accounts"] = accounts
		extra_value["total_accounts_count"] = len(accounts)
		extra_value["add"] = func(a int, b int) int {
			return a + b
		}

		return c.Render(http.StatusOK, "user/index", extra_value)
	})

	e.GET("/home/automations", func(c echo.Context) error {

		extra_value, user_id, err := helper.PageDataCreator(c, "Chito Tiktok", "", "", "", "", true, 2)
		if err != nil {
			c.Redirect(http.StatusSeeOther, "/logout")
		}
		extra_value["body"] = helper.UserSidebar("/home/automations")
		extra_value["automations"] = automation.GetUserAutomation(user_id)
		return c.Render(http.StatusOK, "user/automation/index", extra_value)
	})

	e.GET("/home/automation/create", func(c echo.Context) error {

		extra_value, user_id, err := helper.PageDataCreator(c, "Chito Tiktok", "", "", "", "", true, 2)
		if err != nil {
			c.Redirect(http.StatusSeeOther, "/logout")
		}

		posts := []models.Post{}

		helper.Database.Db.Where("user_id = ?", user_id).Find(&posts)

		extra_value["body"] = helper.UserSidebar("/home/automations")
		extra_value["posts"] = posts
		return c.Render(http.StatusOK, "user/automation/create", extra_value)
	})

	e.GET("/logout", func(c echo.Context) error {
		cookie := new(http.Cookie)
		cookie.Name = "toke"
		cookie.Value = ""
		cookie.Expires = time.Now().Add(time.Duration(-10) * time.Minute)
		c.SetCookie(cookie)
		return c.Redirect(http.StatusSeeOther, "/")
	})

	// Admin

	e.GET("/admin", func(c echo.Context) error {
		send_data, _, err := helper.PageDataCreator(c, "Chito Tiktok", "", "", "", "", true, 1)
		if err != nil {
			c.Redirect(http.StatusSeeOther, "/logout")
		}

		users := user.GetAnyVerifiedUsers(0)

		send_data["users"] = users
		return c.Render(http.StatusOK, "admin/index", send_data)
	})

	e.GET("/admin/verified", func(c echo.Context) error {
		send_data, _, err := helper.PageDataCreator(c, "Chito Tiktok", "404", "", "It looks like you are lost.. Want to go back?", "/", true, 1)

		if err != nil {
			c.Redirect(http.StatusSeeOther, "/logout")
		}

		users := user.GetAnyVerifiedUsers(1)
		send_data["users"] = users
		return c.Render(http.StatusOK, "admin/verified", send_data)
	})

	go automation.TimeSeries()

	// Start Server
	e.Logger.Fatal(e.Start(":1323"))
}
