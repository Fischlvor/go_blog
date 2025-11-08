package middleware

import (
	"auth-service/pkg/global"
	"auth-service/pkg/jwt"
	"auth-service/pkg/utils"
	"strings"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware JWT认证中间件
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			utils.Unauthorized(c, "未提供认证token")
			c.Abort()
			return
		}

		// 检查Bearer前缀
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			utils.Unauthorized(c, "token格式错误")
			c.Abort()
			return
		}

		token := parts[1]

		// 检查黑名单
		blacklistKey := "token:blacklist:" + token
		if global.Redis.Exists(blacklistKey).Val() > 0 {
			utils.Unauthorized(c, "token已失效")
			c.Abort()
			return
		}

		// 解析token
		claims, err := jwt.ParseAccessToken(token, global.RSAPublicKey)
		if err != nil {
			if err == jwt.ErrTokenExpired {
				utils.Unauthorized(c, "token已过期")
			} else {
				utils.Unauthorized(c, "token无效")
			}
			c.Abort()
			return
		}

		// 检查设备黑名单
		deviceBlacklistKey := "device:blacklist:" + claims.DeviceID
		if global.Redis.Exists(deviceBlacklistKey).Val() > 0 {
			utils.Unauthorized(c, "设备已被移除")
			c.Abort()
			return
		}

		// 将用户信息存入context
		c.Set("user_id", claims.UserID)
		c.Set("user_uuid", claims.UserUUID)
		c.Set("app_id", claims.AppID)
		c.Set("device_id", claims.DeviceID)

		c.Next()
	}
}

// GetUserID 从context获取用户ID
func GetUserID(c *gin.Context) uint {
	if userID, exists := c.Get("user_id"); exists {
		return userID.(uint)
	}
	return 0
}

// GetDeviceID 从context获取设备ID
func GetDeviceID(c *gin.Context) string {
	if deviceID, exists := c.Get("device_id"); exists {
		return deviceID.(string)
	}
	return ""
}

// GetAppID 从context获取应用ID
func GetAppID(c *gin.Context) string {
	if appID, exists := c.Get("app_id"); exists {
		return appID.(string)
	}
	return ""
}
