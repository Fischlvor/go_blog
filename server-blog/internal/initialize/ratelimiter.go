package initialize

import (
	"server/pkg/global"

	ratelimiter "github.com/Fischlvor/go-ratelimiter"
	redisstore "github.com/Fischlvor/go-ratelimiter/drivers/store/redis"
	"go.uber.org/zap"
)

// InitRateLimiter 初始化限流器
func InitRateLimiter() *ratelimiter.Limiter {
	// 复用全局 Redis 客户端
	store := redisstore.NewStore(&global.Redis, "ratelimit")

	// 从配置文件加载限流器
	limiter, err := ratelimiter.NewFromFile("./configs/rate_limit.yaml", store)
	if err != nil {
		global.Log.Fatal("初始化限流器失败", zap.Error(err))
	}

	global.Log.Info("✓ 限流器初始化成功")
	return limiter
}
