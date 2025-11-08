package router

import (
	"server/internal/api"

	"github.com/gin-gonic/gin"
)

type AIManagementRouter struct{}

func (s *AIManagementRouter) InitAIManagementRouter(AdminRouter *gin.RouterGroup) {
	aiManagementAdminRouter := AdminRouter.Group("ai-management")
	aiManagementApi := api.ApiGroupApp.AIManagementApi

	// AI模型管理 - 仅管理员
	{
		aiManagementAdminRouter.POST("model", aiManagementApi.CreateAIModel)       // 创建AI模型
		aiManagementAdminRouter.DELETE("model/:id", aiManagementApi.DeleteAIModel) // 删除AI模型
		aiManagementAdminRouter.PUT("model", aiManagementApi.UpdateAIModel)        // 更新AI模型
		aiManagementAdminRouter.GET("model/:id", aiManagementApi.GetAIModel)       // 获取AI模型详情
		aiManagementAdminRouter.GET("model", aiManagementApi.GetAIModelList)       // 获取AI模型列表
	}

	// AI会话管理 - 仅管理员
	{
		aiManagementAdminRouter.GET("session", aiManagementApi.GetAISessionList)       // 获取AI会话列表
		aiManagementAdminRouter.GET("session/:id", aiManagementApi.GetAISession)       // 获取AI会话详情
		aiManagementAdminRouter.DELETE("session/:id", aiManagementApi.DeleteAISession) // 删除AI会话
	}

	// AI消息管理 - 仅管理员
	{
		aiManagementAdminRouter.GET("message", aiManagementApi.GetAIMessageList)       // 获取AI消息列表
		aiManagementAdminRouter.GET("message/:id", aiManagementApi.GetAIMessage)       // 获取AI消息详情
		aiManagementAdminRouter.DELETE("message/:id", aiManagementApi.DeleteAIMessage) // 删除AI消息
	}
}
