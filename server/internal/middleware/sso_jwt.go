package middleware

import (
	"errors"
	"server/internal/api"
	"server/internal/model/database"
	"server/internal/model/request"
	"server/internal/model/response"
	"server/pkg/global"
	"server/pkg/utils"
	"strconv"
	"strings"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
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
		if err := global.DB.Where("uuid = ?", claims.UserUUID).First(&user).Error; err != nil {
			response.NoAuth("用户不存在", c)
			c.Abort()
			return
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
