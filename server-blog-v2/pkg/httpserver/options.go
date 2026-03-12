package httpserver

import (
	"net"
	"time"
)

// Option Server 选项。
type Option func(*Server)

// WithPort 配置监听端口。
func WithPort(port string) Option {
	return func(s *Server) {
		s.address = net.JoinHostPort("", port)
	}
}

// WithPrefork 配置 Prefork 模式。
func WithPrefork(prefork bool) Option {
	return func(s *Server) {
		s.prefork = prefork
	}
}

// WithReadTimeout 配置读超时。
func WithReadTimeout(timeout time.Duration) Option {
	return func(s *Server) {
		s.readTimeout = timeout
	}
}

// WithWriteTimeout 配置写超时。
func WithWriteTimeout(timeout time.Duration) Option {
	return func(s *Server) {
		s.writeTimeout = timeout
	}
}

// WithShutdownTimeout 配置优雅关闭超时。
func WithShutdownTimeout(timeout time.Duration) Option {
	return func(s *Server) {
		s.shutdownTimeout = timeout
	}
}
