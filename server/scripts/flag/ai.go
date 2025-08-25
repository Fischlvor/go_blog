package flag

import (
	"server/internal/initialize"
)

// AI 初始化AI相关的数据库表和默认模型
func AI() error {
	// 初始化AI相关表
	initialize.InitAITables()

	// 初始化默认AI模型配置
	initialize.InitDefaultAIModels()

	return nil
}
