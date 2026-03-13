package middleware

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"time"

	"github.com/gofiber/fiber/v3"

	"server-blog-v2/pkg/logger"
	"server-blog-v2/pkg/redis"
)

const (
	// SessionCookieName Session Cookie 名称。
	SessionCookieName = "blog_session_id"
	// SessionCookieMaxAge Session Cookie 最大存活时间（7天）。
	SessionCookieMaxAge = 7 * 24 * time.Hour
)

// SessionManager Session 管理器，实现 RefreshTokenGetter 接口。
type SessionManager struct {
	redis  *redis.Redis
	logger logger.Interface
}

// NewSessionManager 创建 Session 管理器。
func NewSessionManager(redis *redis.Redis, logger logger.Interface) *SessionManager {
	return &SessionManager{
		redis:  redis,
		logger: logger,
	}
}

// GetRefreshToken 从 Session 获取 refresh_token。
func (s *SessionManager) GetRefreshToken(c fiber.Ctx) string {
	sessionID := c.Cookies(SessionCookieName)
	if sessionID == "" {
		return ""
	}

	refreshToken, err := s.redis.GetRefreshToken(context.Background(), sessionID)
	if err != nil {
		if s.logger != nil {
			s.logger.Debug("session - GetRefreshToken - redis error", "error", err.Error())
		}
		return ""
	}

	return refreshToken
}

// SetRefreshToken 存储 refresh_token 到 Session。
func (s *SessionManager) SetRefreshToken(c fiber.Ctx, token string) {
	sessionID := c.Cookies(SessionCookieName)

	// 如果没有 Session ID，创建新的
	if sessionID == "" {
		sessionID = generateSessionID()
	}

	// 始终设置/刷新 Cookie（确保覆盖旧的，统一 Path 和其他属性）
	c.Cookie(&fiber.Cookie{
		Name:     SessionCookieName,
		Value:    sessionID,
		MaxAge:   int(SessionCookieMaxAge.Seconds()),
		HTTPOnly: true,
		Secure:   false, // 开发环境设为 false，生产环境应设为 true
		SameSite: "Lax",
		Path:     "/",
	})

	// 存储到 Redis
	err := s.redis.SetRefreshToken(context.Background(), sessionID, token, redis.DefaultSessionExpiration)
	if err != nil {
		if s.logger != nil {
			s.logger.Error(err, "session - SetRefreshToken - redis error")
		}
	}
}

// ClearRefreshToken 清除 Session 中的 refresh_token。
func (s *SessionManager) ClearRefreshToken(c fiber.Ctx) {
	sessionID := c.Cookies(SessionCookieName)
	if sessionID == "" {
		return
	}

	// 删除 Redis 中的 Session
	err := s.redis.DeleteSession(context.Background(), sessionID)
	if err != nil {
		if s.logger != nil {
			s.logger.Error(err, "session - ClearRefreshToken - redis error")
		}
	}

	// 清除 Cookie
	c.Cookie(&fiber.Cookie{
		Name:     SessionCookieName,
		Value:    "",
		MaxAge:   -1,
		HTTPOnly: true,
		Secure:   false,
		SameSite: "Lax",
		Path:     "/",
	})
}

// generateSessionID 生成随机 Session ID。
func generateSessionID() string {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		// fallback: 使用时间戳
		return hex.EncodeToString([]byte(time.Now().String()))
	}
	return hex.EncodeToString(bytes)
}
