package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"server/internal/model/response"
	"server/pkg/global"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type AuthApi struct{}

// GetSSOLoginURL 获取SSO登录地址
func (a *AuthApi) GetSSOLoginURL(c *gin.Context) {
	// 生成state用于防CSRF
	state := generateState()

	// 构建SSO登录URL
	redirectURI := c.Query("redirect_uri")
	if redirectURI == "" {
		redirectURI = fmt.Sprintf("http://localhost:%d/callback", global.Config.System.Port) // 默认值使用当前服务端口
	}

	// ✅ 从配置文件读取SSO前端地址
	ssoLoginURL := fmt.Sprintf("%s/login?app_id=%s&redirect_uri=%s&state=%s",
		global.Config.SSO.WebURL,
		global.Config.SSO.ClientID,
		url.QueryEscape(redirectURI),
		state,
	)

	response.OkWithData(gin.H{
		"sso_login_url": ssoLoginURL,
		"state":         state,
	}, c)
}

// SSOCallback SSO回调接口（应用后端接收code，换取token）
func (a *AuthApi) SSOCallback(c *gin.Context) {
	// 获取授权码
	code := c.Query("code")
	state := c.Query("state")
	redirectURI := c.Query("redirect_uri")

	if code == "" {
		response.FailWithMessage("缺少授权码", c)
		return
	}

	// ✅ OAuth 2.0: 使用code向SSO服务换取token
	tokenResp, err := exchangeCodeForToken(code, redirectURI)
	if err != nil {
		global.Log.Error("换取token失败", zap.Error(err))
		response.FailWithMessage("换取token失败: "+err.Error(), c)
		return
	}

	// ✅ 关键：refresh_token存在后端session中，不返回给前端！
	session := sessions.Default(c)
	session.Set("refresh_token", tokenResp.RefreshToken)
	session.Set("refresh_token_expires_at", time.Now().Add(7*24*time.Hour).Unix())
	if err := session.Save(); err != nil {
		global.Log.Error("保存session失败", zap.Error(err))
		response.FailWithMessage("保存会话失败", c)
		return
	}

	// ✅ 返回access_token给前端，state用于前端获取返回URL
	response.OkWithData(gin.H{
		"access_token": tokenResp.AccessToken,
		"token_type":   tokenResp.TokenType,
		"expires_in":   tokenResp.ExpiresIn,
		"user_info":    tokenResp.UserInfo,
		"state":        state,
	}, c)
}

// TokenResponse SSO token响应
type TokenResponse struct {
	AccessToken  string      `json:"access_token"`
	RefreshToken string      `json:"refresh_token"` // 后端持有，不返回给前端
	TokenType    string      `json:"token_type"`
	ExpiresIn    int         `json:"expires_in"`
	UserInfo     interface{} `json:"user_info,omitempty"`
}

// RefreshAccessTokenFromSSO 后端用refresh_token向SSO刷新（内部函数，中间件调用）
func RefreshAccessTokenFromSSO(refreshToken string) (*TokenResponse, error) {
	return refreshAccessToken(refreshToken)
}

// exchangeCodeForToken 使用code向SSO换取token
func exchangeCodeForToken(code, redirectURI string) (*TokenResponse, error) {
	// ✅ 从配置读取SSO服务地址
	ssoURL := fmt.Sprintf("%s/api/auth/token", global.Config.SSO.ServiceURL)

	// ✅ 构建JSON请求体
	requestBody := map[string]string{
		"grant_type":    "authorization_code",
		"code":          code,
		"client_id":     global.Config.SSO.ClientID,
		"client_secret": global.Config.SSO.ClientSecret,
		"redirect_uri":  redirectURI,
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return nil, err
	}

	// ✅ 发送JSON POST请求
	resp, err := http.Post(ssoURL, "application/json", bytes.NewBuffer(jsonData))
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

// refreshAccessToken 刷新AccessToken
func refreshAccessToken(refreshToken string) (*TokenResponse, error) {
	// ✅ 从配置读取SSO服务地址
	ssoURL := fmt.Sprintf("%s/api/auth/token", global.Config.SSO.ServiceURL)

	// ✅ 构建JSON请求体
	requestBody := map[string]string{
		"grant_type":    "refresh_token",
		"refresh_token": refreshToken,
		"client_id":     global.Config.SSO.ClientID,
		"client_secret": global.Config.SSO.ClientSecret,
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return nil, err
	}

	// ✅ 发送JSON POST请求
	resp, err := http.Post(ssoURL, "application/json", bytes.NewBuffer(jsonData))
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

// generateState 生成随机state
func generateState() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}
