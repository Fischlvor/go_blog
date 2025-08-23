package service

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"server/global"
	"server/model/database"
	"server/model/request"
	"server/model/response"
	"server/service/ai"
	"strings"
	"time"
)

// AIChatService AI聊天服务
type AIChatService struct {
	aiManager *ai.AIManager
}

// NewAIChatService 创建AI聊天服务
func NewAIChatService() *AIChatService {
	service := &AIChatService{
		aiManager: ai.NewAIManager(),
	}

	// 初始化AI管理器
	if err := service.aiManager.LoadModelsFromDB(); err != nil {
		global.Log.Error("初始化AI管理器失败: " + err.Error())
	}

	return service
}

// Init 初始化服务
func (s *AIChatService) Init() error {
	return s.aiManager.LoadModelsFromDB()
}

// CreateSession 创建会话
func (s *AIChatService) CreateSession(userID uint, req request.CreateSessionRequest) (*response.ChatSessionResponse, error) {
	session := &database.AIChatSession{
		UserID: userID,
		Title:  req.Title,
		Model:  req.Model,
	}

	if err := global.DB.Create(session).Error; err != nil {
		return nil, fmt.Errorf("create session failed: %w", err)
	}

	return &response.ChatSessionResponse{
		ID:        session.ID,
		Title:     session.Title,
		Model:     session.Model,
		CreatedAt: session.CreatedAt,
		UpdatedAt: session.UpdatedAt,
	}, nil
}

