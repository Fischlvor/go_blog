package router

import (
	"github.com/gin-gonic/gin"

	"server/internal/api/admin"
)

type PublicEmojiRouter struct {
}

func (r *PublicEmojiRouter) InitPublicEmojiRouter(Router *gin.RouterGroup) {
	emojiApi := admin.NewEmojiApi()

	emojiRouter := Router.Group("emoji")
	{
		// 前端公共接口
		emojiRouter.GET("/config", emojiApi.GetEmojiConfig) // 获取前端配置
	}
}
