package v1

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/gofiber/fiber/v3"

	"server-blog-v2/internal/controller/http/bizcode"
	"server-blog-v2/internal/controller/http/shared"
)

// getSSOLoginURL 获取 SSO 登录 URL。
func (v *V1) getSSOLoginURL(c fiber.Ctx) error {
	redirectURI := c.Query("redirect_uri")
	if redirectURI == "" {
		redirectURI = fmt.Sprintf("http://localhost:%d/sso-callback", v.cfg.HTTP.Port)
	}

	returnURL := c.Query("return_url")
	if returnURL == "" {
		returnURL = "/"
	}

	state := fmt.Sprintf(`{"return_url":"%s"}`, returnURL)

	ssoLoginURL := fmt.Sprintf("%s/api/oauth/authorize?app_id=%s&redirect_uri=%s&state=%s",
		v.cfg.SSO.WebURL,
		v.cfg.SSO.ClientID,
		url.QueryEscape(redirectURI),
		url.QueryEscape(state),
	)

	return shared.WriteSuccess(c, shared.WithData(fiber.Map{
		"sso_login_url": ssoLoginURL,
	}))
}

// ssoCallback SSO 回调接口。
func (v *V1) ssoCallback(c fiber.Ctx) error {
	code := c.Query("code")
	state := c.Query("state")
	redirectURI := c.Query("redirect_uri")

	if code == "" {
		return shared.WriteError(c, http.StatusBadRequest, bizcode.ErrorParamMissing, "missing authorization code")
	}

	tokenResp, err := v.exchangeCodeForToken(code, redirectURI)
	if err != nil {
		v.logger.Error(err, "http - v1 - auth - ssoCallback - exchangeCodeForToken")
		return shared.WriteError(c, http.StatusInternalServerError, bizcode.ErrorThirdParty, "failed to exchange token")
	}

	return shared.WriteSuccess(c, shared.WithData(fiber.Map{
		"access_token": tokenResp.AccessToken,
		"token_type":   tokenResp.TokenType,
		"expires_in":   tokenResp.ExpiresIn,
		"user_info":    tokenResp.UserInfo,
		"state":        state,
	}))
}

// TokenResponse SSO token 响应。
type TokenResponse struct {
	AccessToken  string      `json:"access_token"`
	RefreshToken string      `json:"refresh_token"`
	TokenType    string      `json:"token_type"`
	ExpiresIn    int         `json:"expires_in"`
	UserInfo     interface{} `json:"user_info,omitempty"`
}

// exchangeCodeForToken 使用 code 向 SSO 换取 token。
func (v *V1) exchangeCodeForToken(code, redirectURI string) (*TokenResponse, error) {
	ssoURL := fmt.Sprintf("%s/api/auth/token", v.cfg.SSO.ServiceURL)

	requestBody := map[string]string{
		"grant_type":    "authorization_code",
		"code":          code,
		"client_id":     v.cfg.SSO.ClientID,
		"client_secret": v.cfg.SSO.ClientSecret,
		"redirect_uri":  redirectURI,
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return nil, err
	}

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
		return nil, fmt.Errorf("SSO returned error: %s", string(body))
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
		return nil, fmt.Errorf("SSO error: %s", result.Message)
	}

	return result.Data, nil
}
