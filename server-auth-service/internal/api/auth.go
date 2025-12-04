package api

import (
	"auth-service/internal/middleware"
	"auth-service/internal/model/request"
	"auth-service/internal/model/response"
	"auth-service/internal/service"
	"auth-service/pkg/global"
	"auth-service/pkg/utils"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"go.uber.org/zap"
)

type AuthApi struct {
}

// Register 注册
func (h *AuthApi) Register(c *gin.Context) {
	var req request.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	// 验证验证码
	if !captchaStore.Verify(req.CaptchaID, req.Captcha, true) {
		response.Error(c, 1000, "验证码错误")
		return
	}

	// ✅ 注册用户（不自动登录）
	if err := authService.Register(req); err != nil {
		response.Error(c, 1001, err.Error())
		return
	}

	// 返回成功，前端跳转到登录页
	response.SuccessMsg(c, "注册成功，请登录", nil)
}

// Login 登录
func (h *AuthApi) Login(c *gin.Context) {
	var req request.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	// 验证state参数
	stateData, err := utils.ValidateState(req.State)
	if err != nil {
		response.BadRequest(c, "state验证失败: "+err.Error())
		return
	}

	// 验证app_id一致性：state中的app_id必须与请求参数中的app_id一致
	if stateData.AppID != req.AppID {
		response.BadRequest(c, "app_id参数与state中的app_id不一致")
		return
	}

	// 从state中提取其他参数
	req.DeviceID = stateData.DeviceID
	req.RedirectURI = stateData.RedirectURI

	// 判断登录方式：密码登录或验证码登录
	if req.Password != "" {
		// 密码登录：需要图片验证码
		if req.CaptchaID == "" || req.Captcha == "" {
			response.BadRequest(c, "密码登录需要图片验证码")
			return
		}
		// 验证图片验证码
		if !captchaStore.Verify(req.CaptchaID, req.Captcha, true) {
			response.Error(c, 1000, "验证码错误")
			return
		}
	} else if req.VerificationCode != "" {
		// 验证码登录：不需要图片验证码
		// 验证码验证在service层进行
	} else {
		response.BadRequest(c, "请提供密码或邮箱验证码")
		return
	}

	// 登录验证
	resp, err := authService.Login(c, req)
	if err != nil {
		response.Error(c, 1002, err.Error())
		return
	}

	// ✅ SSO: 设置全局 Session Cookie（跨应用单点登录）
	session := sessions.Default(c)
	session.Set("user_uuid", resp.UserInfo.UUID)
	session.Set("sso_device_id", req.DeviceID)           // 存储 SSO 设备 ID
	session.Set("user_agent", c.GetHeader("User-Agent")) // 安全检测
	session.Set("ip_address", c.ClientIP())              // 安全检测
	session.Set("logged_in", true)
	session.Set("logged_in_at", time.Now().Unix())
	if err := session.Save(); err != nil {
		global.Log.Error("保存 Session 失败", zap.Error(err))
	}

	// 检查是否是管理后台登录
	if req.AppID == "manage" {
		// 管理后台登录：直接返回Token，不走OAuth流程
		response.Success(c, gin.H{
			"access_token":  resp.AccessToken,
			"refresh_token": resp.RefreshToken,
			"token_type":    "Bearer",
			"expires_in":    7200, // 2小时
			"redirect_uri":  "http://localhost:3001/manage",
		})
	} else {
		// ✅ OAuth 2.0: 生成授权码（使用UUID）
		code, err := service.GenerateAuthorizationCodeByUUID(resp.UserInfo.UUID, req.AppID, req.RedirectURI, resp.AccessToken, resp.RefreshToken)
		if err != nil {
			response.Error(c, 1003, "生成授权码失败")
			return
		}

		response.Success(c, gin.H{
			"code":         code,
			"redirect_uri": stateData.RedirectURI,
			"return_url":   stateData.ReturnURL,
		})
	}
}

