package initialize

import (
	"github.com/go-redis/redis"
	"go.uber.org/zap"
	"os"
	"server/global"
)

// ConnectRedis 初始化并返回一个 Redis 客户端，支持集群或单节点配置
func ConnectRedis() redis.Client {
	redisCfg := global.Config.Redis
	cfg := &redis.Options{
		Addr:     redisCfg.Address,  // 设置 Redis 服务器地址
		Password: redisCfg.Password, // 设置 Redis 密码
		DB:       redisCfg.DB,       // 设置使用的数据库索引
	}

	// 使用单节点配置创建 Redis 客户端
	client := redis.NewClient(cfg)

	// Ping Redis 服务器以检查连接是否正常
	_, err := client.Ping().Result()
	if err != nil {
		global.Log.Error("Failed to connect to Redis:", zap.Error(err))
		os.Exit(1)
	}

	return *client
}
