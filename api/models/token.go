package models

import (
	"login_api/storage"
	"time"
)

type UserToken struct {
	Token
	Sub        string `json:"sub"`
	AccessUuid string `json:"access_uuid"`
	Email      string `json:"email"`
	TotpActive bool   `json:"totp_active"`
}

type Token struct {
	Typ   string `json:"typ"`
	Iat   int64  `json:"iat"`
	Nbf   int64  `json:"nbf"`
	Iss   string `json:"iss"`
	Exp   int64  `json:"exp"`
	Scope string `json:"scope"`
}

func (u *UserToken) CreateAuth() error {
	at := time.Unix(u.Exp, 0) //converting Unix to UTC(to Time object)
	now := time.Now()
	client, err := storage.GetRedis()
	if err != nil {
		return err
	}
	errAccess := client.Set(u.AccessUuid, u.Sub, at.Sub(now)).Err()
	if errAccess != nil {
		return errAccess
	}
	return nil
}
