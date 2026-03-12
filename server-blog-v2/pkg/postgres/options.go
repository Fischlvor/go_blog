package postgres

import "time"

// Option Postgres 选项。
type Option func(*Postgres)

// WithConnAttempts 配置连接重试次数。
func WithConnAttempts(attempts int) Option {
	return func(m *Postgres) {
		m.connAttempts = attempts
	}
}

// WithConnTimeout 配置连接重试间隔。
func WithConnTimeout(timeout time.Duration) Option {
	return func(m *Postgres) {
		m.connTimeout = timeout
	}
}
