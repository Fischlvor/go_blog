// Package httpserver Fiber v3 HTTP Server 封装。
package httpserver

import (
	"context"
	"errors"
	"time"

	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v3"
	"golang.org/x/sync/errgroup"

	"server-blog-v2/pkg/logger"
)

const (
	_defaultAddress         = ":8080"
	_defaultReadTimeout     = 5 * time.Second
	_defaultWriteTimeout    = 5 * time.Second
	_defaultShutdownTimeout = 3 * time.Second
)

// Server HTTP Server 实例。
type Server struct {
	ctx context.Context
	eg  *errgroup.Group

	App    *fiber.App
	notify chan error

	address         string
	prefork         bool
	readTimeout     time.Duration
	writeTimeout    time.Duration
	shutdownTimeout time.Duration

	logger logger.Interface
}

// New 创建 Server。
func New(l logger.Interface, opts ...Option) *Server {
	group, ctx := errgroup.WithContext(context.Background())
	group.SetLimit(1)

	s := &Server{
		ctx:             ctx,
		eg:              group,
		App:             nil,
		notify:          make(chan error, 1),
		address:         _defaultAddress,
		readTimeout:     _defaultReadTimeout,
		writeTimeout:    _defaultWriteTimeout,
		shutdownTimeout: _defaultShutdownTimeout,
		logger:          l,
	}

	for _, opt := range opts {
		opt(s)
	}

	app := fiber.New(fiber.Config{
		ReadTimeout:  s.readTimeout,
		WriteTimeout: s.writeTimeout,
		JSONDecoder:  json.Unmarshal,
		JSONEncoder:  json.Marshal,
		BodyLimit:    50 * 1024 * 1024, // 50MB
	})

	s.App = app

	return s
}

// Start 启动 HTTP Server。
func (s *Server) Start() {
	s.eg.Go(func() error {
		err := s.App.Listen(s.address, fiber.ListenConfig{EnablePrefork: s.prefork})
		if err != nil {
			s.notify <- err

			close(s.notify)

			return err
		}

		return nil
	})

	s.logger.Info("http server - Server - Started")
}

// Notify 返回服务启动与运行错误通道。
func (s *Server) Notify() <-chan error {
	return s.notify
}

// Shutdown 优雅关闭 HTTP Server。
func (s *Server) Shutdown() error {
	var shutdownErrors []error

	err := s.App.ShutdownWithTimeout(s.shutdownTimeout)
	if err != nil && !errors.Is(err, context.Canceled) {
		s.logger.Error(err, "http server - Server - Shutdown - s.App.ShutdownWithTimeout")

		shutdownErrors = append(shutdownErrors, err)
	}

	err = s.eg.Wait()
	if err != nil && !errors.Is(err, context.Canceled) {
		s.logger.Error(err, "http server - Server - Shutdown - s.eg.Wait")

		shutdownErrors = append(shutdownErrors, err)
	}

	s.logger.Info("http server - Server - Shutdown")

	return errors.Join(shutdownErrors...)
}
