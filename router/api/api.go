package api

import (
	accountController "beepbop/controller/account"
	postController "beepbop/controller/post"
	userController "beepbop/controller/user"

	"github.com/labstack/echo/v4"
)

// func hello_api(c echo.Context) error {
// 	return c.String(http.StatusOK, "Hello, Api!")
// }

func Api(e *echo.Echo) {
	e.POST("/home/posts/create", postController.CreatePost)
	e.GET("/home/posts/delete/:id", postController.DeletePost)

	e.POST("/home/automation/create", userController.CreateAutomation)
	e.GET("/home/automation/stop/:id", userController.StopAutomation)

	e.POST("/home/account/refresh", accountController.RefreshAccount)
	e.GET("/home/account/clear/:id", accountController.ClearAccount)
	e.POST("/home/account/contain", accountController.ContainedAccount)
	e.POST("/home/account/automate", accountController.AutomateDisplay)
	e.GET("/home/tiktok/post/delete/:tiktok_id/:id/:postid", accountController.RemovePost)

	e.POST("/home/name/create", accountController.NameAdd)
	e.POST("/home/proxy/create", accountController.ProxyAdd)
	e.GET("/home/proxy/refresh/:id", accountController.ProxyRefresh)
	e.GET("/home/proxy/delete/:id", accountController.ProxyDelete)

	g := e.Group("/user")
	// g.GET("/", hello_api)
	g.POST("/signin", userController.Login)
	g.POST("/signup", userController.Signup)
	g.GET("/accept/:id", userController.Accept)
	g.GET("/delete/:id", userController.Delete)
	g.POST("/check_username", userController.CheckUsername)

	// g.POST("/verify_otp", otpcontroller.VerifyOtp)

	g.GET("/me", userController.Me)
}