// GetSessions 获取会话列表
func (s *AIChatService) GetSessions(userID uint, req request.GetSessionsRequest) ([]response.ChatSessionResponse, int64, error) {
	var sessions []database.AIChatSession
	var total int64

	// 获取总数
	if err := global.DB.Model(&database.AIChatSession{}).Where("user_id = ?", userID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取会话列表
	query := global.DB.Where("user_id = ?", userID).Order("updated_at DESC")
	if req.Page > 0 && req.PageSize > 0 {
		offset := (req.Page - 1) * req.PageSize
		query = query.Offset(offset).Limit(req.PageSize)
	}

	if err := query.Find(&sessions).Error; err != nil {
		return nil, 0, err
	}

	// 转换为响应格式
	var responses []response.ChatSessionResponse
	for _, session := range sessions {
		// 获取最后一条消息
		var lastMessage database.AIChatMessage
		global.DB.Where("session_id = ?", session.ID).Order("created_at DESC").First(&lastMessage)

		// 获取消息数量
		var messageCount int64
		global.DB.Model(&database.AIChatMessage{}).Where("session_id = ?", session.ID).Count(&messageCount)

		responses = append(responses, response.ChatSessionResponse{
			ID:           session.ID,
			Title:        session.Title,
			Model:        session.Model,
			CreatedAt:    session.CreatedAt,
			UpdatedAt:    session.UpdatedAt,
			LastMessage:  lastMessage.Content,
			MessageCount: int(messageCount),
		})
	}

	return responses, total, nil
}

// GetMessages 获取消息列表
func (s *AIChatService) GetMessages(userID uint, req request.GetMessagesRequest) ([]response.ChatMessageResponse, int64, error) {
	// 验证会话所有权
	var session database.AIChatSession
	if err := global.DB.Where("id = ? AND user_id = ?", req.SessionID, userID).First(&session).Error; err != nil {
		return nil, 0, fmt.Errorf("session not found or access denied")
	}

	var messages []database.AIChatMessage
	var total int64

	// 获取总数
	if err := global.DB.Model(&database.AIChatMessage{}).Where("session_id = ?", req.SessionID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取消息列表
	query := global.DB.Where("session_id = ?", req.SessionID).Order("created_at ASC")
	if req.Page > 0 && req.PageSize > 0 {
		offset := (req.Page - 1) * req.PageSize
		query = query.Offset(offset).Limit(req.PageSize)
	}

	if err := query.Find(&messages).Error; err != nil {
		return nil, 0, err
	}

	// 转换为响应格式
	var responses []response.ChatMessageResponse
	for _, message := range messages {
		responses = append(responses, response.ChatMessageResponse{
			ID:        message.ID,
			Role:      message.Role,
			Content:   message.Content,
			Tokens:    message.Tokens,
			CreatedAt: message.CreatedAt,
		})
	}

	return responses, total, nil
}

// SendMessage 发送消息
func (s *AIChatService) SendMessage(userID uint, req request.SendMessageRequest) (*response.ChatResponse, error) {
	// 验证会话所有权
	var session database.AIChatSession
	if err := global.DB.Where("id = ? AND user_id = ?", req.SessionID, userID).First(&session).Error; err != nil {
		return nil, fmt.Errorf("session not found or access denied")
	}

	// 获取历史消息
	var messages []database.AIChatMessage
	if err := global.DB.Where("session_id = ?", req.SessionID).Order("created_at ASC").Find(&messages).Error; err != nil {
		return nil, err
	}

	// 保存用户消息
	userMessage := &database.AIChatMessage{
		SessionID: req.SessionID,
		Role:      "user",
		Content:   req.Content,
		Tokens:    0,
	}

	if err := global.DB.Create(userMessage).Error; err != nil {
		return nil, err
	}

	// 构建AI请求消息
	var aiMessages []ai.Message
	for _, msg := range messages {
		aiMessages = append(aiMessages, ai.Message{
			Role:    msg.Role,
			Content: msg.Content,
		})
	}
	aiMessages = append(aiMessages, ai.Message{
		Role:    "user",
		Content: req.Content,
	})

	// 调用AI服务
	aiReq := ai.ChatRequest{
		Model:       session.Model,
		Messages:    aiMessages,
		MaxTokens:   4096,
		Temperature: 0.7,
	}

	ctx := context.Background()
	aiResp, err := s.aiManager.Chat(ctx, session.Model, aiReq)
	if err != nil {
		return nil, fmt.Errorf("ai chat failed: %w", err)
	}

	// 保存AI回复
	aiMessage := &database.AIChatMessage{
		SessionID: req.SessionID,
		Role:      "assistant",
		Content:   aiResp.Content,
		Tokens:    aiResp.Tokens,
	}

	if err := global.DB.Create(aiMessage).Error; err != nil {
		return nil, err
	}

	// 更新会话时间
	global.DB.Model(&session).Update("updated_at", time.Now())

	return &response.ChatResponse{
		Message: response.ChatMessageResponse{
			ID:        aiMessage.ID,
			Role:      aiMessage.Role,
			Content:   aiMessage.Content,
			Tokens:    aiMessage.Tokens,
			CreatedAt: aiMessage.CreatedAt,
		},
		SessionID:   req.SessionID,
		TotalTokens: aiResp.Tokens,
	}, nil
}

// SendMessageStream 流式发送消息
func (s *AIChatService) SendMessageStream(userID uint, req request.SendMessageRequest, writer io.Writer) error {
	// 验证会话所有权
	var session database.AIChatSession
	if err := global.DB.Where("id = ? AND user_id = ?", req.SessionID, userID).First(&session).Error; err != nil {
		return fmt.Errorf("session not found or access denied")
	}

	// 获取历史消息
	var messages []database.AIChatMessage
	if err := global.DB.Where("session_id = ?", req.SessionID).Order("created_at ASC").Find(&messages).Error; err != nil {
		return err
	}

	// 保存用户消息
	userMessage := &database.AIChatMessage{
		SessionID: req.SessionID,
		Role:      "user",
		Content:   req.Content,
		Tokens:    0,
	}

	if err := global.DB.Create(userMessage).Error; err != nil {
		return err
	}

	// 构建AI请求消息
	var aiMessages []ai.Message
	for _, msg := range messages {
		aiMessages = append(aiMessages, ai.Message{
			Role:    msg.Role,
			Content: msg.Content,
		})
	}
	aiMessages = append(aiMessages, ai.Message{
		Role:    "user",
		Content: req.Content,
	})

	// 调用AI流式服务
	aiReq := ai.ChatRequest{
		Model:       session.Model,
		Messages:    aiMessages,
		MaxTokens:   4096,
		Temperature: 0.7,
		Stream:      true,
	}

	ctx := context.Background()
	stream, err := s.aiManager.ChatStream(ctx, session.Model, aiReq)
	if err != nil {
		return fmt.Errorf("ai chat stream failed: %w", err)
	}
	defer stream.Close()

	// 创建AI消息记录
	aiMessage := &database.AIChatMessage{
		SessionID: req.SessionID,
		Role:      "assistant",
		Content:   "",
		Tokens:    0,
	}

	if err := global.DB.Create(aiMessage).Error; err != nil {
		return err
	}

	// 处理流式响应
	scanner := bufio.NewScanner(stream)
	var fullContent strings.Builder
	var totalTokens int

	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "data: ") {
			data := strings.TrimPrefix(line, "data: ")
			if data == "[DONE]" {
				break
			}

			var streamResp struct {
				Choices []struct {
					Delta struct {
						Content string `json:"content"`
					} `json:"delta"`
				} `json:"choices"`
			}

			if err := json.Unmarshal([]byte(data), &streamResp); err != nil {
				continue
			}

			if len(streamResp.Choices) > 0 && streamResp.Choices[0].Delta.Content != "" {
				content := streamResp.Choices[0].Delta.Content
				fullContent.WriteString(content)

				// 发送流式响应
				streamResponse := response.StreamingChatResponse{
					Content:     content,
					IsComplete:  false,
					SessionID:   req.SessionID,
					MessageID:   aiMessage.ID,
					TotalTokens: totalTokens,
				}

				jsonData, _ := json.Marshal(streamResponse)
				writer.Write([]byte("data: " + string(jsonData) + "\n\n"))
				// 确保数据立即发送
				if flusher, ok := writer.(interface{ Flush() }); ok {
					flusher.Flush()
				}
			}
		}
	}

	// 更新AI消息内容
	aiMessage.Content = fullContent.String()
	global.DB.Save(aiMessage)

	// 发送完成信号
	completeResponse := response.StreamingChatResponse{
		Content:     "",
		IsComplete:  true,
		SessionID:   req.SessionID,
		MessageID:   aiMessage.ID,
		TotalTokens: totalTokens,
	}

	jsonData, _ := json.Marshal(completeResponse)
	writer.Write([]byte("data: " + string(jsonData) + "\n\n"))
	// 确保数据立即发送
	if flusher, ok := writer.(interface{ Flush() }); ok {
		flusher.Flush()
	}

	// 更新会话时间
	global.DB.Model(&session).Update("updated_at", time.Now())

	return nil
}

