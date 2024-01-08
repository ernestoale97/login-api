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
	"login_api/pkg/totp"
	"login_api/pkg/validation"
	"login_api/storage"
	"net/http"
)

func ActivateTotp(c echo.Context) error {
	// load config
	_, err := config.LoadConfig()
	if err != nil {
		log.Println("Error loading config:", err.Error())
		return c.JSON(
			http.StatusInternalServerError,
			&echo.Map{
				"success": false,
				"message": "There was an unknown error",
			},
		)
	}
	var data models.SignupInput
	if err := c.Bind(&data); err != nil {
		log.Println("Error on validating request data:", err.Error())
		return c.JSON(
			http.StatusBadRequest,
			&echo.Map{
				"success": false,
				"message": "There was an error validating your data. Please retry",
			},
		)
	}
	validate := validation.GetInputValidationInstance()
	if err := validate.Struct(data); err != nil {
		return c.JSON(
			http.StatusBadRequest,
			&echo.Map{
				"success": false,
				"message": "There was an error in data input. Fix them",
			},
		)
	}
	// connect db
	db, err := storage.ConnectDB()
	if err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			&echo.Map{
				"success": false,
				"message": "There was an error with database",
			},
		)
	}
	var user = models.User{}
	err = db.Where("email = ?", data.Email).First(&user).Error
	if err == nil {
		return c.JSON(
			http.StatusBadRequest,
			&echo.Map{
				"success": false,
				"message": "User email already exists",
			},
		)
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return c.JSON(
			http.StatusInternalServerError,
			&echo.Map{
				"success": false,
				"message": "There was an error with database",
			},
		)
	}
	hash, err := password_validator.HashPassword(data.Password)
	if err != nil {
		fmt.Println("Error creating password hash:", err.Error())
		return c.JSON(
			http.StatusBadGateway,
			&echo.Map{
				"success": false,
				"message": "There was an error creating the user",
			},
		)
	}
	user = models.User{
		Email:      data.Email,
		Password:   hash,
		TotpActive: false,
	}
	result := db.Create(&user)
	if result.Error != nil || result.RowsAffected <= 0 {
		return c.JSON(
			http.StatusBadGateway,
			&echo.Map{
				"success": false,
				"message": "There was an error creating the user",
			},
		)
	}
	// signup successful
	return c.JSON(
		http.StatusCreated,
		&echo.Map{
			"success": true,
			"data": &echo.Map{
				"user": &echo.Map{
					"user_uuid": &user.UserUuid,
					"email":     &user.Email,
				},
			},
		},
	)
}

func Signup(c echo.Context) error {
	// load config
	_, err := config.LoadConfig()
	if err != nil {
		log.Println("Error loading config:", err.Error())
		return c.JSON(
			http.StatusInternalServerError,
			&echo.Map{
				"success": false,
				"message": "There was an unknown error",
			},
		)
	}
	var data models.SignupInput
	if err := c.Bind(&data); err != nil {
		log.Println("Error on validating request data:", err.Error())
		return c.JSON(
			http.StatusBadRequest,
			&echo.Map{
				"success": false,
				"message": "There was an error validating your data. Please retry",
			},
		)
	}
	validate := validation.GetInputValidationInstance()
	if err := validate.Struct(data); err != nil {
		return c.JSON(
			http.StatusBadRequest,
			&echo.Map{
				"success": false,
				"message": "There was an error in data input. Fix them",
			},
		)
	}
	// connect db
	db, err := storage.ConnectDB()
	if err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			&echo.Map{
				"success": false,
				"message": "There was an error with database",
			},
		)
	}
	var user = models.User{}
	err = db.Where("email = ?", data.Email).First(&user).Error
	if err == nil {
		return c.JSON(
			http.StatusBadRequest,
			&echo.Map{
				"success": false,
				"message": "User email already exists",
			},
		)
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return c.JSON(
			http.StatusInternalServerError,
			&echo.Map{
				"success": false,
				"message": "There was an error with database",
			},
		)
	}
	hash, err := password_validator.HashPassword(data.Password)
	if err != nil {
		fmt.Println("Error creating password hash:", err.Error())
		return c.JSON(
			http.StatusBadGateway,
			&echo.Map{
				"success": false,
				"message": "There was an error creating the user",
			},
		)
	}
	user = models.User{
		Email:      data.Email,
		Password:   hash,
		TotpActive: false,
		TotpSecret: totp.GenerateSecret(),
	}
	result := db.Create(&user)
	if result.Error != nil || result.RowsAffected <= 0 {
		return c.JSON(
			http.StatusBadGateway,
			&echo.Map{
				"success": false,
				"message": "There was an error creating the user",
			},
		)
	}
	// signup successful
	return c.JSON(
		http.StatusCreated,
		&echo.Map{
			"success": true,
			"data": &echo.Map{
				"user": &echo.Map{
					"user_uuid": &user.UserUuid,
					"email":     &user.Email,
				},
			},
		},
	)
}
