package middleware

import (
	"server/pkg/utils"

	ratelimiter "github.com/Fischlvor/go-ratelimiter"
	ginmiddleware "github.com/Fischlvor/go-ratelimiter/drivers/middleware/gin"
	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
)

// RateLimitKeyGetter 自定义 KeyGetter
// 优先使用用户UUID，其次使用设备ID，最后fallback到IP
func RateLimitKeyGetter(c *gin.Context) (path, method, ip, userID string) {
	path = c.Request.URL.Path
	method = c.Request.Method
	ip = c.ClientIP()

	// 1. 优先使用用户UUID（已登录用户）
	userUUID := utils.GetUUID(c)
	if userUUID != uuid.Nil {
		userID = userUUID.String()
		return
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
