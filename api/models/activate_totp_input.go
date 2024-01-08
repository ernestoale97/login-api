package models

type ActivateTotpInput struct {
	Secret string `json:"secret" validate:"required"` // secret url in base64 encoding
	Code   string `json:"code" validate:"required"`   // actual generated CODE from GAuthenticator
}
