package controllers

import (
	"errors"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
	"log"
	"login_api/api/models"
	"login_api/pkg/config"
	jwt_verifier "login_api/pkg/jwt"
	"login_api/pkg/password_validator"
	"login_api/pkg/totp"
	validation "login_api/pkg/validation"
	"login_api/storage"
	"net/http"
	"strconv"
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
		jwtString, err := user.GenerateJWT(true)
		if err != nil {
			return c.JSON(
				http.StatusBadGateway,
				&echo.Map{
					"message": "There was an error generating token for user",
				},
			)
		}
		log.Println(jwtString)
		return c.JSON(
			http.StatusOK,
			&echo.Map{
				"data": &echo.Map{
					"mfa_token":    jwtString,
					"mfa_required": true,
				},
			},
		)
	}
	jwtString, err := user.GenerateJWT(false)
	if err != nil {
		return c.JSON(
			http.StatusBadGateway,
			&echo.Map{
				"message": "There was an error generating token for user",
			},
		)
	}
	log.Println(jwtString)
	// user logged in and totp not active
	return c.JSON(
		http.StatusOK,
		&echo.Map{
			"data": &echo.Map{
				"access_token": jwtString,
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
	_, err = storage.ConnectDB()
	if err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			&echo.Map{
				"message": "There was an error with database",
			},
		)
	}
	// validate token
	token, err := jwt_verifier.IsValidToken(c, "verify-totp")
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
	var user = models.User{}
	err = db.Where("user_uuid = ?", sub).First(&user).Error
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
	if !user.TotpActive {
		return c.JSON(
			http.StatusForbidden,
			&echo.Map{
				"message": "User Totp is not active",
			},
		)
	}
	// validate the totp against token owner secret
	if !totp.Verify(user.TotpSecret, strconv.Itoa(data.Totp)) {
		return c.JSON(
			http.StatusUnauthorized,
			&echo.Map{
				"message": "Invalid totp",
			},
		)
	}
	jwtString, err := user.GenerateJWT(false)
	if err != nil {
		return c.JSON(
			http.StatusBadGateway,
			&echo.Map{
				"message": "There was an error generating token for user",
			},
		)
	}
	// user logged in and totp is active
	return c.JSON(
		http.StatusOK,
		&echo.Map{
			"data": &echo.Map{
				"access_token": jwtString,
				"mfa_required": false,
			},
		},
	)
}
