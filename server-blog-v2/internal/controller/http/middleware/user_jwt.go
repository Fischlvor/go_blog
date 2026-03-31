package middleware

import (
	"context"
	"crypto/rsa"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"

	"server-blog-v2/internal/controller/http/bizcode"
	"server-blog-v2/internal/controller/http/shared"
	"server-blog-v2/internal/repo/webapi"
	"server-blog-v2/pkg/logger"
)

// claimsKey Context 中存储 Claims 的 key。
type claimsKey struct{}

// AccessClaims JWT Claims。
// 注意：SSO 服务的 token 中不包含 role_id，需要从数据库查询
type AccessClaims struct {
	UUID   string `json:"user_uuid"`
	RoleID uint   `json:"role_id"` // 从数据库填充，token 中可能没有
	jwt.RegisteredClaims
}

// UserIDString 返回用户 UUID。
func (c *AccessClaims) UserIDString() string {
	return c.UUID
}

// WithAccessClaims 将 Claims 存入 Context。
func WithAccessClaims(ctx context.Context, claims *AccessClaims) context.Context {
	return context.WithValue(ctx, claimsKey{}, claims)
}

// AccessClaimsFromContext 从 Context 获取 Claims。
func AccessClaimsFromContext(ctx context.Context) (*AccessClaims, bool) {
	claims, ok := ctx.Value(claimsKey{}).(*AccessClaims)
	return claims, ok
}

// JWTConfig JWT 中间件配置。
type JWTConfig struct {
	PublicKey *rsa.PublicKey
}

// NewUserJWTMiddleware 创建基础 JWT 认证中间件（必须登录）。
//
// Deprecated: 当前路由已统一切换到 `NewSSOUserJWTMiddleware`，以支持 refresh_token 自动刷新。
// 该中间件暂保留仅用于兼容/回滚。
//
//nolint:unused // 保留用于兼容/回滚
func NewUserJWTMiddleware(publicKey *rsa.PublicKey) fiber.Handler {
	return func(c fiber.Ctx) error {
		claims, err := parseToken(c, publicKey)
		if err != nil {
			if errors.Is(err, jwt.ErrTokenExpired) {
				return shared.WriteError(c, http.StatusUnauthorized, bizcode.ErrorTokenExpired, "token expired")
			}
			return shared.WriteError(c, http.StatusUnauthorized, bizcode.ErrorTokenInvalid, "invalid token")
		}

		c.Locals("claims", claims)
		c.SetContext(WithAccessClaims(c.Context(), claims))
		return c.Next()
	}
}

// NewOptionalUserJWTMiddleware 创建可选 JWT 认证中间件（不强制登录，仅本地解析）。
//
// Deprecated: 请使用 `NewOptionalSSOUserJWTMiddleware`，以支持 refresh_token 自动刷新、查询角色和用户同步。
//
//nolint:unused // 保留用于兼容/回滚
func NewOptionalUserJWTMiddleware(publicKey *rsa.PublicKey) fiber.Handler {
	return func(c fiber.Ctx) error {
		claims, err := parseToken(c, publicKey)
		if err == nil && claims != nil {
			c.Locals("claims", claims)
			c.SetContext(WithAccessClaims(c.Context(), claims))
		}
		return c.Next()
	}
}

// NewAdminJWTMiddleware 创建基础管理员 JWT 认证中间件。
//
// Deprecated: 当前 Admin 路由已切换到 `NewSSOAdminJWTMiddleware`，以支持 refresh_token 自动刷新。
// 该中间件暂保留仅用于兼容/回滚。
//
//nolint:unused // 保留用于兼容/回滚
func NewAdminJWTMiddleware(publicKey *rsa.PublicKey, userRoleGetter ...UserRoleGetter) fiber.Handler {
	return func(c fiber.Ctx) error {
		claims, err := parseToken(c, publicKey)
		if err != nil {
			return shared.WriteError(c, http.StatusUnauthorized, bizcode.ErrorTokenInvalid, "invalid token")
		}

		// 如果 token 中没有 role_id，从数据库查询
		if claims.RoleID == 0 && len(userRoleGetter) > 0 && userRoleGetter[0] != nil {
			roleID, err := userRoleGetter[0].GetRoleByUUID(c.Context(), claims.UUID)
			if err == nil {
				claims.RoleID = uint(roleID)
			}
		}

		// RoleID = 2 为管理员
		if claims.RoleID != 2 {
			return shared.WriteError(c, http.StatusForbidden, bizcode.ErrorPermissionDenied, "permission denied")
		}

		c.Locals("claims", claims)
		c.SetContext(WithAccessClaims(c.Context(), claims))
		return c.Next()
	}
}

// UserRoleGetter 用户角色查询接口。
type UserRoleGetter interface {
	GetRoleByUUID(ctx context.Context, uuid string) (int, error)
}

// UserCreator 用户创建接口（用于从 SSO 同步创建用户）。
type UserCreator interface {
	CreateFromSSO(ctx context.Context, uuid, nickname, email, avatar, address, signature string, registerSource int) error
}

