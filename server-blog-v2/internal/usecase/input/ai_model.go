package input

// ListAIModels AI 模型列表参数。
type ListAIModels struct {
	PageParams
	Name     string
	Provider string
}

// CreateAIModel 创建 AI 模型参数。
type CreateAIModel struct {
	Name        string
	DisplayName string
	Provider    string
	Endpoint    string
	ApiKey      string
	MaxTokens   int
	Temperature float64
	IsActive    bool
}

// UpdateAIModel 更新 AI 模型参数。
type UpdateAIModel struct {
	ID          int64
	Name        string
	DisplayName string
	Provider    string
	Endpoint    string
	ApiKey      string
	MaxTokens   int
	Temperature float64
	IsActive    bool
}
