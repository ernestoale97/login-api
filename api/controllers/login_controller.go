package controllers

import (
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
	"log"
	"login_api/api/models"
	"login_api/pkg/config"
	"login_api/pkg/password_validator"
	validation "login_api/pkg/validation"
	"login_api/storage"
	"net/http"
)

func Login(c echo.Context) error {
	// load config
	_, err := config.LoadConfig()
	if err != nil {
		log.Println("Error loading config:", err.Error())
		return c.JSON(
			http.StatusInternalServerError,
			&echo.Map{
				"message": "There was an unknown error",
			},
		)
	}
	var data models.LoginInput
	if err := c.Bind(&data); err != nil {
		log.Println("Error on validating request data:", err.Error())
		return c.JSON(
			http.StatusBadRequest,
			&echo.Map{
				"message": "There was an error validating your data. Please retry",
			},
		)
	}
	validate := validation.GetInputValidationInstance()
	if err := validate.Struct(data); err != nil {
		return err
	}
	// connect db
	db, err := storage.ConnectDB()
	if err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			&echo.Map{
				"message": "There was an error with database",
			},
		)
	}
	var user = models.User{}
	err = db.Where("email = ?", data.Email).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return c.JSON(
			http.StatusUnauthorized,
			&echo.Map{
				"message": "User not exists",
			},
		)
	}
	if err != nil {
		return c.JSON(
			http.StatusBadGateway,
			&echo.Map{
				"message": "There was an error with database",
			},
		)
	}
	if !password_validator.CheckPasswordHash(data.Password, user.Password) {
		return c.JSON(
			http.StatusUnauthorized,
			&echo.Map{
				"message": "Wrong credentials",
			},
		)
	}
	if user.TotpActive {
		// if user has totp then send token but with scope only for verifyOtp
		return c.JSON(
			http.StatusOK,
			&echo.Map{
				"data": &echo.Map{
					"mfa_token":    "token with only scope verifyOtp",
					"mfa_required": true,
				},
			},
		)
	}
	// user logged in and totp not active
	return c.JSON(
		http.StatusOK,
		&echo.Map{
			"data": &echo.Map{
				"access_token": "okkkkkkk",
				"mfa_required": false,
			},
		},
	)
}

func VerifyTotp(c echo.Context) error {
	// load config
	_, err := config.LoadConfig()
	if err != nil {
		log.Println("Error loading config:", err.Error())
		return c.JSON(
			http.StatusInternalServerError,
			&echo.Map{
				"message": "There was an unknown error",
			},
		)
	}
	var data models.VerifyTotpInput
	if err := c.Bind(&data); err != nil {
		log.Println("Error on validating request data:", err.Error())
		return c.JSON(
			http.StatusBadRequest,
			&echo.Map{
				"message": "There was an error validating your data. Please retry",
			},
		)
	}
	validate := validation.GetInputValidationInstance()
	if err := validate.Struct(data); err != nil {
		return err
	}
	// connect db
	db, err := storage.ConnectDB()
	if err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			&echo.Map{
				"message": "There was an error with database",
			},
		)
	}
	fmt.Println(db)
	// here check the totp of user owner of token with scope verifyTotp
	// if success then return access_token with scope user
	return c.JSON(
		http.StatusOK,
		&echo.Map{
			"data": &echo.Map{
				"access_token": "okkkkkkk",
			},
		},
	)
}
