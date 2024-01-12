package controllers

import (
	"errors"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
	"log"
	"login_api/api/models"
	"login_api/pkg/config"
	jwtVerifier "login_api/pkg/jwt"
	"login_api/pkg/password_validator"
	"login_api/pkg/totp"
	"login_api/pkg/validation"
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
				"message": "Hubo un error desconocido",
			},
		)
	}
	var data models.LoginInput
	if err := c.Bind(&data); err != nil {
		log.Println("Error al validar los datos de entrada:", err.Error())
		return c.JSON(
			http.StatusBadRequest,
			&echo.Map{
				"message": "Error al validar los datos de entrada. Rectifíquelos",
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
				"message": "Hubo un error con la base de datos",
			},
		)
	}
	var user = models.User{}
	err = db.Where("email = ?", data.Email).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return c.JSON(
			http.StatusNotFound,
			&echo.Map{
				"message": "El usuario no existe",
			},
		)
	}
	if err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			&echo.Map{
				"message": "Hubo un error con la base de datos",
			},
		)
	}
	if !password_validator.CheckPasswordHash(data.Password, user.Password) {
		return c.JSON(
			http.StatusUnauthorized,
			&echo.Map{
				"message": "Credenciales inválidas",
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
					"message": "Hubo un error al generar el token",
				},
			)
		}
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
				"message": "Hubo un error al generar el token",
			},
		)
	}
	// user logged in and totp not active
	return c.JSON(
		http.StatusOK,
		&echo.Map{
			"data": &echo.Map{
				"access_token": jwtString,
				"email":        user.Email,
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
				"message": "Hubo un error desconocido",
			},
		)
	}
	var data models.VerifyTotpInput
	if err := c.Bind(&data); err != nil {
		log.Println("Error al validar los datos de entrada:", err.Error())
		return c.JSON(
			http.StatusBadRequest,
			&echo.Map{
				"message": "Error al validar los datos de entrada. Rectifíquelos",
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
				"message": "Hubo un error con la base de datos",
			},
		)
	}
	// validate token
	token, err := jwtVerifier.IsValidToken(c, "verify-totp")
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
				"message": "Hubo un error con la base de datos",
			},
		)
	}
	var sub = token.Claims.(jwt.MapClaims)["sub"]
	var user = models.User{}
	err = db.Where("user_uuid = ?", sub).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return c.JSON(
			http.StatusNotFound,
			&echo.Map{
				"message": "El usuario no existe",
			},
		)
	}
	if err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			&echo.Map{
				"message": "Hubo un error con la base de datos",
			},
		)
	}
	if !user.TotpActive {
		return c.JSON(
			http.StatusForbidden,
			&echo.Map{
				"message": "TOTP del usuario no activo",
			},
		)
	}
	// validate the totp against token owner secret
	if !totp.Verify(user.TotpSecret, strconv.Itoa(data.Totp)) {
		return c.JSON(
			http.StatusUnauthorized,
			&echo.Map{
				"message": "TOTP inválido",
			},
		)
	}
	jwtString, err := user.GenerateJWT(false)
	if err != nil {
		return c.JSON(
			http.StatusBadGateway,
			&echo.Map{
				"message": "Hubo un error al generar el token",
			},
		)
	}
	// user logged in and totp is active
	return c.JSON(
		http.StatusOK,
		&echo.Map{
			"data": &echo.Map{
				"access_token": jwtString,
				"email":        user.Email,
				"mfa_required": false,
			},
		},
	)
}
