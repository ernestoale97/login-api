package routes

import (
	"github.com/labstack/echo/v4"
	"login_api/api/controllers"
)

func LoginRoutes(app *echo.Echo) {
	app.POST("/login", controllers.Login)
	app.POST("/signup", controllers.Signup)

	// verify totp using a verifyTotp scoped token
	app.POST("/verify-totp", controllers.VerifyTotp)
}
