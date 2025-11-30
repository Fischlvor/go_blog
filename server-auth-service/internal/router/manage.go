package router

import (
	"auth-service/internal/api"
	"auth-service/internal/middleware"

	"github.com/gin-gonic/gin"
)

type ManageRouter struct{}

func (m *ManageRouter) InitManageRouter(apiGroup *gin.RouterGroup) {
	// 需要鉴权的接口
	authenticated := apiGroup.Group("", middleware.AuthMiddleware())
	{
		manage := authenticated.Group("/manage")
		{
			manageApi := api.ApiGroupApp.ManageApi

			// 设备管理
			devices := manage.Group("/devices")
			{
				devices.GET("list", manageApi.GetDevices)
				devices.POST("/kick", manageApi.KickDevice)
				devices.POST("/logout", manageApi.Logout)
				devices.POST("/sso-logout", manageApi.SSOLogout)
				devices.POST("/logout-all", manageApi.LogoutAllDevices)
			}

			// 日志和用户信息
			manage.GET("/logs", manageApi.GetLogs)
			manage.GET("/profile", manageApi.GetProfile)
		}
	}
}
