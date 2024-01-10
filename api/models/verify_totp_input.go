package models

type VerifyTotpInput struct {
	Totp int `json:"totp" validate:"required"` // actual generated CODE from GAuthenticator
}
