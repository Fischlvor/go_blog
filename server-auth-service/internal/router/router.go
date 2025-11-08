package router

import (
	"auth-service/internal/handler"
	"auth-service/pkg/middleware"

	"github.com/gin-gonic/gin"
)

// Setup 设置路由
func Setup(r *gin.Engine) {
	// 全局中间件
	r.Use(middleware.Cors())

	// API组
	api := r.Group("/api")
	{
		// 基础接口
		base := api.Group("/base")
		{
			captchaHandler := handler.NewCaptchaHandler()
			base.GET("/captcha", captchaHandler.GetCaptcha)
		}

		// 认证相关（无需鉴权）
		auth := api.Group("/auth")
		{
			authHandler := handler.NewAuthHandler()
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.POST("/token", authHandler.RefreshToken) // OAuth 2.0标准token端点
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
		}
	}

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})
}
