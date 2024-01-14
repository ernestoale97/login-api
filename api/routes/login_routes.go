package routes

import (
	"github.com/labstack/echo/v4"
	"login_api/api/controllers"
)

func LoginRoutes(app *echo.Echo) {
	app.POST("/login", controllers.Login)
	app.GET("/logout", controllers.Logout)
	app.POST("/signup", controllers.Signup)
	app.POST("/verify-totp", controllers.VerifyTotp)
}
