package initialize

import (
	"server/internal/model/database"
	"server/pkg/global"
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
// 从配置文件读取默认模型配置
func InitDefaultAIModels() {
	// 检查是否已有模型配置
	var count int64
	global.DB.Model(&database.AIModel{}).Count(&count)
	if count > 0 {
		return // 已有配置，跳过
	}

	// 从配置文件读取默认模型
	if len(global.Config.AI.DefaultModels) == 0 {
		global.Log.Warn("配置文件中未定义默认AI模型")
		return
	}

	// 创建默认的AI模型配置
	for _, cfg := range global.Config.AI.DefaultModels {
		model := database.AIModel{
			Name:        cfg.Name,
			DisplayName: cfg.DisplayName,
			Provider:    cfg.Provider,
			Endpoint:    cfg.Endpoint,
			ApiKey:      cfg.ApiKey,
			MaxTokens:   cfg.MaxTokens,
			Temperature: cfg.Temperature,
			IsActive:    cfg.IsActive,
		}
		if err := global.DB.Create(&model).Error; err != nil {
			global.Log.Error("创建默认AI模型失败: " + err.Error())
		}
	}

	global.Log.Info("默认AI模型配置初始化完成")
}
