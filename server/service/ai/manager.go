package ai

import (
	"context"
	"fmt"
	"io"
	"server/global"
	"server/model/database"
	"sync"
)

// AIManager AI服务管理器
type AIManager struct {
	providers map[string]AIModelProvider
	mu        sync.RWMutex
}

// NewAIManager 创建AI管理器
func NewAIManager() *AIManager {
	return &AIManager{
		providers: make(map[string]AIModelProvider),
	}
}

// RegisterProvider 注册AI提供商
func (m *AIManager) RegisterProvider(name string, provider AIModelProvider) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.providers[name] = provider
}

// GetProvider 获取AI提供商
func (m *AIManager) GetProvider(name string) (AIModelProvider, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	provider, exists := m.providers[name]
	if !exists {
		return nil, fmt.Errorf("provider %s not found", name)
	}

	if !provider.IsAvailable() {
		return nil, fmt.Errorf("provider %s is not available", name)
	}

	return provider, nil
}

// GetAvailableProviders 获取可用的提供商列表
func (m *AIManager) GetAvailableProviders() []string {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var available []string
	for name, provider := range m.providers {
		if provider.IsAvailable() {
			available = append(available, name)
		}
	}
	return available
}

// GetAvailableModels 获取可用的模型配置
func (m *AIManager) GetAvailableModels() ([]ModelConfig, error) {
	// 从数据库获取活跃的模型配置
	var models []database.AIModel
	if err := global.DB.Where("is_active = ?", true).Find(&models).Error; err != nil {
		return nil, err
	}

	var configs []ModelConfig
	for _, model := range models {
		configs = append(configs, ModelConfig{
			Name:        model.Name,
			DisplayName: model.DisplayName,
			Provider:    model.Provider,
			Endpoint:    model.Endpoint,
			ApiKey:      model.ApiKey,
			MaxTokens:   model.MaxTokens,
			Temperature: model.Temperature,
			IsActive:    model.IsActive,
		})
	}

	return configs, nil
}

// LoadModelsFromDB 从数据库加载模型配置
func (m *AIManager) LoadModelsFromDB() error {
	var models []database.AIModel
	if err := global.DB.Where("is_active = ?", true).Find(&models).Error; err != nil {
		return fmt.Errorf("load models from db failed: %w", err)
	}

	for _, model := range models {
		config := ModelConfig{
			Name:        model.Name,
			DisplayName: model.DisplayName,
			Provider:    model.Provider,
			Endpoint:    model.Endpoint,
			ApiKey:      model.ApiKey,
			MaxTokens:   model.MaxTokens,
			Temperature: model.Temperature,
			IsActive:    model.IsActive,
		}

		// 根据模型名称选择对应的提供商
		switch model.Name {
		case "deepseek-r1":
			provider := NewQiniuDeepSeekProvider(config)
			m.RegisterProvider(model.Name, provider)
		// 可以在这里添加其他模型
		default:
			global.Log.Warn("unsupported model: " + model.Name)
		}
	}

	return nil
}

// Chat 聊天
func (m *AIManager) Chat(ctx context.Context, modelName string, req ChatRequest) (*ChatResponse, error) {
	provider, err := m.GetProvider(modelName)
	if err != nil {
		return nil, err
	}

	return provider.Chat(ctx, req)
}

// ChatStream 流式聊天
func (m *AIManager) ChatStream(ctx context.Context, modelName string, req ChatRequest) (io.ReadCloser, error) {
	provider, err := m.GetProvider(modelName)
	if err != nil {
		return nil, err
	}

	return provider.ChatStream(ctx, req)
}