// RefreshToken OAuth 2.0 token端点（支持authorization_code和refresh_token）
func (h *AuthApi) RefreshToken(c *gin.Context) {
	var req request.TokenExchangeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	// 从数据库验证 client_id 和 client_secret
	app, err := authService.GetAppByKey(req.ClientID)
	if err != nil || app.AppSecret != req.ClientSecret {
		response.Unauthorized(c, "client_id或client_secret错误")
		return
	}

	switch req.GrantType {
	case "authorization_code":
		// ✅ 用授权码换取token
		authCode, err := service.ValidateAndConsumeCode(req.Code, req.ClientID, req.RedirectURI)
		if err != nil {
			response.Error(c, 1003, err.Error())
			return
		}

		response.Success(c, gin.H{
			"access_token":  authCode.AccessToken,
			"refresh_token": authCode.RefreshToken,
			"token_type":    "Bearer",
			"expires_in":    7200,
		})

	case "refresh_token":
		// ✅ 用refresh_token刷新
		refreshReq := request.RefreshTokenRequest{
			GrantType:    req.GrantType,
			RefreshToken: req.RefreshToken,
			ClientID:     req.ClientID,
			ClientSecret: req.ClientSecret,
		}
		resp, err := authService.RefreshToken(refreshReq)
		if err != nil {
			response.Error(c, 1003, err.Error())
			return
		}
		response.Success(c, resp)

	default:
		response.BadRequest(c, "不支持的grant_type")
	}
}

// Logout 登出
func (h *AuthApi) Logout(c *gin.Context) {
	// 从Authorization header获取token
	token := c.GetHeader("Authorization")
	if len(token) > 7 && token[:7] == "Bearer " {
		token = token[7:]
	}

	// 获取客户端信息
	ipAddress := c.ClientIP()
	userAgent := c.GetHeader("User-Agent")

	if err := authService.Logout(token, ipAddress, userAgent); err != nil {
		response.Error(c, 1004, err.Error())
		return
	}

	// 注意：应用退出不清除 SSO Session，保持其他应用的静默登录能力
	// 只有 /api/manage/sso-logout 才清除 SSO Session
	// session := sessions.Default(c)
	// session.Clear()
	// session.Save()

	response.SuccessMsg(c, "登出成功", nil)
}

// GetUserInfo 获取用户信息
func (h *AuthApi) GetUserInfo(c *gin.Context) {
	userUUID := middleware.GetUserUUID(c)
	if userUUID == uuid.Nil {
		response.Unauthorized(c, "未登录")
		return
	}

	userInfo, err := authService.GetUserInfoByUUID(userUUID)
	if err != nil {
		response.Error(c, 1005, err.Error())
		return
	}

	response.Success(c, userInfo)
}

// UpdateUserInfo 更新用户信息
func (h *AuthApi) UpdateUserInfo(c *gin.Context) {
	userUUID := middleware.GetUserUUID(c)
	if userUUID == uuid.Nil {
		response.Unauthorized(c, "未登录")
		return
	}

	var req request.UpdateUserInfoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	if err := authService.UpdateUserInfoByUUID(userUUID, req); err != nil {
		response.Error(c, 1006, err.Error())
		return
	}

	response.SuccessMsg(c, "更新成功", nil)
}

// UpdatePassword 修改密码
func (h *AuthApi) UpdatePassword(c *gin.Context) {
	userUUID := middleware.GetUserUUID(c)
	if userUUID == uuid.Nil {
		response.Unauthorized(c, "未登录")
		return
	}

	var req request.UpdatePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	if err := authService.UpdatePasswordByUUID(userUUID, req); err != nil {
		response.Error(c, 1007, err.Error())
		return
	}

	response.SuccessMsg(c, "密码修改成功", nil)
}

