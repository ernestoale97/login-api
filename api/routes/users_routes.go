package routes

import (
	"github.com/labstack/echo/v4"
	"login_api/api/controllers"
)

func UsersRoutes(app *echo.Echo) {
	// generates TOTP secret and TOTP url
	app.GET("/users/:uuid/totp", controllers.Login)
	// activates otp
	app.POST("/users/:uuid/totp", controllers.ActivateTotp)
}
