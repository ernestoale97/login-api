package models

type ActivateTotpInput struct {
	Totp int `json:"totp" validate:"required"` // actual generated CODE from GAuthenticator
}
