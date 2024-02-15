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

	// testing
	app.GET("/locales", controllers.GetLocales)
	app.GET("/register-config", controllers.GetConfig)
	app.POST("/register", controllers.Register)
	app.POST("/login", controllers.AppLogin)
	app.POST("/auth/userinfo", controllers.UserInfo)
	app.GET("/auth/listActiveWallet", controllers.ActiveWallets)
	app.GET("/auth/transactionHistoryByWallet", controllers.TransactionHistoryByWallet)
	app.GET("/auth/getListChallenge", controllers.ChallengesList)
	app.POST("/auth/getMetricDetails", controllers.ChallengeMetrics)

}
