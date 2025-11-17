package utils

import (
	"auth-service/pkg/global"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"
)

// StateData OAuth state参数结构
type StateData struct {
	Nonce       string `json:"nonce"`        // 防重放攻击的随机值
	AppID       string `json:"app_id"`       // 应用ID
	DeviceID    string `json:"device_id"`    // 设备ID
	RedirectURI string `json:"redirect_uri"` // 回调地址
	ReturnURL   string `json:"return_url"`   // 用户目标页面
	IAT         int64  `json:"iat"`          // 签发时间戳（秒）
	EXP         int64  `json:"exp"`          // 过期时间戳（秒）
}

// EncodeState 编码state参数
func EncodeState(data *StateData) (string, error) {
	// 验证必需字段
	if data.Nonce == "" {
		return "", errors.New("nonce不能为空")
	}
	if data.AppID == "" {
		return "", errors.New("app_id不能为空")
	}
	if data.RedirectURI == "" {
		return "", errors.New("redirect_uri不能为空")
	}

	// 设置时间戳（如果未设置）
	if data.IAT == 0 {
		data.IAT = time.Now().Unix()
	}
	if data.EXP == 0 {
		data.EXP = data.IAT + 300 // 默认5分钟过期
	}

	// JSON序列化
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return "", fmt.Errorf("JSON序列化失败: %w", err)
	}

	// Base64编码
	encoded := base64.StdEncoding.EncodeToString(jsonBytes)
	return encoded, nil
}

// DecodeState 解码state参数
func DecodeState(state string) (*StateData, error) {
	if state == "" {
		return nil, errors.New("state参数为空")
	}

	// Base64解码
	jsonBytes, err := base64.StdEncoding.DecodeString(state)
	if err != nil {
		return nil, fmt.Errorf("Base64解码失败: %w", err)
	}

	// JSON反序列化
	var data StateData
	if err := json.Unmarshal(jsonBytes, &data); err != nil {
		return nil, fmt.Errorf("JSON解析失败: %w", err)
	}

	return &data, nil
}

// ValidateState 验证state参数（包含nonce验证、过期验证、redirect_uri白名单验证）
func ValidateState(state string) (*StateData, error) {
	// 1. 解码state
	data, err := DecodeState(state)
	if err != nil {
		return nil, err
	}

	// 2. 验证必需字段
	if data.Nonce == "" {
		return nil, errors.New("state中缺少nonce")
	}
	if data.AppID == "" {
		return nil, errors.New("state中缺少app_id")
	}
	if data.RedirectURI == "" {
		return nil, errors.New("state中缺少redirect_uri")
	}

	// 3. 验证过期时间
	now := time.Now().Unix()
	if data.EXP < now {
		return nil, errors.New("state已过期")
	}

	// 4. 验证并消费nonce（防重放攻击）
	if err := ValidateAndConsumeNonce(data.Nonce); err != nil {
		return nil, fmt.Errorf("nonce验证失败: %w", err)
	}

	// 5. 验证redirect_uri白名单
	if err := ValidateRedirectURI(data.AppID, data.RedirectURI); err != nil {
		return nil, fmt.Errorf("redirect_uri验证失败: %w", err)
	}

	return data, nil
}

// ValidateAndConsumeNonce 验证并消费nonce（确保只能使用一次）
func ValidateAndConsumeNonce(nonce string) error {
	if nonce == "" {
		return errors.New("nonce为空")
	}

	key := fmt.Sprintf("oauth_nonce:%s", nonce)

	// 检查nonce是否存在
	exists, err := global.Redis.Exists(key).Result()
	if err != nil {
		return fmt.Errorf("Redis查询失败: %w", err)
	}

	if exists > 0 {
		// nonce已被使用
		return errors.New("nonce已被使用（可能是重放攻击）")
	}

	// 标记nonce为已使用（存储5分钟，与state过期时间一致）
	err = global.Redis.Set(key, "used", 5*time.Minute).Err()
	if err != nil {
		return fmt.Errorf("Redis存储失败: %w", err)
	}

	return nil
}

// ValidateRedirectURI 验证redirect_uri是否在应用的白名单中
func ValidateRedirectURI(appID, redirectURI string) error {
	if appID == "" {
		return errors.New("app_id为空")
	}
	if redirectURI == "" {
		return errors.New("redirect_uri为空")
	}

	// 从数据库查询应用的redirect_uris白名单
	var app struct {
		RedirectURIs string
	}
	err := global.DB.Table("sso_applications").
		Select("redirect_uris").
		Where("app_key = ? AND status = 1", appID).
		First(&app).Error

	if err != nil {
		return fmt.Errorf("应用不存在或已禁用: %w", err)
	}

	if app.RedirectURIs == "" {
		return errors.New("应用未配置redirect_uris白名单")
	}

	// 解析白名单（逗号分隔）
	allowedURIs := strings.Split(app.RedirectURIs, ",")
	for i := range allowedURIs {
		allowedURIs[i] = strings.TrimSpace(allowedURIs[i])
	}

	// 解析请求的redirect_uri，只比较base URL（不包含查询参数）
	parsedURI, err := url.Parse(redirectURI)
	if err != nil {
		return fmt.Errorf("redirect_uri格式错误: %w", err)
	}
	// 获取base URL（scheme + host + path）
	baseURI := parsedURI.Scheme + "://" + parsedURI.Host + parsedURI.Path

	// 验证redirect_uri是否在白名单中
	for _, allowed := range allowedURIs {
		// 也解析白名单中的URL，只比较base部分
		parsedAllowed, err := url.Parse(allowed)
		if err != nil {
			continue // 跳过格式错误的白名单项
		}
		baseAllowed := parsedAllowed.Scheme + "://" + parsedAllowed.Host + parsedAllowed.Path

		if baseAllowed == baseURI {
			return nil // 验证通过
		}
	}

	return fmt.Errorf("redirect_uri不在白名单中: %s", baseURI)
}
