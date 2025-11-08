package service

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"server/internal/model/database"
	"server/internal/model/request"
	"server/internal/model/response"
	"server/internal/service/ai"
	"server/pkg/global"
	"strings"
	"sync"
	"time"

	"github.com/gofrs/uuid"
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
func (s *AIChatService) CreateSession(userUUID uuid.UUID, req request.CreateSessionRequest) (*response.ChatSessionResponse, error) {
	session := &database.AIChatSession{
		UserUUID: userUUID,
		Title:    req.Title,
		Model:    req.Model,
	}

	if err := global.DB.Create(session).Error; err != nil {
		return nil, fmt.Errorf("create session failed: %w", err)
	}

	return &response.ChatSessionResponse{
		ID:        session.ID,
		Title:     session.Title,
		Model:     session.Model,
		UpdatedAt: session.UpdatedAt,
	}, nil
}

// GetSessions 获取会话列表
func (s *AIChatService) GetSessions(userUUID uuid.UUID, req request.GetSessionsRequest) ([]response.ChatSessionResponse, int64, error) {
	var sessions []database.AIChatSession
	var total int64

	// 获取总数
	if err := global.DB.Model(&database.AIChatSession{}).Where("user_uuid = ?", userUUID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取会话列表，按updated_at倒序排列
	query := global.DB.Where("user_uuid = ?", userUUID).Order("updated_at DESC")
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
		responses = append(responses, response.ChatSessionResponse{
			ID:        session.ID,
			Title:     session.Title,
			Model:     session.Model,
			UpdatedAt: session.UpdatedAt,
		})
	}

	return responses, total, nil
}

