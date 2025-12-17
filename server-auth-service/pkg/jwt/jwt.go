package jwt

import (
	"crypto/rsa"
	"errors"
	"time"

	"github.com/gofrs/uuid"
	"github.com/golang-jwt/jwt/v4"
)

var (
	ErrTokenExpired = errors.New("token已过期")
	ErrTokenInvalid = errors.New("token无效")
)

// AccessTokenClaims AccessToken载荷
// 注意：AppID 存储的是 app_key（字符串，如 "mcp"、"blog"），不是数据库中的数字 ID
// 查询数据库时需要先通过 GetAppByKey(AppID) 获取 app.ID
type AccessTokenClaims struct {
	UserUUID  uuid.UUID `json:"user_uuid"`
	AppID     string    `json:"app_id"` // app_key（字符串），不是数字 ID
	DeviceID  string    `json:"device_id"`
	TokenType string    `json:"token_type"` // access_token
	jwt.RegisteredClaims
}

// RefreshTokenClaims RefreshToken载荷
// 注意：AppID 存储的是 app_key（字符串），不是数据库中的数字 ID
type RefreshTokenClaims struct {
	UserUUID  uuid.UUID `json:"user_uuid"`
	AppID     string    `json:"app_id"` // app_key（字符串），不是数字 ID
	DeviceID  string    `json:"device_id"`
	TokenType string    `json:"token_type"` // refresh_token
	jwt.RegisteredClaims
}

// CreateAccessToken 创建AccessToken
func CreateAccessToken(
	userUUID uuid.UUID,
	appID string,
	deviceID string,
	expiry time.Duration,
	issuer string,
	privateKey *rsa.PrivateKey,
) (string, error) {
	now := time.Now()
	claims := AccessTokenClaims{
		UserUUID:  userUUID,
		AppID:     appID,
		DeviceID:  deviceID,
		TokenType: "access_token",
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    issuer,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(expiry)),
			NotBefore: jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	return token.SignedString(privateKey)
}

// CreateRefreshToken 创建RefreshToken
func CreateRefreshToken(
	userUUID uuid.UUID,
	appID string,
	deviceID string,
	expiry time.Duration,
	issuer string,
	privateKey *rsa.PrivateKey,
) (string, error) {
	now := time.Now()
	claims := RefreshTokenClaims{
		UserUUID:  userUUID,
		AppID:     appID,
		DeviceID:  deviceID,
		TokenType: "refresh_token",
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    issuer,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(expiry)),
			NotBefore: jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	return token.SignedString(privateKey)
}

// ParseAccessToken 解析AccessToken
func ParseAccessToken(tokenString string, publicKey *rsa.PublicKey) (*AccessTokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &AccessTokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, ErrTokenInvalid
		}
		return publicKey, nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrTokenExpired
		}
		return nil, ErrTokenInvalid
	}

	if claims, ok := token.Claims.(*AccessTokenClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, ErrTokenInvalid
}

// ParseAccessTokenIgnoreExpiry 解析AccessToken（忽略过期）
func ParseAccessTokenIgnoreExpiry(tokenString string, publicKey *rsa.PublicKey) (*AccessTokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &AccessTokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, ErrTokenInvalid
		}
		return publicKey, nil
	}, jwt.WithoutClaimsValidation())

	if err != nil {
		return nil, ErrTokenInvalid
	}

	if claims, ok := token.Claims.(*AccessTokenClaims); ok {
		return claims, nil
	}

	return nil, ErrTokenInvalid
}

// ParseRefreshToken 解析RefreshToken
func ParseRefreshToken(tokenString string, publicKey *rsa.PublicKey) (*RefreshTokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &RefreshTokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, ErrTokenInvalid
		}
		return publicKey, nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrTokenExpired
		}
		return nil, ErrTokenInvalid
	}

	if claims, ok := token.Claims.(*RefreshTokenClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, ErrTokenInvalid
}