// GetUserByUUID 根据UUID获取用户信息（服务间调用）
func (h *AuthApi) GetUserByUUID(c *gin.Context) {
	userUUIDStr := c.Param("uuid")
	if userUUIDStr == "" {
		response.BadRequest(c, "缺少UUID参数")
		return
	}

	userInfo, err := authService.GetUserByUUID(userUUIDStr)
	if err != nil {
		response.Error(c, 1008, err.Error())
		return
	}

	response.Success(c, userInfo)
}

// QQLogin QQ登录
func (h *AuthApi) QQLogin(c *gin.Context) {
	var req request.QQLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	// 获取客户端信息
	ipAddress := c.ClientIP()
	userAgent := c.GetHeader("User-Agent")

	// 调用QQ登录服务
	resp, err := authService.QQLogin(req, ipAddress, userAgent)
	if err != nil {
		response.Error(c, 1009, err.Error())
		return
	}

	// ✅ SSO: 设置全局 Session Cookie（跨应用单点登录）
	session := sessions.Default(c)
	session.Set("user_uuid", resp.UserInfo.UUID)
	session.Set("sso_device_id", req.DeviceID)           // 存储 SSO 设备 ID
	session.Set("user_agent", c.GetHeader("User-Agent")) // 安全检测
	session.Set("ip_address", c.ClientIP())              // 安全检测
	session.Set("logged_in", true)
	session.Set("logged_in_at", time.Now().Unix())
	if err := session.Save(); err != nil {
		global.Log.Error("保存 Session 失败", zap.Error(err))
	}

	// ✅ OAuth 2.0: 生成授权码（使用UUID）
	code, err := service.GenerateAuthorizationCodeByUUID(resp.UserInfo.UUID, req.AppID, req.RedirectURI, resp.AccessToken, resp.RefreshToken)
	if err != nil {
		response.Error(c, 1003, "生成授权码失败")
		return
	}

	response.Success(c, gin.H{
		"code":         code,
		"redirect_uri": req.RedirectURI,
	})
}

// QQCallback QQ授权回调（GET方式，QQ服务端回调）
func (h *AuthApi) QQCallback(c *gin.Context) {
	code := c.Query("code")
	appID := c.Query("app_id")
	state := c.Query("state")
	if code == "" {
		response.BadRequest(c, "缺少code")
		return
	}
	if strings.TrimSpace(appID) == "" {
		response.BadRequest(c, "缺少app_id")
		return
	}
	if state == "" {
		response.BadRequest(c, "缺少state参数")
		return
	}

	// 验证state参数
	stateData, err := utils.ValidateState(state)
	if err != nil {
		response.BadRequest(c, "state验证失败: "+err.Error())
		return
	}

	// 验证app_id是否匹配
	if stateData.AppID != appID {
		response.BadRequest(c, "app_id不匹配")
		return
	}

	// 客户端信息
	ipAddress := c.ClientIP()
	userAgent := c.GetHeader("User-Agent")

	// 从 User-Agent 解析设备名称
	deviceName := parseDeviceNameFromUserAgent(userAgent)

	// 组装请求
	req := request.QQLoginRequest{
		Code:        code,
		AppID:       stateData.AppID,
		RedirectURI: stateData.RedirectURI,
		DeviceID:    stateData.DeviceID,
		DeviceType:  "web",
		DeviceName:  deviceName,
	}

	// 调用服务登录
	resp, err := authService.QQLogin(req, ipAddress, userAgent)
	if err != nil {
		response.Error(c, 1009, err.Error())
		return
	}

	// 生成授权码并重定向到 redirect_uri?code=...
	authCode, err := service.GenerateAuthorizationCodeByUUID(resp.UserInfo.UUID, req.AppID, req.RedirectURI, resp.AccessToken, resp.RefreshToken)
	if err != nil {
		response.Error(c, 1003, "生成授权码失败")
		return
	}
	// 重定向到回调地址，携带code和return_url
	redirectURL := fmt.Sprintf("%s?code=%s", stateData.RedirectURI, authCode)
	if stateData.ReturnURL != "" {
		redirectURL = redirectURL + "&return_url=" + url.QueryEscape(stateData.ReturnURL)
	}
	c.Redirect(302, redirectURL)
}