// DeleteSession 删除会话
func (s *AIChatService) DeleteSession(userID uint, req request.DeleteSessionRequest) error {
	// 验证会话所有权
	var session database.AIChatSession
	if err := global.DB.Where("id = ? AND user_id = ?", req.SessionID, userID).First(&session).Error; err != nil {
		return fmt.Errorf("session not found or access denied")
	}

	// 删除会话及其消息
	tx := global.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 删除消息
	if err := tx.Where("session_id = ?", req.SessionID).Delete(&database.AIChatMessage{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 删除会话
	if err := tx.Delete(&session).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

// UpdateSession 更新会话
func (s *AIChatService) UpdateSession(userID uint, req request.UpdateSessionRequest) error {
	// 验证会话所有权
	var session database.AIChatSession
	if err := global.DB.Where("id = ? AND user_id = ?", req.SessionID, userID).First(&session).Error; err != nil {
		return fmt.Errorf("session not found or access denied")
	}

	// 更新会话标题
	session.Title = req.Title
	return global.DB.Save(&session).Error
}

// GetAvailableModels 获取可用模型列表
func (s *AIChatService) GetAvailableModels() ([]response.AIModelResponse, error) {
	// 通过AIManager获取模型配置
	modelConfigs, err := s.aiManager.GetAvailableModels()
	if err != nil {
		return nil, err
	}

	var responses []response.AIModelResponse
	for _, config := range modelConfigs {
		responses = append(responses, response.AIModelResponse{
			Name:        config.Name,
			DisplayName: config.DisplayName,
			Provider:    config.Provider,
			MaxTokens:   config.MaxTokens,
			Temperature: config.Temperature,
			IsActive:    config.IsActive,
		})
	}

	return responses, nil
}
