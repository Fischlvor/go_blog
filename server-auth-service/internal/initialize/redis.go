package initialize

import (
	"auth-service/pkg/global"
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis"
)

// InitRedis 初始化Redis
func InitRedis() *redis.Client {
	conf := global.Config.Redis

	client := redis.NewClient(&redis.Options{
		Addr:     conf.Address,
		Password: conf.Password,
		DB:       conf.DB,

		// 连接池配置
		PoolSize:     10, // 连接池大小
		MinIdleConns: 2,  // 最小空闲连接数

		// 超时配置
		DialTimeout:  5 * time.Second, // 连接超时
		ReadTimeout:  3 * time.Second, // 读取超时
		WriteTimeout: 3 * time.Second, // 写入超时
		PoolTimeout:  4 * time.Second, // 连接池超时

		// 空闲连接检查
		IdleTimeout:        5 * time.Minute, // 空闲连接超时
		IdleCheckFrequency: 1 * time.Minute, // 空闲检查频率

		// 重试配置
		MaxRetries:      3, // 最大重试次数
		MinRetryBackoff: 8 * time.Millisecond,
		MaxRetryBackoff: 512 * time.Millisecond,
	})

	// 测试连接
	if err := client.Ping().Err(); err != nil {
		panic(fmt.Sprintf("连接Redis失败: %v", err))
	}

	log.Println("Redis连接成功")
	return client
}
