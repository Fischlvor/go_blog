package handler

import (
	"auth-service/internal/service"
	"auth-service/pkg/global"
	"auth-service/pkg/utils"
	"fmt"
	"net/url"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"go.uber.org/zap"
)

type OAuthHandler struct {
	authService *service.AuthService
}

func NewOAuthHandler() *OAuthHandler {
	return &OAuthHandler{
		authService: &service.AuthService{},
	}
}

// Authorize OAuth 2.0 授权端点（检查 Session，实现静默登录）
// GET /api/oauth/authorize?app_id=blog&redirect_uri=xxx&state=xxx
func (h *OAuthHandler) Authorize(c *gin.Context) {
	appID := c.Query("app_id")
	redirectURI := c.Query("redirect_uri")
	state := c.Query("state")

	// 参数验证
	if appID == "" || redirectURI == "" {
		utils.BadRequest(c, "缺少必要参数")
		return
	}

	// ✅ 检查全局 Session
	session := sessions.Default(c)
	userUUID := session.Get("user_uuid")
	ssoDeviceID := session.Get("sso_device_id")
	sessionUserAgent := session.Get("user_agent")
	sessionIPAddress := session.Get("ip_address")
	loggedIn := session.Get("logged_in")

	if userUUID != nil && loggedIn == true {
		// ✅ 用户已登录 SSO，进行安全检测和静默登录
		userUUIDStr, ok := userUUID.(string)
		if !ok {
			utils.Error(c, 1014, "Session 数据格式错误")
			return
		}

		// 获取 SSO 设备 ID
		ssoDeviceIDStr := ""
		if ssoDeviceID != nil {
			ssoDeviceIDStr, _ = ssoDeviceID.(string)
		}
		if ssoDeviceIDStr == "" {
			global.Log.Warn("Session 中缺少 sso_device_id，要求重新登录")
			c.Redirect(302, "/login")
			return
		}

		// 安全检测：检查 User-Agent 和 IP 地址变化
		currentUserAgent := c.GetHeader("User-Agent")
		currentIPAddress := c.ClientIP()

		sessionUAStr, _ := sessionUserAgent.(string)
		sessionIPStr, _ := sessionIPAddress.(string)

		// 检测异常行为
		if sessionUAStr != "" && sessionUAStr != currentUserAgent {
			global.Log.Warn("SSO Session 异常：User-Agent 变化",
				zap.String("user_uuid", userUUIDStr),
				zap.String("session_ua", sessionUAStr),
				zap.String("current_ua", currentUserAgent),
			)
			// 中安全模式：记录日志但允许继续
		}

		if sessionIPStr != "" && sessionIPStr != currentIPAddress {
			global.Log.Warn("SSO Session 异常：IP 地址变化",
				zap.String("user_uuid", userUUIDStr),
				zap.String("session_ip", sessionIPStr),
				zap.String("current_ip", currentIPAddress),
			)
			// 中安全模式：记录日志但允许继续
		}

		// 检查设备过期状态（滑动过期）
		err := h.authService.CheckDeviceExpiry(ssoDeviceIDStr)
		if err != nil {
			global.Log.Warn("设备已过期，清除 Session",
				zap.String("user_uuid", userUUIDStr),
				zap.String("device_id", ssoDeviceIDStr),
				zap.Error(err),
			)

			// 清除 Session
			session.Clear()
			session.Save()

			// 重定向到登录页面
			loginURL := fmt.Sprintf("/login?app_id=%s&redirect_uri=%s", appID, url.QueryEscape(redirectURI))
			if state != "" {
				loginURL += "&state=" + url.QueryEscape(state)
			}
			c.Redirect(302, loginURL)
			return
		}

		global.Log.Info("✓ SSO 静默登录成功",
			zap.String("user_uuid", userUUIDStr),
			zap.String("sso_device_id", ssoDeviceIDStr),
			zap.String("app_id", appID),
		)

		// 生成新的 AccessToken 和 RefreshToken（使用 SSO 设备 ID）
		tokenResp, err := h.authService.GenerateTokensForUser(userUUIDStr, appID, ssoDeviceIDStr)
		if err != nil {
			utils.Error(c, 1015, "生成 Token 失败: "+err.Error())
			return
		}

		// 记录静默登录日志
		userUUIDParsed, _ := uuid.FromString(userUUIDStr)
		// 获取应用 ID
		app, _ := h.authService.GetAppByKey(appID)
		if app != nil {
			h.authService.LogAction(userUUIDParsed, app.ID, "silent_login", ssoDeviceIDStr, "SSO静默登录成功", 1)
		}

		// 生成授权码
		code, err := service.GenerateAuthorizationCodeByUUID(
			userUUIDStr,
			appID,
			redirectURI,
			tokenResp.AccessToken,
			tokenResp.RefreshToken,
		)
		if err != nil {
			utils.Error(c, 1016, "生成授权码失败")
			return
		}

		// 构建回调 URL
		callbackURL := redirectURI + "?code=" + code
		if state != "" {
			callbackURL += "&state=" + state
		}

		// 重定向回应用（带授权码）
		c.Redirect(302, callbackURL)
		return
	}

	// ❌ 用户未登录，重定向到 SSO 登录页面
	// 构建登录页面 URL，保留原始参数以便登录后继续授权流程
	loginURL := fmt.Sprintf("/login?app_id=%s&redirect_uri=%s",
		appID,
		url.QueryEscape(redirectURI),
	)
	if state != "" {
		loginURL += "&state=" + url.QueryEscape(state)
	}

	global.Log.Info("✗ SSO 需要登录，重定向到登录页面",
		zap.String("app_id", appID),
		zap.String("redirect_uri", redirectURI),
	)

	// 重定向到登录页面（相对路径，同域名下）
	c.Redirect(302, loginURL)
}
