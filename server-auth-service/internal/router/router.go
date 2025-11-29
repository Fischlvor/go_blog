package router

import (
	"auth-service/internal/handler"
	"auth-service/pkg/global"
	"auth-service/pkg/middleware"
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
	api := r.Group("/api")
	{
		// 基础接口
		base := api.Group("/base")
		{
			captchaHandler := handler.NewCaptchaHandler()
			base.GET("/captcha", captchaHandler.GetCaptcha)
			authHandler := handler.NewAuthHandler()
			base.GET("/qqLoginURL", authHandler.QQLoginURL) // QQ登录URL
		}

		// 认证相关（无需鉴权）
		auth := api.Group("/auth")
		{
			authHandler := handler.NewAuthHandler()
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.POST("/token", authHandler.RefreshToken)                                  // OAuth 2.0标准token端点
			auth.POST("/oauth/qq/login", authHandler.QQLogin)                              // QQ登录
			auth.GET("/oauth/qq/callback", authHandler.QQCallback)                         // QQ回调（GET：QQ服务端回调）
			auth.POST("/sendEmailVerificationCode", authHandler.SendEmailVerificationCode) // 发送邮箱验证码
			auth.POST("/forgotPassword", authHandler.ForgotPassword)                       // 忘记密码
		}

		// ✅ OAuth 2.0 授权端点（SSO 静默登录）
		oauth := api.Group("/oauth")
		{
			oauthHandler := handler.NewOAuthHandler()
			oauth.GET("/authorize", oauthHandler.Authorize) // 授权端点（检查 Session）
		}

		// 服务间调用接口（使用客户端认证）
		internal := api.Group("/internal", middleware.ClientAuthMiddleware())
		{
			authHandler := handler.NewAuthHandler()
			internal.GET("/user/:uuid", authHandler.GetUserByUUID)
		}

		// 需要鉴权的接口
		authenticated := api.Group("", middleware.AuthMiddleware())
		{
			// 用户信息
			user := authenticated.Group("/user")
			{
				authHandler := handler.NewAuthHandler()
				user.GET("/info", authHandler.GetUserInfo)
				user.PUT("/info", authHandler.UpdateUserInfo)
				user.PUT("/password", authHandler.UpdatePassword)
				user.POST("/logout", authHandler.Logout)
			}

			// 设备管理
			device := authenticated.Group("/devices")
			{
				deviceHandler := handler.NewDeviceHandler()
				device.GET("", deviceHandler.GetDevices)
				device.DELETE("/:device_id", deviceHandler.KickDevice)
			}

			// SSO管理后台
			manage := authenticated.Group("/manage")
			{
				manageHandler := handler.NewManageHandler()
				manage.GET("/devices", manageHandler.GetDevices)           // 获取设备列表
				manage.POST("/kick-device", manageHandler.KickDevice)      // 踢出设备
				manage.POST("/sso-logout", manageHandler.SSOLogout)        // SSO退出
				manage.POST("/logout-all", manageHandler.LogoutAllDevices) // 退出所有设备
				manage.GET("/logs", manageHandler.GetLogs)                 // 操作日志
				manage.GET("/profile", manageHandler.GetProfile)           // 用户信息
			}
		}
	}

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})
}
