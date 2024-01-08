package models

import (
	"github.com/google/uuid"
	"github.com/skip2/go-qrcode"
	"github.com/xlzd/gotp"
	"gorm.io/gorm"
)

type User struct {
	ID         uint   `gorm:"primaryKey"`
	UserUuid   string `gorm:"unique"`
	Email      string `gorm:"unique"`
	Password   string
	TotpActive bool
	TotpSecret string
}

type TotpInfo struct {
	uri    string
	qr     string
	secret string
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	u.UserUuid = uuid.New().String()
	return
}

// GenerateTotpInfo otpauth://totp/issuerName:email?secret=secretOfUser&issuer=NewsLogin
func (u *User) GenerateTotpInfo() (*TotpInfo, error) {
	uri := gotp.NewDefaultTOTP(
		u.TotpSecret,
	).ProvisioningUri(
		u.Email,
		"NewsLogin",
	)
	code, err := qrcode.New(uri, qrcode.Medium)
	if err != nil {
		return nil, err
	}
	return &TotpInfo{
		secret: u.TotpSecret,
		uri:    uri,
		qr:     code.Content,
	}, nil
}
