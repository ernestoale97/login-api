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
	app.GET("/generalConfig", controllers.GeneralConfig)
	app.POST("/register", controllers.Register)
	app.POST("/login", controllers.AppLogin)
	app.POST("/auth/userinfo", controllers.UserInfo)
	app.GET("/auth/listActiveWallet", controllers.ActiveWallets)
	app.POST("/auth/transactionHistoryByWallet", controllers.TransactionHistoryByWallet)
	app.GET("/auth/getListChallenge", controllers.ChallengesList)
	app.POST("/auth/getMetricDetails", controllers.ChallengeMetrics)
	app.GET("/auth/generate2FAToken", controllers.Generate2faToken)
	app.POST("/auth/verify2FACode", controllers.Verify2faToken)
	app.POST("/auth/disable2FA", controllers.Disable2faToken)
	app.GET("/auth/enableEmailVerification", controllers.EnableEmailVerification)
	app.GET("/auth/disableEmailVerification", controllers.DisableEmailVerification)
	app.GET("/auth/getBanner", controllers.GetBanner)
	app.GET("/auth/loadAlerts", controllers.LoadAlerts)
	app.GET("/auth/getModules", controllers.GetModules)
	app.POST("/auth/updatedImgprofile", controllers.UpdateImageProfile)

}
