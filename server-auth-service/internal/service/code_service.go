package service

import (
	"auth-service/internal/model/appTypes"
	"auth-service/pkg/global"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"
)

// GenerateAuthorizationCodeByUUID 生成授权码（使用UUID）
func GenerateAuthorizationCodeByUUID(userUUID, appID, redirectURI, accessToken, refreshToken string) (string, error) {
	// 生成随机授权码
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	code := hex.EncodeToString(bytes)

	// 存储到Redis，有效期5分钟
	key := fmt.Sprintf("auth_code:%s", code)

	// ✅ 使用JSON序列化存储（支持特殊字符）
	authCode := &appTypes.AuthorizationCode{
		Code:         code,
		UserUUID:     userUUID,
		AppID:        appID,
		RedirectURI:  redirectURI,
		ExpiresAt:    time.Now().Add(5 * time.Minute),
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	value, err := json.Marshal(authCode)
	if err != nil {
		return "", err
	}

	err = global.Redis.Set(key, value, 5*time.Minute).Err()
	if err != nil {
		return "", err
	}

	return code, nil
}

// ValidateAndConsumeCode 验证并消费授权码（一次性使用）
func ValidateAndConsumeCode(code, appID, redirectURI string) (*appTypes.AuthorizationCode, error) {
	key := fmt.Sprintf("auth_code:%s", code)

	// 获取授权码信息
	value, err := global.Redis.Get(key).Result()
	if err != nil {
		return nil, fmt.Errorf("授权码无效或已过期")
	}

	// ✅ 使用JSON反序列化
	var authCode appTypes.AuthorizationCode
	if err := json.Unmarshal([]byte(value), &authCode); err != nil {
		return nil, fmt.Errorf("授权码数据格式错误: %v", err)
	}

	// 验证appID和redirectURI
	if authCode.AppID != appID {
		return nil, fmt.Errorf("应用ID不匹配")
	}
	if authCode.RedirectURI != redirectURI {
		return nil, fmt.Errorf("回调地址不匹配")
	}

	// 删除授权码（一次性使用）
	global.Redis.Del(key)

	return &authCode, nil
}
