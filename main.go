package main

import (
	"crypto/tls"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"log"
	"login_api/api/middlewares"
	"login_api/api/models"
	"login_api/api/routes"
	"login_api/pkg/config"
	"login_api/storage"
	"net/http"
)

func main() {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{
		InsecureSkipVerify: true,
	}
	// Echo instance
	app := echo.New()
	// Adding logger
	app.Use(middlewares.LogRequest)
	// Adding panic recover
	app.Use(middleware.Recover())
	// Routes
	routes.LoginRoutes(app)
	routes.UsersRoutes(app)
	// Allow from all origins
	app.Use(middleware.CORSWithConfig(
		middleware.CORSConfig{
			AllowOrigins: []string{"*"},
			AllowHeaders: []string{
				echo.HeaderOrigin,
				echo.HeaderContentType,
				echo.HeaderAccept,
			},
		},
	))
	conf, _ := config.LoadConfig()
	// Perform database migration
	db, err := storage.ConnectDB()
	if err != nil {
		log.Fatal(err)
	}
	// auto migrate users
	err = db.AutoMigrate(&models.User{})
	if err != nil {
		log.Fatal(err)
	}
	// Start server
	app.Logger.Fatal(app.Start(fmt.Sprintf(":%s", conf.AppPort)))
}
