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
	jwtVerifier "login_api/pkg/jwt"
	"login_api/pkg/password_validator"
	"login_api/pkg/totp"
	"login_api/pkg/validation"
	"login_api/storage"
	"net/http"
	"strconv"
)

func Signup(c echo.Context) error {
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
	var data models.SignupInput
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
		return c.JSON(
			http.StatusBadRequest,
			&echo.Map{
				"message": "Error al validar los datos de entrada. Rectifíquelos",
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
	var user = models.User{}
	err = db.Where("email = ?", data.Email).First(&user).Error
	if err == nil {
		return c.JSON(
			http.StatusBadRequest,
			&echo.Map{
				"message": "El correo ya existe",
			},
		)
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return c.JSON(
			http.StatusInternalServerError,
			&echo.Map{
				"message": "Hubo un error con la base de datos",
			},
		)
	}
	hash, err := password_validator.HashPassword(data.Password)
	if err != nil {
		fmt.Println("Error creating password hash:", err.Error())
		return c.JSON(
			http.StatusBadGateway,
			&echo.Map{
				"message": "Hubo un error al crear el usuario",
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
				"message": "Hubo un error al crear el usuario",
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

func GenerateTotp(c echo.Context) error {
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
	// validate token
	token, err := jwtVerifier.IsValidToken(c, "user")
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
	if sub != c.Param("uuid") {
		return c.JSON(
			http.StatusForbidden,
			&echo.Map{
				"message": "UUID no coincide con el token",
			},
		)
	}
	var user = models.User{}
	err = db.Where("user_uuid = ?", sub).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return c.JSON(
			http.StatusInternalServerError,
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
	totpInfo, err := user.GenerateTotpInfo()
	if err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			&echo.Map{
				"message": "Hubo un error al generar la información del TOTP",
			},
		)
	}
	// generated totp info
	return c.JSON(
		http.StatusOK,
		totpInfo,
	)
}

func ActivateTotp(c echo.Context) error {
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
	var data models.ActivateTotpInput
	if err := c.Bind(&data); err != nil {
		log.Println("Error al validar los datos de entrada:", err.Error())
		return c.JSON(
			http.StatusBadRequest,
			&echo.Map{
				"message": "Error al validar los datos de entrada. Rectifíquelos",
			},
		)
	}
	// validate data input fields
	validate := validation.GetInputValidationInstance()
	if err := validate.Struct(data); err != nil {
		return err
	}
	// validate token: scope, expiration, owner and signature
	token, err := jwtVerifier.IsValidToken(c, "user")
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
			http.StatusBadGateway,
			&echo.Map{
				"message": "Hubo un error con la base de datos",
			},
		)
	}
	// compare url user_uuid with token owner
	var sub = token.Claims.(jwt.MapClaims)["sub"]
	if sub != c.Param("uuid") {
		return c.JSON(
			http.StatusForbidden,
			&echo.Map{
				"message": "UUID no coincide con el token",
			},
		)
	}
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
	// now check received totp code against user secret
	if !totp.Verify(user.TotpSecret, strconv.Itoa(data.Totp)) {
		return c.JSON(
			http.StatusUnauthorized,
			&echo.Map{
				"message": "TOTP inválido",
			},
		)
	}
	// code is valid so set 2fa activated successfully
	res := db.Model(&user).Where("user_uuid = ?", sub).Update("totp_active", true)
	if res.RowsAffected <= 0 {
		return c.JSON(
			http.StatusBadGateway,
			&echo.Map{
				"message": "Hubo un error activando el 2FA",
			},
		)
	}
	return c.JSON(
		http.StatusNoContent,
		nil,
	)
}

func DisableTotp(c echo.Context) error {
	// TODO implement disable totp
	// TODO on disable remove user TOTP secret (it revokes all generated new keys from old devices)
	return nil
}
