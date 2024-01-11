package controllers

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
	"log"
	"login_api/api/models"
	"login_api/pkg/config"
	jwt_verifier "login_api/pkg/jwt"
	"login_api/pkg/password_validator"
	"login_api/pkg/totp"
	"login_api/pkg/validation"
	"login_api/storage"
	"net/http"
)

// ActivateTotp headers: token of user; params:
func ActivateTotp(c echo.Context) error {
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
	var data models.SignupInput
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
		return c.JSON(
			http.StatusBadRequest,
			&echo.Map{
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
				"message": "User email already exists",
			},
		)
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return c.JSON(
			http.StatusInternalServerError,
			&echo.Map{
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
				"message": "There was an error creating the user",
			},
		)
	}
	// signup successful
	return c.JSON(
		http.StatusNoContent,
		nil,
	)
}

func GenerateTotp(c echo.Context) error {
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
	// validate token
	token, err := jwt_verifier.IsValidToken(c, "user")
	if err != nil {
		return c.JSON(
			http.StatusUnauthorized,
			&echo.Map{
				"message": err.Error(),
			},
		)
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
	var sub = token.Claims.(jwt.MapClaims)["sub"]
	if sub != c.Param("uuid") {
		return c.JSON(
			http.StatusForbidden,
			&echo.Map{
				"message": "UUID not match with token",
			},
		)
	}
	var user = models.User{}
	err = db.Where("user_uuid = ?", sub).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return c.JSON(
			http.StatusInternalServerError,
			&echo.Map{
				"message": "User not exists",
			},
		)
	}
	if err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			&echo.Map{
				"message": "There was an error with database",
			},
		)
	}
	totpInfo, err := user.GenerateTotpInfo()
	log.Println(totpInfo)
	if err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			&echo.Map{
				"message": "There was an error generating totp info",
			},
		)
	}
	// generated totp info
	return c.JSON(
		http.StatusOK,
		totpInfo,
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
				"message": "There was an error validating your data. Please retry",
			},
		)
	}
	validate := validation.GetInputValidationInstance()
	if err := validate.Struct(data); err != nil {
		return c.JSON(
			http.StatusBadRequest,
			&echo.Map{
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
				"message": "User email already exists",
			},
		)
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return c.JSON(
			http.StatusInternalServerError,
			&echo.Map{
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
				"message": "There was an error creating the user",
			},
		)
	}
	// signup successful
	return c.JSON(
		http.StatusCreated,
		&echo.Map{
			"data": &echo.Map{
				"user": &echo.Map{
					"user_uuid": &user.UserUuid,
					"email":     &user.Email,
				},
			},
		},
	)
}
