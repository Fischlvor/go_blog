package handler

import (
	"auth-service/internal/service"
	"auth-service/pkg/global"
	"auth-service/pkg/utils"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
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
	loggedIn := session.Get("logged_in")

	if userUUID != nil && loggedIn == true {
		// ✅ 用户已登录 SSO，直接生成授权码（静默登录）
		userUUIDStr, ok := userUUID.(string)
		if !ok {
			utils.Error(c, 1014, "Session 数据格式错误")
			return
		}

		global.Log.Info("✓ SSO 静默登录",
			zap.String("user_uuid", userUUIDStr),
			zap.String("app_id", appID),
		)

		// 生成新的 AccessToken 和 RefreshToken
		tokenResp, err := h.authService.GenerateTokensForUser(userUUIDStr, appID)
		if err != nil {
			utils.Error(c, 1015, "生成 Token 失败: "+err.Error())
			return
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

	// ❌ 用户未登录，返回需要登录的提示
	// 前端会根据这个响应跳转到登录页面
	utils.Success(c, gin.H{
		"need_login":   true,
		"app_id":       appID,
		"redirect_uri": redirectURI,
		"state":        state,
		"message":      "需要登录",
	})
}
