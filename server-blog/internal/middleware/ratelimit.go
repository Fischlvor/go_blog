package middleware

import (
	"server/internal/model/request"

	ratelimiter "github.com/Fischlvor/go-ratelimiter"
	ginmiddleware "github.com/Fischlvor/go-ratelimiter/drivers/middleware/gin"
	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
)

// RateLimitKeyGetter 自定义 KeyGetter
// 优先使用用户UUID，其次使用设备ID，最后fallback到IP
// 注意：此中间件在私有路由组中位于 SSOJWTAuth 之后执行，可从 context 获取 claims
func RateLimitKeyGetter(c *gin.Context) (path, method, ip, userID string) {
	path = c.Request.URL.Path
	method = c.Request.Method
	ip = c.ClientIP()

	// 1. 优先从 context 获取用户UUID（SSO中间件设置的claims）
	if claims, exists := c.Get("claims"); exists {
		if jwtClaims, ok := claims.(*request.JwtCustomClaims); ok && jwtClaims.UUID != uuid.Nil {
			userID = jwtClaims.UUID.String()
			return
		}
	}

	// 2. 其次使用设备ID（未登录用户）
	userID = c.GetHeader("X-Device-Id")

	// 3. 都没有则userID为空，会自动按IP限流
	return
}

// RateLimitMiddleware 返回限流中间件
func RateLimitMiddleware(limiter *ratelimiter.Limiter) gin.HandlerFunc {
	return ginmiddleware.NewMiddleware(
		limiter,
		ginmiddleware.WithKeyGetter(RateLimitKeyGetter),
	)
}
