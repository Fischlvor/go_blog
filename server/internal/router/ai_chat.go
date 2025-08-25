package router

import (
	"server/internal/api"
	"server/internal/middleware"

	"github.com/gin-gonic/gin"
)

type AIChatRouter struct {
}

func (a *AIChatRouter) InitAIChatRouter(Router *gin.RouterGroup, PublicRouter *gin.RouterGroup, AdminRouter *gin.RouterGroup) {
	aiChatRouter := Router.Group("ai-chat").Use(middleware.JWTAuth())
	aiChatPublicRouter := PublicRouter.Group("ai-chat")
	// aiChatAdminRouter := AdminRouter.Group("ai-chat") // 暂时注释，后续需要时再启用
	aiChatApi := api.ApiGroupApp.AIChatApi

	{
		// 需要登录的接口
		aiChatRouter.POST("session", aiChatApi.CreateSession)            // 创建会话
		aiChatRouter.GET("sessions", aiChatApi.GetSessions)              // 获取会话列表
		aiChatRouter.GET("session/:id", aiChatApi.GetSessionDetail)      // 获取会话详情
		aiChatRouter.GET("messages", aiChatApi.GetMessages)              // 获取消息列表
		aiChatRouter.POST("message", aiChatApi.SendMessage)              // 发送消息
		aiChatRouter.POST("message/stream", aiChatApi.SendMessageStream) // 流式发送消息
		aiChatRouter.DELETE("session", aiChatApi.DeleteSession)          // 删除会话
		aiChatRouter.PUT("session", aiChatApi.UpdateSession)             // 更新会话
	}

	{
		// 公开接口
		aiChatPublicRouter.GET("models", aiChatApi.GetAvailableModels) // 获取可用模型列表
	}

	{
		// 管理员接口（如果需要管理AI模型配置）
		// aiChatAdminRouter.POST("model", aiChatApi.CreateModel)
		// aiChatAdminRouter.PUT("model", aiChatApi.UpdateModel)
		// aiChatAdminRouter.DELETE("model", aiChatApi.DeleteModel)
	}
}
