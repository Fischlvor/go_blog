package redis

import "time"

// Option Redis 选项。
type Option func(*Redis)

// WithReadTimeout 配置读超时。
func WithReadTimeout(timeout time.Duration) Option {
	return func(r *Redis) {
		r.readTimeout = timeout
	}
}

// WithWriteTimeout 配置写超时。
func WithWriteTimeout(timeout time.Duration) Option {
	return func(r *Redis) {
		r.writeTimeout = timeout
	}
}

// WithConnAttempts 配置连接重试次数。
func WithConnAttempts(attempts int) Option {
	return func(r *Redis) {
		r.connAttempts = attempts
	}
}

// WithConnTimeout 配置连接重试间隔。
func WithConnTimeout(timeout time.Duration) Option {
	return func(r *Redis) {
		r.connTimeout = timeout
	}
}