// GetMessages 获取消息列表
func (s *AIChatService) GetMessages(userUUID uuid.UUID, req request.GetMessagesRequest) ([]response.ChatMessageResponse, int64, error) {
	// 验证会话所有权
	var session database.AIChatSession
	if err := global.DB.Where("id = ? AND user_uuid = ?", req.SessionID, userUUID).First(&session).Error; err != nil {
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
func (s *AIChatService) SendMessage(userUUID uuid.UUID, req request.SendMessageRequest) (*response.ChatResponse, error) {
	// 验证会话所有权
	var session database.AIChatSession
	if err := global.DB.Where("id = ? AND user_uuid = ?", req.SessionID, userUUID).First(&session).Error; err != nil {
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
func (s *AIChatService) SendMessageStream(userUUID uuid.UUID, req request.SendMessageRequest, writer io.Writer) error {
	// 验证会话所有权
	var session database.AIChatSession
	if err := global.DB.Where("id = ? AND user_uuid = ?", req.SessionID, userUUID).First(&session).Error; err != nil {
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

	// 异步生成会话标题（如果是第一次对话且标题还是默认值）
	var wg sync.WaitGroup
	if len(messages) == 0 {
		wg.Add(1)
		go s.generateSessionTitleAsync(userUUID, req.SessionID, req.Content, session.Model, &wg)
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
				Usage struct {
					TotalTokens int `json:"total_tokens"`
				} `json:"usage"`
			}

			if err := json.Unmarshal([]byte(data), &streamResp); err != nil {
				continue
			}

			// 更新token统计（如果存在）
			if streamResp.Usage.TotalTokens > 0 {
				totalTokens = streamResp.Usage.TotalTokens
			}

			if len(streamResp.Choices) > 0 && streamResp.Choices[0].Delta.Content != "" {
				content := streamResp.Choices[0].Delta.Content
				fullContent.WriteString(content)

				// 发送流式响应
				streamResponse := response.StreamingChatResponse{
					Content:   content,
					SessionID: req.SessionID,
					MessageID: aiMessage.ID,
					EventID:   response.EventMessage, // 正常消息
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

	// 更新AI消息内容和token数量
	aiMessage.Content = fullContent.String()
	aiMessage.Tokens = totalTokens
	global.DB.Save(aiMessage)

	// 等待标题生成完成
	wg.Wait()

	// 如果生成了新标题，发送一个特殊的流式消息作为标记
	if len(messages) == 0 {
		// 发送标题生成完成的标记消息
		titleCompleteResponse := response.StreamingChatResponse{
			Content:   "",
			SessionID: req.SessionID,
			MessageID: aiMessage.ID,
			EventID:   response.EventTitleGenerated, // 标题生成完成
		}

		jsonData, _ := json.Marshal(titleCompleteResponse)
		writer.Write([]byte("data: " + string(jsonData) + "\n\n"))
		// 确保数据立即发送
		if flusher, ok := writer.(interface{ Flush() }); ok {
			flusher.Flush()
		}
	}

	// 发送完成信号（统一结束）
	completeResponse := response.StreamingChatResponse{
		Content:   "",
		SessionID: req.SessionID,
		MessageID: aiMessage.ID,
		EventID:   response.EventComplete, // 流式响应完成
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
func (s *AIChatService) DeleteSession(userUUID uuid.UUID, req request.DeleteSessionRequest) error {
	// 验证会话所有权
	var session database.AIChatSession
	if err := global.DB.Where("id = ? AND user_uuid = ?", req.SessionID, userUUID).First(&session).Error; err != nil {
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
func (s *AIChatService) UpdateSession(userUUID uuid.UUID, req request.UpdateSessionRequest) error {
	// 验证会话所有权
	var session database.AIChatSession
	if err := global.DB.Where("id = ? AND user_uuid = ?", req.SessionID, userUUID).First(&session).Error; err != nil {
		return fmt.Errorf("session not found or access denied")
	}

	// 更新会话标题
	session.Title = req.Title
	return global.DB.Save(&session).Error
}

// GetSessionDetail 获取会话详情
func (s *AIChatService) GetSessionDetail(userUUID uuid.UUID, sessionID uint) (*response.ChatSessionResponse, error) {
	var session database.AIChatSession
	if err := global.DB.Where("id = ? AND user_uuid = ?", sessionID, userUUID).First(&session).Error; err != nil {
		return nil, fmt.Errorf("session not found or access denied")
	}

	return &response.ChatSessionResponse{
		ID:        session.ID,
		Title:     session.Title,
		Model:     session.Model,
		UpdatedAt: session.UpdatedAt,
	}, nil
}

// 异步生成会话标题
func (s *AIChatService) generateSessionTitleAsync(userUUID uuid.UUID, sessionID uint, userQuestion string, model string, wg *sync.WaitGroup) {
	defer wg.Done()

	// 构建AI提示词来生成标题
	prompt := fmt.Sprintf(`请根据以下用户问题生成一个简洁的对话标题（不超过20个字符）：

用户问题：%s

要求：
1. 标题要简洁明了，体现对话主题
2. 不超过20个字符
3. 只返回标题，不要其他内容

标题：`, userQuestion)

	// 调用AI服务生成标题
	aiReq := ai.ChatRequest{
		Model:       model, // 使用传入的模型
		Messages:    []ai.Message{{Role: "user", Content: prompt}},
		MaxTokens:   100,
		Temperature: 0.3, // 使用较低的温度确保标题的一致性
	}

	ctx := context.Background()
	aiResp, err := s.aiManager.Chat(ctx, model, aiReq)
	if err != nil {
		global.Log.Error("异步生成标题失败: " + err.Error())
		return
	}

	// 清理AI返回的标题（去除多余的空格、换行等）
	title := strings.TrimSpace(aiResp.Content)
	title = strings.ReplaceAll(title, "\n", "")
	title = strings.ReplaceAll(title, "\r", "")

	// 如果标题过长，截取前20个字符
	if len([]rune(title)) > 20 {
		title = string([]rune(title)[:20])
	}

	// 更新数据库中的会话标题
	if err := global.DB.Model(&database.AIChatSession{}).Where("id = ?", sessionID).Update("title", title).Error; err != nil {
		global.Log.Error("更新会话标题失败: " + err.Error())
		return
	}
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
