package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// QiniuDeepSeekProvider 七牛云DeepSeek提供商
type QiniuDeepSeekProvider struct {
	config ModelConfig
	client *http.Client
}

// NewQiniuDeepSeekProvider 创建七牛云DeepSeek提供商
func NewQiniuDeepSeekProvider(config ModelConfig) *QiniuDeepSeekProvider {
	return &QiniuDeepSeekProvider{
		config: config,
		client: &http.Client{
			Timeout: 60 * time.Second,
		},
	}
}

// Chat 普通聊天
func (p *QiniuDeepSeekProvider) Chat(ctx context.Context, req ChatRequest) (*ChatResponse, error) {
	// 构建七牛云DeepSeek请求
	qiniuReq := map[string]interface{}{
		"model":       req.Model,
		"messages":    req.Messages,
		"max_tokens":  req.MaxTokens,
		"temperature": req.Temperature,
		"stream":      false,
	}

	jsonData, err := json.Marshal(qiniuReq)
	if err != nil {
		return nil, fmt.Errorf("marshal request failed: %w", err)
	}

	// 创建HTTP请求
	httpReq, err := http.NewRequestWithContext(ctx, "POST", p.config.Endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("create request failed: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+p.config.ApiKey)

	// 发送请求
	resp, err := p.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("send request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("qiniu deepseek api error: %s, body: %s", resp.Status, string(body))
	}

	// 解析响应
	var qiniuResp struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
			FinishReason string `json:"finish_reason"`
		} `json:"choices"`
		Usage struct {
			TotalTokens int `json:"total_tokens"`
		} `json:"usage"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&qiniuResp); err != nil {
		return nil, fmt.Errorf("decode response failed: %w", err)
	}

	if len(qiniuResp.Choices) == 0 {
		return nil, fmt.Errorf("no choices in response")
	}

	return &ChatResponse{
		Content:      qiniuResp.Choices[0].Message.Content,
		Tokens:       qiniuResp.Usage.TotalTokens,
		FinishReason: qiniuResp.Choices[0].FinishReason,
	}, nil
}

// ChatStream 流式聊天
func (p *QiniuDeepSeekProvider) ChatStream(ctx context.Context, req ChatRequest) (io.ReadCloser, error) {
	// 构建七牛云DeepSeek请求
	qiniuReq := map[string]interface{}{
		"model":       req.Model,
		"messages":    req.Messages,
		"max_tokens":  req.MaxTokens,
		"temperature": req.Temperature,
		"stream":      true,
	}

	jsonData, err := json.Marshal(qiniuReq)
	if err != nil {
		return nil, fmt.Errorf("marshal request failed: %w", err)
	}

	// 创建HTTP请求
	httpReq, err := http.NewRequestWithContext(ctx, "POST", p.config.Endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("create request failed: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+p.config.ApiKey)

	// 发送请求
	resp, err := p.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("send request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("qiniu deepseek api error: %s, body: %s", resp.Status, string(body))
	}

	return resp.Body, nil
}

// GetName 获取提供商名称
func (p *QiniuDeepSeekProvider) GetName() string {
	return "qiniu-deepseek"
}

// IsAvailable 检查是否可用
func (p *QiniuDeepSeekProvider) IsAvailable() bool {
	return p.config.ApiKey != "" && p.config.Endpoint != ""
}
