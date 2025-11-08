package handler

import (
	"auth-service/internal/model/request"
	"auth-service/internal/service"
	"auth-service/pkg/middleware"
	"auth-service/pkg/utils"

	"github.com/gin-gonic/gin"
	"github.com/mojocn/base64Captcha"
)

type AuthHandler struct {
	authService *service.AuthService
	store       base64Captcha.Store
}

func NewAuthHandler() *AuthHandler {
	return &AuthHandler{
		authService: &service.AuthService{},
		store:       base64Captcha.DefaultMemStore,
	}
}

// Register 注册
func (h *AuthHandler) Register(c *gin.Context) {
	var req request.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	// 验证验证码
	if !h.store.Verify(req.CaptchaID, req.Captcha, true) {
		utils.Error(c, 1000, "验证码错误")
		return
	}

	// ✅ 注册用户（不自动登录）
	if err := h.authService.Register(req); err != nil {
		utils.Error(c, 1001, err.Error())
		return
	}

	// 返回成功，前端跳转到登录页
	utils.SuccessMsg(c, "注册成功，请登录", nil)
}

// Login 登录
func (h *AuthHandler) Login(c *gin.Context) {
	var req request.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	// 验证验证码
	if !h.store.Verify(req.CaptchaID, req.Captcha, true) {
		utils.Error(c, 1000, "验证码错误")
		return
	}

	// 获取客户端信息
	ipAddress := c.ClientIP()
	userAgent := c.GetHeader("User-Agent")

	// 登录验证
	resp, err := h.authService.Login(req, ipAddress, userAgent)
	if err != nil {
		utils.Error(c, 1002, err.Error())
		return
	}

	// ✅ OAuth 2.0: 生成授权码
	code, err := service.GenerateAuthorizationCode(resp.UserInfo.UserID, req.AppID, req.RedirectURI, resp.AccessToken, resp.RefreshToken)
	if err != nil {
		utils.Error(c, 1003, "生成授权码失败")
		return
	}

	utils.Success(c, gin.H{
		"code":         code,
		"redirect_uri": req.RedirectURI,
	})
}

// RefreshToken OAuth 2.0 token端点（支持authorization_code和refresh_token）
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req request.TokenExchangeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	// 验证client_id和client_secret
	if req.ClientID != "blog" || req.ClientSecret != "blog_secret_2025_go_blog_system" {
		utils.Unauthorized(c, "client_id或client_secret错误")
		return
	}

	switch req.GrantType {
	case "authorization_code":
		// ✅ 用授权码换取token
		authCode, err := service.ValidateAndConsumeCode(req.Code, req.ClientID, req.RedirectURI)
		if err != nil {
			utils.Error(c, 1003, err.Error())
			return
		}

		utils.Success(c, gin.H{
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
		resp, err := h.authService.RefreshToken(refreshReq)
		if err != nil {
			utils.Error(c, 1003, err.Error())
			return
		}
		utils.Success(c, resp)

	default:
		utils.BadRequest(c, "不支持的grant_type")
	}
}

// Logout 登出
func (h *AuthHandler) Logout(c *gin.Context) {
	// 从Authorization header获取token
	token := c.GetHeader("Authorization")
	if len(token) > 7 && token[:7] == "Bearer " {
		token = token[7:]
	}

	if err := h.authService.Logout(token); err != nil {
		utils.Error(c, 1004, err.Error())
		return
	}

	utils.SuccessMsg(c, "登出成功", nil)
}

// GetUserInfo 获取用户信息
func (h *AuthHandler) GetUserInfo(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		utils.Unauthorized(c, "未登录")
		return
	}

	userInfo, err := h.authService.GetUserInfo(userID)
	if err != nil {
		utils.Error(c, 1005, err.Error())
		return
	}

	utils.Success(c, userInfo)
}

// UpdateUserInfo 更新用户信息
func (h *AuthHandler) UpdateUserInfo(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		utils.Unauthorized(c, "未登录")
		return
	}

	var req request.UpdateUserInfoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	if err := h.authService.UpdateUserInfo(userID, req); err != nil {
		utils.Error(c, 1006, err.Error())
		return
	}

	utils.SuccessMsg(c, "更新成功", nil)
}

// UpdatePassword 修改密码
func (h *AuthHandler) UpdatePassword(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		utils.Unauthorized(c, "未登录")
		return
	}

	var req request.UpdatePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	if err := h.authService.UpdatePassword(userID, req); err != nil {
		utils.Error(c, 1007, err.Error())
		return
	}

	utils.SuccessMsg(c, "密码修改成功", nil)
}

// GetUserByUUID 根据UUID获取用户信息（服务间调用）
func (h *AuthHandler) GetUserByUUID(c *gin.Context) {
	userUUIDStr := c.Param("uuid")
	if userUUIDStr == "" {
		utils.BadRequest(c, "缺少UUID参数")
		return
	}

	userInfo, err := h.authService.GetUserByUUID(userUUIDStr)
	if err != nil {
		utils.Error(c, 1008, err.Error())
		return
	}

	utils.Success(c, userInfo)
}
