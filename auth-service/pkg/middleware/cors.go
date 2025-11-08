package middleware

import (
	"auth-service/pkg/global"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// Cors 跨域中间件
func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		origin := c.Request.Header.Get("Origin")

		// 判断是否允许该来源
		allowOrigin := "*"
		if !global.Config.System.AllowAllOrigins {
			// 生产环境：检查白名单
			allowed := false
			for _, allowedOrigin := range global.Config.System.AllowedOrigins {
				if origin == allowedOrigin || strings.HasPrefix(origin, allowedOrigin) {
					allowed = true
					allowOrigin = origin
					break
				}
			}
			if !allowed {
				c.AbortWithStatus(http.StatusForbidden)
				return
			}
		} else {
			// 开发环境：允许所有来源
			if origin != "" {
				allowOrigin = origin
			}
		}

		// 设置跨域响应头
		c.Header("Access-Control-Allow-Origin", allowOrigin)
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With, X-App-ID, X-Device-ID")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Content-Type, X-Access-Token, X-Refresh-Token, X-New-Access-Token, X-Token-Expires-In")
		c.Header("Access-Control-Max-Age", "86400")

		// 放行OPTIONS请求
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}
