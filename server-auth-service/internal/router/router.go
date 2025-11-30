package router

import (
	"auth-service/internal/middleware"
	"auth-service/pkg/global"
	"os"
	"strconv"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
)

// Setup 设置路由
func Setup(r *gin.Engine) {
	// 全局中间件
	r.Use(middleware.Cors())

	// ✅ SSO Session 中间件（全局登录状态）
	// 使用 NewStoreWithDB 支持指定数据库
	store, err := redis.NewStoreWithDB(
		10,
		"tcp",
		global.Config.Redis.Address,
		"",                                    // username（Redis 6.0+，留空表示默认用户）
		global.Config.Redis.Password,          // password
		strconv.Itoa(global.Config.Redis.DB),  // db
		[]byte("sso_session_secret_key_2025"), // keyPairs
	)

	if err != nil {
		// 如果连接失败，打印详细错误信息
		panic("初始化 Session 存储失败: " + err.Error() + ", 请检查 Redis 配置是否正确")
	}

	sessionOptions := sessions.Options{
		Path:     "/",
		MaxAge:   7 * 24 * 3600, // 7 天
		HttpOnly: true,          // 防止 XSS
		SameSite: 2,             // Lax 模式
	}

	// 根据环境变量配置 Session
	env := os.Getenv("APP_ENV")
	if env == "prod" || env == "production" {
		// 生产环境：跨子域名 SSO
		sessionOptions.Domain = ".hsk423.cn"
		sessionOptions.Secure = true // HTTPS only
	}
	// 开发环境：本地测试
	// Domain 留空（Cookie 只在当前域名有效）
	// Secure = false（支持 HTTP）

	store.Options(sessionOptions)
	r.Use(sessions.Sessions("sso_session", store))

	// API组
	apiGroup := r.Group("/api")
	{
		// 初始化各个路由
		RouterGroupApp.BaseRouter.InitBaseRouter(apiGroup)
		RouterGroupApp.AuthRouter.InitAuthRouter(apiGroup)
		RouterGroupApp.OAuthRouter.InitOAuthRouter(apiGroup)
		RouterGroupApp.UserRouter.InitUserRouter(apiGroup)
		RouterGroupApp.ManageRouter.InitManageRouter(apiGroup)
	}

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})
}
