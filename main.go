package main

import (
	"beepbop/automation"
	"beepbop/controller/user"
	"beepbop/database"
	"beepbop/helper"
	"beepbop/models"
	"beepbop/router"
	"fmt"
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
		return c.Render(http.StatusOK, "index", helper.PageDataCreator(c, "Chito Tiktok", "", "", "", "", true, 0))
	})

	e.GET("/signin", func(c echo.Context) error {
		return c.Render(http.StatusOK, "sign_in", helper.PageDataCreator(c, "Chito Tiktok", "", "", "", "", true, 0))
	})

	e.GET("/signup", func(c echo.Context) error {
		return c.Render(http.StatusOK, "sign_up", helper.PageDataCreator(c, "Chito Tiktok", "", "", "", "", true, 0))
	})

	e.GET("/404", func(c echo.Context) error {
		return c.Render(http.StatusOK, "404",
			helper.PageDataCreator(c, "Chito Tiktok", "404", "", "It looks like you are lost.. Want to go back?", "/", false, 0))
	})

	e.GET("/home", func(c echo.Context) error {
		user_claims, err := helper.JWT(c)

		if err != nil || user_claims.Role != 2 {
			c.Redirect(http.StatusSeeOther, "/404")
		}

		var accounts []models.Account

		helper.Database.Db.Table("accounts").Joins("INNER JOIN users ON accounts.user_id = users.id").
			Where("users.token = ? And accounts.type_of_account = ? And accounts.cleared = 0", user_claims.Token, "disposable").
			Select("accounts.*").
			Find(&accounts)

		extra_value := helper.PageDataCreator(c, "Chito Tiktok", "", "", "", "", true, 2)
		extra_value["body"] = helper.UserSidebar("/home")
		extra_value["accounts"] = accounts
		extra_value["total_accounts_count"] = len(accounts)
		extra_value["add"] = func(a int, b int) int {
			return a + b
		}

		return c.Render(http.StatusOK, "user/index", extra_value)
	})

	e.GET("/home/tiktok/posts/:screenName/:id", func(c echo.Context) error {
		user_claims, err := helper.JWT(c)

		if err != nil || user_claims.Role != 2 {
			c.Redirect(http.StatusSeeOther, "/404")
		}

		id := c.Param("id")
		screenName := c.Param("screenName")

		var posts []models.AccountPost

		helper.Database.Db.Preload("Post").Find(&posts, "account_id = ?", id)

		extra_value := helper.PageDataCreator(c, "Chito Tiktok", "", "", "", "", true, 2)
		extra_value["body"] = helper.UserSidebar("")
		extra_value["posts"] = posts
		extra_value["ScreenName"] = screenName
		extra_value["id"] = id

		return c.Render(http.StatusOK, "user/post", extra_value)
	})

	e.GET("/home/posts", func(c echo.Context) error {
		user_claims, err := helper.JWT(c)

		if err != nil || user_claims.Role != 2 {
			c.Redirect(http.StatusSeeOther, "/404")
		}

		user := models.User{}
		helper.Database.Db.Preload("Posts").First(&user, "token = ?", user_claims.Token)

		extra_value := helper.PageDataCreator(c, "Chito Tiktok", "", "", "", "", true, 2)
		extra_value["body"] = helper.UserSidebar("/home/posts")
		extra_value["posts"] = user.Posts
		return c.Render(http.StatusOK, "user/posts/all", extra_value)
	})

	e.GET("/home/posts/create", func(c echo.Context) error {
		user_claims, err := helper.JWT(c)

		if err != nil || user_claims.Role != 2 {
			c.Redirect(http.StatusSeeOther, "/404")
		}

		extra_value := helper.PageDataCreator(c, "Create Post", "", "", "", "", true, 2)
		extra_value["body"] = helper.UserSidebar("/home/posts")

		// posts := post.GetStoredPost(1)
		return c.Render(http.StatusOK, "user/posts/create", extra_value)
	})

	e.GET("/home/contained", func(c echo.Context) error {
		user_claims, err := helper.JWT(c)

		if err != nil || user_claims.Role != 2 {
			c.Redirect(http.StatusSeeOther, "/404")
		}

		var accounts []models.Account

		helper.Database.Db.Table("accounts").Joins("INNER JOIN users ON accounts.user_id = users.id").
			Where("users.token = ? And accounts.type_of_account = ? And accounts.cleared = 0", user_claims.Token, "contained").
			Select("accounts.*").
			Find(&accounts)

		extra_value := helper.PageDataCreator(c, "Chito Tiktok", "", "", "", "", true, 2)
		extra_value["body"] = helper.UserSidebar("/home/contained")
		extra_value["accounts"] = accounts
		extra_value["total_accounts_count"] = len(accounts)
		extra_value["add"] = func(a int, b int) int {
			return a + b
		}

		return c.Render(http.StatusOK, "user/index", extra_value)
	})

	e.GET("/home/automations", func(c echo.Context) error {
		user_claims, err := helper.JWT(c)

		if err != nil || user_claims.Role != 2 {
			c.Redirect(http.StatusSeeOther, "/404")
		}

		user := models.User{}

		helper.Database.Db.Select("id").Where("token = ?", user_claims.Token).Find(&user)

		fmt.Print(user.Id)

		extra_value := helper.PageDataCreator(c, "Chito Tiktok", "", "", "", "", true, 2)
		extra_value["body"] = helper.UserSidebar("/home/automations")
		extra_value["automations"] = automation.GetUserAutomation(user.Id)
		return c.Render(http.StatusOK, "user/automation/index", extra_value)
	})

	e.GET("/home/automation/create", func(c echo.Context) error {
		user_claims, err := helper.JWT(c)

		if err != nil || user_claims.Role != 2 {
			c.Redirect(http.StatusSeeOther, "/404")
		}

		user := models.User{}

		helper.Database.Db.Preload("Posts").First(&user, "token = ?", user_claims.Token)

		extra_value := helper.PageDataCreator(c, "Chito Tiktok", "", "", "", "", true, 2)
		extra_value["body"] = helper.UserSidebar("/home/automations")
		extra_value["posts"] = user.Posts
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
		user_claims, err := helper.JWT(c)

		if err != nil || user_claims.Role != 1 {
			c.Redirect(http.StatusSeeOther, "/404")
		}

		send_data := helper.PageDataCreator(c, "Chito Tiktok", "", "", "", "", true, 1)

		users := user.GetAnyVerifiedUsers(0)

		send_data["users"] = users
		return c.Render(http.StatusOK, "admin/index", send_data)
	})

	e.GET("/admin/verified", func(c echo.Context) error {
		user_claims, err := helper.JWT(c)

		if err != nil || user_claims.Role != 1 {
			c.Redirect(http.StatusSeeOther, "/404")
		}

		send_data := helper.PageDataCreator(c, "Chito Tiktok", "404", "", "It looks like you are lost.. Want to go back?", "/", true, 1)

		users := user.GetAnyVerifiedUsers(1)
		send_data["users"] = users
		return c.Render(http.StatusOK, "admin/verified", send_data)
	})

	go automation.TimeSeries()

	// Start Server
	e.Logger.Fatal(e.Start(":1323"))
}
