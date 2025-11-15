package router

import (
	"github.com/gin-gonic/gin"

	"server/internal/api/admin"
)

type EmojiRouter struct {
}

func (r *EmojiRouter) InitEmojiRouter(Router *gin.RouterGroup) {
	emojiApi := admin.NewEmojiApi()

	emojiRouter := Router.Group("emoji")
	{
		// 表情列表和操作
		emojiRouter.GET("/list", emojiApi.GetEmojiList)        // 获取表情列表
		emojiRouter.POST("/upload", emojiApi.UploadEmoji)      // 批量上传表情
		emojiRouter.DELETE("/:id", emojiApi.DeleteEmoji)       // 删除表情
		emojiRouter.PUT("/:id/restore", emojiApi.RestoreEmoji) // 恢复表情

		// 表情组管理
		emojiRouter.GET("/groups", emojiApi.GetEmojiGroups)          // 获取表情组列表
		emojiRouter.POST("/groups", emojiApi.CreateEmojiGroup)       // 创建表情组
		emojiRouter.PUT("/groups/:id", emojiApi.UpdateEmojiGroup)    // 更新表情组
		emojiRouter.DELETE("/groups/:id", emojiApi.DeleteEmojiGroup) // 删除表情组

		// 雪碧图管理
		emojiRouter.GET("/sprites", emojiApi.GetSpriteList)         // 获取雪碧图列表
		emojiRouter.POST("/regenerate", emojiApi.RegenerateSprites) // 重新生成雪碧图
		emojiRouter.GET("/task/:id", emojiApi.GetTaskStatus)        // 获取任务状态
	}
}
