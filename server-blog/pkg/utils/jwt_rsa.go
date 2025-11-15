package utils

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"os"

	"github.com/gofrs/uuid"
	"github.com/golang-jwt/jwt/v4"
)

// SSOAccessTokenClaims SSO的AccessToken载荷
type SSOAccessTokenClaims struct {
	UserUUID  uuid.UUID `json:"user_uuid"`
	AppID     string    `json:"app_id"`
	DeviceID  string    `json:"device_id"`
	TokenType string    `json:"token_type"` // access_token
	jwt.RegisteredClaims
}

var (
	SSOPublicKey *rsa.PublicKey
)

// LoadSSOPublicKey 加载SSO公钥
func LoadSSOPublicKey(path string) error {
	keyData, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	block, _ := pem.Decode(keyData)
	if block == nil {
		return errors.New("无法解析PEM格式的公钥")
	}

	publicKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return err
	}

	rsaPublicKey, ok := publicKey.(*rsa.PublicKey)
	if !ok {
		return errors.New("不是RSA公钥")
	}

	SSOPublicKey = rsaPublicKey
	return nil
}

// ParseSSOAccessToken 解析SSO的AccessToken（使用RSA公钥验证）
func ParseSSOAccessToken(tokenString string) (*SSOAccessTokenClaims, error) {
	if SSOPublicKey == nil {
		return nil, errors.New("SSO公钥未加载")
	}

	token, err := jwt.ParseWithClaims(tokenString, &SSOAccessTokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, TokenInvalid
		}
		return SSOPublicKey, nil
	})

	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			switch {
			case ve.Errors&jwt.ValidationErrorMalformed != 0:
				return nil, TokenMalformed
			case ve.Errors&jwt.ValidationErrorExpired != 0:
				return nil, TokenExpired
			case ve.Errors&jwt.ValidationErrorNotValidYet != 0:
				return nil, TokenNotValidYet
			default:
				return nil, TokenInvalid
			}
		}
		return nil, TokenInvalid
	}

	if claims, ok := token.Claims.(*SSOAccessTokenClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, TokenInvalid
}
