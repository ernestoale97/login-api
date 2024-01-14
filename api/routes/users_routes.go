package routes

import (
	"github.com/labstack/echo/v4"
	"login_api/api/controllers"
)

// UsersRoutes defines app user routes
func UsersRoutes(app *echo.Echo) {
	// generates TOTP secret and TOTP url
	app.GET("/users/:uuid/totp", controllers.GenerateTotp)
	// activates totp
	app.POST("/users/:uuid/totp", controllers.ActivateTotp)
	// disables totp
	app.DELETE("/users/:uuid/totp", controllers.DisableTotp)
}
