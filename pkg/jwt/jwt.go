package jwt

import (
	"crypto/x509"
	"encoding/pem"
	"errors"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"log"
	"login_api/storage"
	"os"
	"strings"
)

func getTokenPayload(token string) (string, error) {
	if len(strings.Trim(token, " ")) == 0 {
		return "", errors.New("invalid token")
	}
	tokens := strings.Split(token, " ")
	if len(tokens) < 2 {
		return "", errors.New("invalid token")
	}
	return tokens[1], nil
}

func getKey() (any, error) {
	publicKeyPath := "./env/public_key.pem"
	keyData, err := os.ReadFile(publicKeyPath)
	if err != nil {
		log.Println("error reading key file from assets", err.Error())
		return nil, err
	}
	publicKeyBlock, _ := pem.Decode(keyData)
	if publicKeyBlock == nil {
		return nil, errors.New("error decoding pem with keyData")
	}
	key, err := x509.ParsePKIXPublicKey(publicKeyBlock.Bytes)
	if err != nil {
		log.Println("error parsing key from key block bytes", err.Error())
		return nil, errors.New("error_key")
	}
	return key, nil
}

// function for parsing, validating and get a jwt.Token object
func parseValidateAndGetToken(bToken string) (*jwt.Token, error) {
	tokenPayload, err := getTokenPayload(bToken)
	if err != nil {
		return nil, err
	}
	key, err := getKey()
	if err != nil {
		return nil, err
	}
	token, err := jwt.Parse(tokenPayload, func(token *jwt.Token) (interface{}, error) {
		return key, nil
	})
	if err != nil {
		if err.Error() == "crypto/rsa: verification error" {
			return nil, errors.New("invalid token")
		}
		return nil, err
	}
	return token, nil
}

func checkAccessUuid(accessUuid string) error {
	client, err := storage.GetRedis()
	if err != nil {
		return err
	}
	_, err = client.Get(accessUuid).Result()
	if err != nil {
		return err
	}
	return nil
}

func IsValidToken(c echo.Context, scope string) (*jwt.Token, error) {
	tokenHeader := c.Request().Header.Get("authorization")
	// parse and validate token and receive an jwt.Token object
	token, err := parseValidateAndGetToken(tokenHeader)
	if err != nil {
		return nil, err
	}
	// if token is valid return no error
	if token.Valid {
		claims := token.Claims.(jwt.MapClaims)
		if scope == "user" {
			err := checkAccessUuid(claims["access_uuid"].(string))
			if err != nil {
				return nil, errors.New("token is no longer valid")
			}
		}
		if claims["scope"] != scope {
			return nil, errors.New("not valid scope")
		}
		if claims["iss"] != "login-news-api" {
			return nil, errors.New("not valid issuer")
		}
		return token, nil
	} else if ve, ok := err.(*jwt.ValidationError); ok {
		if ve.Errors&jwt.ValidationErrorMalformed != 0 {
			// if token is malformed
			return nil, errors.New("token invalid. Token is malformed")
		} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
			// if token is expired or not active yet due to claims nbf and iat
			return nil, errors.New("token is expired or not active yet")
		} else {
			// if not known error return it
			return nil, err
		}
	} else {
		// if not known error return it
		return nil, err
	}
}
