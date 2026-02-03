package config

// AI AI模型配置
type AI struct {
	DefaultModels []AIModel `json:"default_models" yaml:"default_models"`
}

// AIModel 单个AI模型配置
type AIModel struct {
	Name        string  `json:"name" yaml:"name"`                 // 模型名称
	DisplayName string  `json:"display_name" yaml:"display_name"` // 显示名称
	Provider    string  `json:"provider" yaml:"provider"`         // 提供商
	Endpoint    string  `json:"endpoint" yaml:"endpoint"`         // API端点
	ApiKey      string  `json:"api_key" yaml:"api_key"`           // API密钥
	MaxTokens   int     `json:"max_tokens" yaml:"max_tokens"`     // 最大token数
	Temperature float64 `json:"temperature" yaml:"temperature"`   // 温度参数
	IsActive    bool    `json:"is_active" yaml:"is_active"`       // 是否启用
}
