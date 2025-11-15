package middleware

import (
	"net/http"
	"path"
	"strings"

	"server/internal/model/appTypes"
	"server/internal/model/response"
	"server/pkg/utils"

	"github.com/gin-gonic/gin"
)

// token白名单：固定令牌 -> 允许放行的接口前缀（按前缀匹配）
var fixedTokenWhitelist = map[string][]string{
	"37395c61-a2ec-464e-9567-ce6fa92630f7": {
		"/api/image/upload",
	},
}

// isWhitelistedToken 校验是否命中固定token白名单
func isWhitelistedToken(authHeader, fullPath, urlPath, method string) bool {
	parts := strings.SplitN(authHeader, " ", 2)
	if !(len(parts) == 2 && parts[0] == "Bearer") {
		return false
	}
	token := parts[1]
	allowedPrefixes, ok := fixedTokenWhitelist[token]
	if !ok {
		return false
	}
	// 仅允许POST方法，避免误放行
	if method != http.MethodPost {
		return false
	}
	requestPath := fullPath
	if requestPath == "" {
		requestPath = path.Clean(urlPath)
	}
	for _, p := range allowedPrefixes {
		if strings.HasPrefix(requestPath, p) {
			return true
		}
	}
	return false
}

// AdminAuth 是一个中间件，用于检查用户是否具有管理员权限
func AdminAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 固定token白名单放行（与SSOJWTAuth保持一致）
		authHeader := c.GetHeader("Authorization")
		if isWhitelistedToken(authHeader, c.FullPath(), c.Request.URL.Path, c.Request.Method) {
			c.Next()
			return
		}
		// 从上下文中获取用户的角色ID
		roleID := utils.GetRoleID(c)

		// 检查用户是否为管理员
		if roleID != appTypes.Admin {
			// 如果不是管理员，返回访问被拒绝的响应
			response.Forbidden("Access denied. Admin privileges are required", c)
			// 中止请求处理
			c.Abort()
			return
		}

		// 如果用户是管理员，继续执行后续处理器
		c.Next()
	}
}
