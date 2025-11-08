package initialize

import (
	"server/pkg/global"
	"server/internal/model/database"
)

// InitAITables 初始化AI相关表
func InitAITables() {
	// 创建AI模型配置表
	err := global.DB.AutoMigrate(&database.AIModel{})
	if err != nil {
		global.Log.Error("创建AI模型表失败: " + err.Error())
		return
	}

	// 创建AI聊天会话表
	err = global.DB.AutoMigrate(&database.AIChatSession{})
	if err != nil {
		global.Log.Error("创建AI聊天会话表失败: " + err.Error())
		return
	}

	// 创建AI聊天消息表
	err = global.DB.AutoMigrate(&database.AIChatMessage{})
	if err != nil {
		global.Log.Error("创建AI聊天消息表失败: " + err.Error())
		return
	}

	global.Log.Info("AI相关表初始化完成")
}

// InitDefaultAIModels 初始化默认AI模型配置
func InitDefaultAIModels() {
	// 检查是否已有模型配置
	var count int64
	global.DB.Model(&database.AIModel{}).Count(&count)
	if count > 0 {
		return // 已有配置，跳过
	}

	// 创建默认的AI模型配置
	defaultModels := []database.AIModel{
		{
			Name:        "deepseek-r1",
			DisplayName: "DeepSeek R1 (七牛云)",
			Provider:    "qiniu",
			Endpoint:    "https://openai.qiniu.com/v1/chat/completions",
			ApiKey:      "sk-a9f485fe4b59d27f38c8dee28a36143212e98d5e388bf1a37d84c7f6d2f64c0b",
			MaxTokens:   4096,
			Temperature: 0.7,
			IsActive:    true,
		},
	}

	for _, model := range defaultModels {
		if err := global.DB.Create(&model).Error; err != nil {
			global.Log.Error("创建默认AI模型失败: " + err.Error())
		}
	}

	global.Log.Info("默认AI模型配置初始化完成")
}
