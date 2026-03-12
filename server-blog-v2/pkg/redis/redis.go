// Package redis go-redis 连接封装。
package redis

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	_defaultReadTimeout  = 3 * time.Second
	_defaultWriteTimeout = 3 * time.Second
	_defaultConnAttempts = 10
	_defaultConnTimeout  = time.Second
)

// Redis go-redis 连接容器。
type Redis struct {
	readTimeout  time.Duration
	writeTimeout time.Duration

	connAttempts int
	connTimeout  time.Duration

	RDB *redis.Client
}

// New 创建 Redis 连接。
func New(host string, port int, password string, db int, opts ...Option) (*Redis, error) {
	r := &Redis{
		readTimeout:  _defaultReadTimeout,
		writeTimeout: _defaultWriteTimeout,
		connAttempts: _defaultConnAttempts,
		connTimeout:  _defaultConnTimeout,
	}

	for _, opt := range opts {
		opt(r)
	}

	r.RDB = redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%d", host, port),
		Password:     password,
		DB:           db,
		WriteTimeout: r.writeTimeout,
		ReadTimeout:  r.readTimeout,
	})

	var err error
	attempts := r.connAttempts
	for attempts > 0 {
		ctx, cancel := context.WithTimeout(context.Background(), r.connTimeout)
		err = r.RDB.Ping(ctx).Err()
		cancel()
		if err == nil {
			break
		}

		log.Printf("Redis: connect retry, attempts left: %d", attempts)
		time.Sleep(r.connTimeout)
		attempts--
	}

	if err != nil {
		return nil, fmt.Errorf("redis - New - connAttempts == 0: %w", err)
	}

	return r, nil
}

// Close 关闭 Redis 连接。
func (r *Redis) Close() {
	if r.RDB != nil {
		_ = r.RDB.Close()
		fmt.Println("Redis: connection closed")
	}
}

// Client Redis 客户端接口。
type Client interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value string, expiration time.Duration) error
}

// Get 获取缓存值。
func (r *Redis) Get(ctx context.Context, key string) (string, error) {
	return r.RDB.Get(ctx, key).Result()
}

// Set 设置缓存值。
func (r *Redis) Set(ctx context.Context, key string, value string, expiration time.Duration) error {
	return r.RDB.Set(ctx, key, value, expiration).Err()
}
