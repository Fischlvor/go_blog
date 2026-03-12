package webapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// SSOClient SSO 服务客户端。
type SSOClient struct {
	serviceURL   string
	clientID     string
	clientSecret string
}

// NewSSOClient 创建 SSO 客户端。
func NewSSOClient(serviceURL, clientID, clientSecret string) *SSOClient {
	return &SSOClient{
		serviceURL:   serviceURL,
		clientID:     clientID,
		clientSecret: clientSecret,
	}
}

// TokenResponse SSO Token 响应。
type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
	TokenType    string `json:"token_type"`
}

// SSOUserInfo SSO 用户信息。
type SSOUserInfo struct {
	Nickname       string `json:"nickname"`
	Email          string `json:"email"`
	Avatar         string `json:"avatar"`
	Address        string `json:"address"`
	Signature      string `json:"signature"`
	RegisterSource int    `json:"register_source"` // 0:email 1:qq 2:wechat 3:github
}

// RefreshAccessToken 使用 refresh_token 刷新 access_token。
func (c *SSOClient) RefreshAccessToken(refreshToken string) (*TokenResponse, error) {
	url := fmt.Sprintf("%s/api/auth/token", c.serviceURL)

	requestBody := map[string]string{
		"grant_type":    "refresh_token",
		"refresh_token": refreshToken,
		"client_id":     c.clientID,
		"client_secret": c.clientSecret,
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return nil, err
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("SSO返回错误: %s", string(body))
	}

	var result struct {
		Code    int            `json:"code"`
		Message string         `json:"message"`
		Data    *TokenResponse `json:"data"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	if result.Code != 0 {
		return nil, fmt.Errorf("SSO错误: %s", result.Message)
	}

	return result.Data, nil
}

// GetUserInfo 从 SSO 获取用户信息。
func (c *SSOClient) GetUserInfo(userUUID string) (*SSOUserInfo, error) {
	url := fmt.Sprintf("%s/api/internal/user/%s", c.serviceURL, userUUID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("X-Client-ID", c.clientID)
	req.Header.Set("X-Client-Secret", c.clientSecret)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("SSO返回错误状态: %d, body: %s", resp.StatusCode, string(body))
	}

	var apiResp struct {
		Code    int          `json:"code"`
		Message string       `json:"msg"`
		Data    *SSOUserInfo `json:"data"`
	}

	if err := json.Unmarshal(body, &apiResp); err != nil {
		return nil, err
	}

	if apiResp.Code != 0 {
		return nil, fmt.Errorf("SSO接口返回错误: %s", apiResp.Message)
	}

	return apiResp.Data, nil
}
