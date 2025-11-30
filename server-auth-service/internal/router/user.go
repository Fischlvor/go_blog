package router

import (
	"auth-service/internal/api"
	"auth-service/internal/middleware"

	"github.com/gin-gonic/gin"
)

type UserRouter struct{}

func (u *UserRouter) InitUserRouter(apiGroup *gin.RouterGroup) {
	// 服务间调用接口（使用客户端认证）
	internal := apiGroup.Group("/internal", middleware.ClientAuthMiddleware())
	{
		internalUser := internal.Group("/user")
		{
			authApi := api.ApiGroupApp.AuthApi
			internalUser.GET("/:uuid", authApi.GetUserByUUID)
		}
	}

	// 需要鉴权的接口
	authenticated := apiGroup.Group("", middleware.AuthMiddleware())
	{
		user := authenticated.Group("/user")
		{
			authApi := api.ApiGroupApp.AuthApi
			user.GET("/info", authApi.GetUserInfo)
			user.POST("/updateInfo", authApi.UpdateUserInfo)
			user.POST("/updatePassword", authApi.UpdatePassword)
			user.POST("/logout", authApi.Logout)
		}
	}
}
