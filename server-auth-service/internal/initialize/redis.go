package initialize

import (
	"auth-service/pkg/global"
	"fmt"
	"log"

	"github.com/go-redis/redis"
)

// InitRedis 初始化Redis
func InitRedis() *redis.Client {
	conf := global.Config.Redis

	client := redis.NewClient(&redis.Options{
		Addr:     conf.Address,
		Password: conf.Password,
		DB:       conf.DB,
	})

	// 测试连接
	if err := client.Ping().Err(); err != nil {
		panic(fmt.Sprintf("连接Redis失败: %v", err))
	}

	log.Println("Redis连接成功")
	return client
}