// RefreshTokenGetter 获取 refresh_token 的接口。
type RefreshTokenGetter interface {
	GetRefreshToken(c fiber.Ctx) string
	SetRefreshToken(c fiber.Ctx, token string)
	ClearRefreshToken(c fiber.Ctx)
}

// SSOJWTConfig SSO JWT 中间件配置。
type SSOJWTConfig struct {
	PublicKey          *rsa.PublicKey
	SSOClient          *webapi.SSOClient
	UserRoleGetter     UserRoleGetter
	UserCreator        UserCreator
	RefreshTokenGetter RefreshTokenGetter
	Logger             logger.Interface
}

// parseToken 解析 JWT Token。
func parseToken(c fiber.Ctx, publicKey *rsa.PublicKey) (*AccessClaims, error) {
	authorization := c.Get("Authorization")
	if authorization == "" {
		return nil, errors.New("missing authorization header")
	}

	parts := strings.SplitN(authorization, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		return nil, errors.New("invalid authorization format")
	}
	tokenStr := parts[1]

	token, err := jwt.ParseWithClaims(tokenStr, &AccessClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return publicKey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*AccessClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

// GetUserUUID 从 Context 获取用户 UUID。
func GetUserUUID(c fiber.Ctx) string {
	claims, ok := AccessClaimsFromContext(c.Context())
	if !ok || claims == nil {
		return ""
	}
	return claims.UUID
}

// GetOptionalUserUUID 从 Context 获取可选的用户 UUID，未登录时返回 nil。
func GetOptionalUserUUID(c fiber.Ctx) *string {
	uuid := GetUserUUID(c)
	if uuid == "" {
		return nil
	}
	return &uuid
}

// GetRoleID 从 Context 获取角色 ID。
func GetRoleID(c fiber.Ctx) uint {
	claims, ok := AccessClaimsFromContext(c.Context())
	if !ok || claims == nil {
		return 0
	}
	return claims.RoleID
}

// IsAdmin 检查是否为管理员。
func IsAdmin(c fiber.Ctx) bool {
	return GetRoleID(c) == 2
}

// NewSSOJWTMiddleware 创建 SSO JWT 认证中间件（支持自动刷新 token 和用户同步）。
func NewSSOUserJWTMiddleware(cfg SSOJWTConfig) fiber.Handler {
	return func(c fiber.Ctx) error {
		claims, err := parseToken(c, cfg.PublicKey)
		if err != nil {
			if errors.Is(err, jwt.ErrTokenExpired) {
				// Token 过期，尝试自动刷新
				newClaims, refreshErr := autoRefreshToken(c, cfg)
				if refreshErr != nil {
					if cfg.Logger != nil {
						cfg.Logger.Error(refreshErr, "middleware - sso jwt - auto refresh failed")
					}
					return shared.WriteError(c, http.StatusUnauthorized, bizcode.ErrorTokenExpired, "token expired, please login again")
				}
				claims = newClaims
			} else {
				return shared.WriteError(c, http.StatusUnauthorized, bizcode.ErrorTokenInvalid, "invalid token")
			}
		}

		// 从数据库查询用户角色
		if cfg.UserRoleGetter != nil {
			roleID, err := cfg.UserRoleGetter.GetRoleByUUID(c.Context(), claims.UUID)
			if err != nil {
				// 用户不存在，尝试从 SSO 同步创建
				if cfg.UserCreator != nil && cfg.SSOClient != nil {
					if createErr := createUserFromSSO(c.Context(), claims.UUID, cfg); createErr != nil {
						if cfg.Logger != nil {
							cfg.Logger.Error(createErr, "middleware - sso jwt - create user failed")
						}
						return shared.WriteError(c, http.StatusUnauthorized, bizcode.ErrorUnauthorized, "user not found")
					}
					// 重新查询角色
					roleID, _ = cfg.UserRoleGetter.GetRoleByUUID(c.Context(), claims.UUID)
				}
			}
			claims.RoleID = uint(roleID)
		}

		c.Locals("claims", claims)
		c.SetContext(WithAccessClaims(c.Context(), claims))
		return c.Next()
	}
}

// NewOptionalSSOUserJWTMiddleware 创建可选 SSO JWT 认证中间件（不强制登录）。
// 有 token 时走完整 SSO 流程（自动刷新、查角色、同步用户）；无 token 或解析失败时直接放行。
func NewOptionalSSOUserJWTMiddleware(cfg SSOJWTConfig) fiber.Handler {
	return func(c fiber.Ctx) error {
		claims, err := parseToken(c, cfg.PublicKey)
		if err != nil {
			if errors.Is(err, jwt.ErrTokenExpired) {
				// token 过期，尝试自动刷新；刷新失败则当作未登录放行
				newClaims, refreshErr := autoRefreshToken(c, cfg)
				if refreshErr != nil {
					if cfg.Logger != nil {
						cfg.Logger.Debug("middleware - optional sso jwt - refresh failed, continue as guest", "error", refreshErr.Error())
					}
					return c.Next()
				}
				claims = newClaims
			} else {
				// 无 token 或 token 无效，直接放行
				return c.Next()
			}
		}

		// 从数据库查询用户角色
		if cfg.UserRoleGetter != nil {
			roleID, err := cfg.UserRoleGetter.GetRoleByUUID(c.Context(), claims.UUID)
			if err != nil {
				// 用户不存在，尝试从 SSO 同步创建
				if cfg.UserCreator != nil && cfg.SSOClient != nil {
					if createErr := createUserFromSSO(c.Context(), claims.UUID, cfg); createErr != nil {
						if cfg.Logger != nil {
							cfg.Logger.Error(createErr, "middleware - optional sso jwt - create user failed")
						}
						// 同步失败，当作未登录放行
						return c.Next()
					}
					roleID, _ = cfg.UserRoleGetter.GetRoleByUUID(c.Context(), claims.UUID)
				}
			}
			claims.RoleID = uint(roleID)
		}

		c.Locals("claims", claims)
		c.SetContext(WithAccessClaims(c.Context(), claims))
		return c.Next()
	}
}

// NewSSOAdminJWTMiddleware 创建 SSO 管理员 JWT 认证中间件。
func NewSSOAdminJWTMiddleware(cfg SSOJWTConfig) fiber.Handler {
	return func(c fiber.Ctx) error {
		claims, err := parseToken(c, cfg.PublicKey)
		if err != nil {
			if errors.Is(err, jwt.ErrTokenExpired) {
				newClaims, refreshErr := autoRefreshToken(c, cfg)
				if refreshErr != nil {
					return shared.WriteError(c, http.StatusUnauthorized, bizcode.ErrorTokenExpired, "token expired, please login again")
				}
				claims = newClaims
			} else {
				return shared.WriteError(c, http.StatusUnauthorized, bizcode.ErrorTokenInvalid, "invalid token")
			}
		}

		// 从数据库查询用户角色
		if cfg.UserRoleGetter != nil {
			roleID, err := cfg.UserRoleGetter.GetRoleByUUID(c.Context(), claims.UUID)
			if err != nil {
				// 用户不存在，尝试从 SSO 同步创建
				if cfg.UserCreator != nil && cfg.SSOClient != nil {
					if createErr := createUserFromSSO(c.Context(), claims.UUID, cfg); createErr != nil {
						return shared.WriteError(c, http.StatusUnauthorized, bizcode.ErrorUnauthorized, "user not found")
					}
					roleID, _ = cfg.UserRoleGetter.GetRoleByUUID(c.Context(), claims.UUID)
				}
			}
			claims.RoleID = uint(roleID)
		}

		// 检查管理员权限
		if claims.RoleID != 2 {
			return shared.WriteError(c, http.StatusForbidden, bizcode.ErrorPermissionDenied, "permission denied")
		}

		c.Locals("claims", claims)
		c.SetContext(WithAccessClaims(c.Context(), claims))
		return c.Next()
	}
}

// autoRefreshToken 自动刷新 Token。
func autoRefreshToken(c fiber.Ctx, cfg SSOJWTConfig) (*AccessClaims, error) {
	if cfg.RefreshTokenGetter == nil || cfg.SSOClient == nil {
		return nil, errors.New("refresh token not supported")
	}

	refreshToken := cfg.RefreshTokenGetter.GetRefreshToken(c)
	if refreshToken == "" {
		return nil, errors.New("refresh token not found")
	}

	// 向 SSO 刷新 token
	tokenResp, err := cfg.SSOClient.RefreshAccessToken(refreshToken)
	if err != nil {
		cfg.RefreshTokenGetter.ClearRefreshToken(c)
		return nil, err
	}

	// 在响应头返回新 token
	c.Set("X-New-Access-Token", tokenResp.AccessToken)
	c.Set("X-Token-Expires-In", strconv.Itoa(tokenResp.ExpiresIn))

	// 更新 refresh_token（如果 SSO 返回了新的）
	if tokenResp.RefreshToken != "" && tokenResp.RefreshToken != refreshToken {
		cfg.RefreshTokenGetter.SetRefreshToken(c, tokenResp.RefreshToken)
	}

	// 解析新 token
	newClaims, err := parseTokenString(tokenResp.AccessToken, cfg.PublicKey)
	if err != nil {
		return nil, err
	}

	if cfg.Logger != nil {
		cfg.Logger.Info("auto refresh token success", "user_uuid", newClaims.UUID)
	}

	return newClaims, nil
}

// createUserFromSSO 从 SSO 同步创建用户。
func createUserFromSSO(ctx context.Context, userUUID string, cfg SSOJWTConfig) error {
	userInfo, err := cfg.SSOClient.GetUserInfo(userUUID)
	if err != nil {
		return err
	}

	return cfg.UserCreator.CreateFromSSO(
		ctx,
		userUUID,
		userInfo.Nickname,
		userInfo.Email,
		userInfo.Avatar,
		userInfo.Address,
		userInfo.Signature,
		userInfo.RegisterSource,
	)
}

// parseTokenString 解析 JWT Token 字符串。
func parseTokenString(tokenStr string, publicKey *rsa.PublicKey) (*AccessClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &AccessClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return publicKey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*AccessClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}