// SendEmailVerificationCode 发送邮箱验证码
func (h *AuthApi) SendEmailVerificationCode(c *gin.Context) {
	var req request.SendEmailVerificationCodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	// 验证图片验证码
	if !captchaStore.Verify(req.CaptchaID, req.Captcha, true) {
		response.Error(c, 1000, "验证码错误")
		return
	}

	// 发送邮箱验证码
	if err := authService.SendEmailVerificationCode(req.Email); err != nil {
		response.Error(c, 1010, err.Error())
		return
	}

	response.SuccessMsg(c, "验证码已发送", nil)
}

// ForgotPassword 忘记密码
func (h *AuthApi) ForgotPassword(c *gin.Context) {
	var req request.ForgotPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	// 验证邮箱验证码
	if err := authService.ForgotPassword(req); err != nil {
		response.Error(c, 1011, err.Error())
		return
	}

	response.SuccessMsg(c, "密码重置成功", nil)
}

// QQLoginURL 获取QQ登录URL
func (h *AuthApi) QQLoginURL(c *gin.Context) {
	if !global.Config.QQ.Enable {
		response.Error(c, 1012, "QQ登录未启用")
		return
	}
	// 允许前端传入 app_id，拼接到 redirect_uri，确保 QQ 回调时可携带 app_id
	appID := c.Query("app_id")
	if strings.TrimSpace(appID) == "" {
		response.BadRequest(c, "缺少app_id")
		return
	}
	// 支持state参数传递（用于传递设备ID等信息）
	state := c.Query("state")

	redirectURI := global.Config.QQ.RedirectURI
	redirectURI = redirectURI + "?app_id=" + url.QueryEscape(appID)

	authURL := "https://graph.qq.com/oauth2.0/authorize?" +
		"response_type=code&" +
		"client_id=" + global.Config.QQ.AppID + "&" +
		"redirect_uri=" + url.QueryEscape(redirectURI)
	// state参数只需要在授权URL中传递一次，QQ会原样返回
	if state != "" {
		authURL = authURL + "&state=" + url.QueryEscape(state)
	}
	response.Success(c, gin.H{"url": authURL})
}

// parseDeviceNameFromUserAgent 从 User-Agent 解析设备名称
func parseDeviceNameFromUserAgent(userAgent string) string {
	if userAgent == "" {
		return "未知设备"
	}

	ua := strings.ToLower(userAgent)

	// 检测操作系统
	os := "未知系统"
	if strings.Contains(ua, "windows") {
		if strings.Contains(ua, "windows nt 10.0") || strings.Contains(ua, "windows nt 6.3") {
			os = "Windows"
		} else if strings.Contains(ua, "windows nt 6.1") {
			os = "Windows 7"
		} else {
			os = "Windows"
		}
	} else if strings.Contains(ua, "mac os x") || strings.Contains(ua, "macintosh") {
		os = "macOS"
	} else if strings.Contains(ua, "linux") {
		os = "Linux"
	} else if strings.Contains(ua, "android") {
		os = "Android"
	} else if strings.Contains(ua, "iphone") || strings.Contains(ua, "ipad") || strings.Contains(ua, "ipod") {
		os = "iOS"
	}

	// 检测浏览器
	browser := "未知浏览器"
	if strings.Contains(ua, "edg/") {
		browser = "Edge"
	} else if strings.Contains(ua, "chrome/") && !strings.Contains(ua, "edg/") {
		browser = "Chrome"
	} else if strings.Contains(ua, "firefox/") {
		browser = "Firefox"
	} else if strings.Contains(ua, "safari/") && !strings.Contains(ua, "chrome/") {
		browser = "Safari"
	} else if strings.Contains(ua, "opera/") || strings.Contains(ua, "opr/") {
		browser = "Opera"
	}

	return fmt.Sprintf("%s - %s", os, browser)
}
