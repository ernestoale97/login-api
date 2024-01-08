package models

type VerifyTotpInput struct {
	Totp string `json:"totp" validate:"required"` // actual generated CODE from GAuthenticator
}
