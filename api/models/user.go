package models

import (
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"github.com/google/uuid"
	"github.com/skip2/go-qrcode"
	"github.com/xlzd/gotp"
	"gopkg.in/square/go-jose.v2"
	"gopkg.in/square/go-jose.v2/jwt"
	"gorm.io/gorm"
	"log"
	"os"
	"time"
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
	Uri    string `json:"uri"`
	Qr     string `json:"qr"`
	Secret string `json:"secret"`
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	u.UserUuid = uuid.New().String()
	return
}

// GenerateTotpInfo otpauth://totp/issuerName:email?secret=secretOfUser&issuer=NewsLogin
func (u *User) GenerateTotpInfo() (TotpInfo, error) {
	uri := gotp.NewDefaultTOTP(
		u.TotpSecret,
	).ProvisioningUri(
		u.Email,
		"NewsLogin",
	)
	png, err := qrcode.Encode(uri, qrcode.Medium, 256)
	if err != nil {
		return TotpInfo{}, err
	}
	qr := "data:image/png;base64," + base64.StdEncoding.EncodeToString([]byte(png))
	return TotpInfo{
		Secret: u.TotpSecret,
		Uri:    uri,
		Qr:     qr,
	}, nil
}

func (u *User) GenerateJWT(scopeTotp bool) (string, error) {
	privateKeyFile, err := os.ReadFile("env/private_key.pem")
	if err != nil {
		log.Printf("Error loading private key file:%+v", err)
		return "", err
	}
	privateKeyBlock, _ := pem.Decode(privateKeyFile)
	privateKey, err := x509.ParsePKCS1PrivateKey(privateKeyBlock.Bytes)
	if err != nil {
		log.Printf("Error parsing private key:%+v", err)
		return "", err
	}
	key := jose.SigningKey{Algorithm: jose.RS256, Key: privateKey}
	var signerOpts = jose.SignerOptions{}
	signerOpts.WithType("JWT")
	rsaSigner, err := jose.NewSigner(key, &signerOpts)
	if err != nil {
		log.Printf("failed to create signer:%+v", err)
		return "", err
	}
	builder := jwt.Signed(rsaSigner)
	now := time.Now().UTC()
	expireAt := time.Now().Add(time.Hour).UTC()
	scopes := map[string]interface{}{
		"sub":         u.UserUuid,
		"scope":       "user",
		"email":       u.Email,
		"totp_active": u.TotpActive,
	}
	if scopeTotp {
		scopes = map[string]interface{}{
			"scope": "verify-totp",
		}
		expireAt = time.Now().Add(time.Minute * 5).UTC()
	}
	claims := map[string]interface{}{
		"typ": "bearer",
		"iat": now.Unix(),
		"nbf": now.Unix(),
		"iss": "login-news-api",
		"exp": expireAt.Unix(),
		"sub": u.UserUuid,
	}
	builder = builder.Claims(claims).Claims(scopes)
	// validate all ok, sign with the RSA key, and return a compact JWT
	jwtString, err := builder.CompactSerialize()
	if err != nil {
		log.Printf("failed to create JWT:%+v", err)
		return "", err
	}
	log.Printf("Token generated succesfully %+v", jwtString)
	return jwtString, nil
}
