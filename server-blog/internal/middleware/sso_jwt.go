package middleware

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"server/internal/api"
	"server/internal/model/appTypes"
	"server/internal/model/database"
	"server/internal/model/request"
	"server/internal/model/response"
	"server/pkg/global"
	"server/pkg/utils"
	"strconv"
	"strings"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"go.uber.org/zap"
)

// SSOJWTAuth SSO JWT认证中间件（使用RSA公钥验证）
func SSOJWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.NoAuth("未提供认证token", c)
			c.Abort()
			return
		}

		// 检查Bearer前缀
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			response.NoAuth("token格式错误", c)
			c.Abort()
			return
		}

		token := parts[1]

		// 使用RSA公钥解析SSO颁发的AccessToken
		claims, err := utils.ParseSSOAccessToken(token)
		if err != nil {
			if errors.Is(err, utils.TokenExpired) {
				// ✅ Token过期，后端自动刷新！
				newToken, refreshErr := autoRefreshToken(c)
				if refreshErr != nil {
					response.NoAuth("token已过期且刷新失败，请重新登录", c)
					c.Abort()
					return
				}

				// ✅ 刷新成功，在响应头返回新token
				c.Header("X-New-Access-Token", newToken.AccessToken)
				c.Header("X-Token-Expires-In", strconv.Itoa(newToken.ExpiresIn))

				// 重新解析新token
				claims, err = utils.ParseSSOAccessToken(newToken.AccessToken)
				if err != nil {
					response.NoAuth("新token解析失败", c)
					c.Abort()
					return
				}
			} else {
				response.NoAuth("token无效: "+err.Error(), c)
				c.Abort()
				return
			}
		}

		// 检查应用ID是否匹配
		if claims.AppID != "blog" {
			response.NoAuth("token不适用于此应用", c)
			c.Abort()
			return
		}

		// ✅ 从博客数据库查询用户信息（根据UUID）
		var user database.User
		err = global.DB.Where("uuid = ?", claims.UserUUID).First(&user).Error
		if err != nil {
			// 用户不存在，从SSO同步创建
			global.Log.Info("用户在blog数据库不存在，自动创建", zap.String("uuid", claims.UserUUID.String()))
			user, err = createUserFromSSO(claims.UserUUID)
			if err != nil {
				response.NoAuth("创建用户失败: "+err.Error(), c)
				c.Abort()
				return
			}
		}

		// ✅ 将SSO claims转换为应用的JwtCustomClaims格式
		jwtClaims := &request.JwtCustomClaims{
			BaseClaims: request.BaseClaims{
				UserID: user.ID,         // 博客的用户ID
				UUID:   claims.UserUUID, // SSO的UUID
				RoleID: user.RoleID,     // 从博客数据库获取真实角色
			},
		}
		c.Set("claims", jwtClaims)
		c.Set("sso_claims", claims)

		c.Next()
	}
}

// autoRefreshToken 自动刷新Token（从session获取refresh_token）
func autoRefreshToken(c *gin.Context) (*api.TokenResponse, error) {
	// 从session获取refresh_token
	session := sessions.Default(c)
	refreshToken := session.Get("refresh_token")
	if refreshToken == nil {
		return nil, errors.New("未找到refresh_token，请重新登录")
	}

	refreshTokenStr, ok := refreshToken.(string)
	if !ok {
		return nil, errors.New("refresh_token格式错误")
	}

	// 向SSO刷新token
	tokenResp, err := api.RefreshAccessTokenFromSSO(refreshTokenStr)
	if err != nil {
		global.Log.Error("自动刷新token失败", zap.Error(err))
		// 刷新失败，清除session
		session.Delete("refresh_token")
		session.Delete("refresh_token_expires_at")
		session.Save()
		return nil, err
	}

	// 更新session中的refresh_token（如果SSO返回了新的）
	if tokenResp.RefreshToken != "" && tokenResp.RefreshToken != refreshTokenStr {
		session.Set("refresh_token", tokenResp.RefreshToken)
		session.Save()
	}

	global.Log.Info("✓ 自动刷新token成功")
	return tokenResp, nil
}

// createUserFromSSO 从SSO同步创建用户
func createUserFromSSO(userUUID uuid.UUID) (database.User, error) {
	// 调用SSO获取用户信息
	ssoUserInfo, err := fetchSSOUserInfo(userUUID)
	if err != nil {
		return database.User{}, fmt.Errorf("获取SSO用户信息失败: %w", err)
	}

	// 映射注册来源：SSO(0:email 1:qq 2:wechat 3:github) -> Blog(0:Email 1:QQ)
	var registerType appTypes.Register
	switch ssoUserInfo.RegisterSource {
	case 0: // email
		registerType = appTypes.Email
	case 1: // qq
		registerType = appTypes.QQ
	default:
		registerType = appTypes.Email // 其他来源默认为Email
	}

	// 创建博客用户
	user := database.User{
		UUID:      userUUID,
		Username:  ssoUserInfo.Nickname, // 使用昵称作为用户名
		Email:     ssoUserInfo.Email,
		Avatar:    ssoUserInfo.Avatar, // 使用默认头像或SSO头像
		Address:   ssoUserInfo.Address,
		Signature: ssoUserInfo.Signature,
		RoleID:    appTypes.User, // 默认普通用户角色
		Register:  registerType,  // 从SSO映射注册来源
		Freeze:    false,
	}

	// 保存到数据库
	if err := global.DB.Create(&user).Error; err != nil {
		return database.User{}, fmt.Errorf("创建用户失败: %w", err)
	}

	global.Log.Info("✓ 从SSO同步创建用户成功",
		zap.String("uuid", userUUID.String()),
		zap.String("nickname", user.Username))

	return user, nil
}

// SSOUserInfoResponse SSO用户信息响应
type SSOUserInfoResponse struct {
	Nickname       string `json:"nickname"`
	Email          string `json:"email"`
	Avatar         string `json:"avatar"`
	Address        string `json:"address"`
	Signature      string `json:"signature"`
	RegisterSource int    `json:"register_source"` // 0:email 1:qq 2:wechat 3:github
}

// fetchSSOUserInfo 从SSO获取用户信息
func fetchSSOUserInfo(userUUID uuid.UUID) (*SSOUserInfoResponse, error) {
	// 构建请求URL（使用内部接口）
	url := fmt.Sprintf("%s/api/internal/user/%s", global.Config.SSO.ServiceURL, userUUID.String())

	// 创建HTTP请求
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	// 添加客户端认证
	req.Header.Set("X-Client-ID", global.Config.SSO.ClientID)
	req.Header.Set("X-Client-Secret", global.Config.SSO.ClientSecret)
	req.Header.Set("Content-Type", "application/json")

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("SSO返回错误状态: %d, body: %s", resp.StatusCode, string(body))
	}

	// 解析响应
	var apiResp struct {
		Code    int                  `json:"code"`
		Message string               `json:"msg"`
		Data    *SSOUserInfoResponse `json:"data"`
	}

	if err := json.Unmarshal(body, &apiResp); err != nil {
		return nil, err
	}

	if apiResp.Code != 0 {
		return nil, fmt.Errorf("SSO接口返回错误: %s", apiResp.Message)
	}

	return apiResp.Data, nil
}
