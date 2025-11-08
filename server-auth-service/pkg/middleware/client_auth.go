package middleware

import (
	"auth-service/internal/model/entity"
	"auth-service/pkg/global"
	"auth-service/pkg/utils"

	"github.com/gin-gonic/gin"
)

// ClientAuthMiddleware 客户端认证中间件（用于服务间调用）
func ClientAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientID := c.GetHeader("X-Client-ID")
		clientSecret := c.GetHeader("X-Client-Secret")

		if clientID == "" || clientSecret == "" {
			utils.Unauthorized(c, "缺少客户端认证信息")
			c.Abort()
			return
		}

		// 验证客户端ID和Secret
		var app entity.SSOApplication
		err := global.DB.Where("app_key = ? AND status = 1", clientID).First(&app).Error
		if err != nil {
			utils.Unauthorized(c, "无效的客户端ID")
			c.Abort()
			return
		}

		// 验证Secret
		if app.AppSecret != clientSecret {
			utils.Unauthorized(c, "客户端认证失败")
			c.Abort()
			return
		}

		// 认证成功，将应用信息存入context
		c.Set("app_id", app.ID)
		c.Set("app_key", app.AppKey)
		c.Next()
	}
}
