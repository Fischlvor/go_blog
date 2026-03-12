package webapi

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"server-blog-v2/internal/repo"
)

type llmWebAPI struct {
	apiKey  string
	baseURL string
	model   string
	client  *http.Client
}

// NewLLMWebAPI 创建 LLM API 客户端。
func NewLLMWebAPI(apiKey, baseURL, model string) repo.LLMWebAPI {
	return &llmWebAPI{
		apiKey:  apiKey,
		baseURL: baseURL,
		model:   model,
		client:  &http.Client{},
	}
}

type chatRequest struct {
	Model    string          `json:"model"`
	Messages []chatMessage   `json:"messages"`
	Stream   bool            `json:"stream"`
}

type chatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type streamResponse struct {
	Choices []struct {
		Delta struct {
			Content string `json:"content"`
		} `json:"delta"`
		FinishReason *string `json:"finish_reason"`
	} `json:"choices"`
}

func (l *llmWebAPI) ChatStream(ctx context.Context, messages []repo.LLMMessage) (<-chan string, error) {
	// 转换消息格式
	chatMessages := make([]chatMessage, len(messages))
	for i, m := range messages {
		chatMessages[i] = chatMessage{
			Role:    m.Role,
			Content: m.Content,
		}
	}

	// 构建请求
	reqBody := chatRequest{
		Model:    l.model,
		Messages: chatMessages,
		Stream:   true,
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", l.baseURL+"/chat/completions", bytes.NewReader(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+l.apiKey)

	resp, err := l.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("send request: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		return nil, fmt.Errorf("api error: %s, body: %s", resp.Status, string(body))
	}

	// 创建输出 channel
	ch := make(chan string, 100)

	go func() {
		defer close(ch)
		defer resp.Body.Close()

		reader := bufio.NewReader(resp.Body)
		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				if err != io.EOF {
					// 发送错误信息
					ch <- fmt.Sprintf("[ERROR] %v", err)
				}
				return
			}

			line = strings.TrimSpace(line)
			if line == "" || line == "data: [DONE]" {
				continue
			}

			if !strings.HasPrefix(line, "data: ") {
				continue
			}

			data := strings.TrimPrefix(line, "data: ")
			var streamResp streamResponse
			if err := json.Unmarshal([]byte(data), &streamResp); err != nil {
				continue
			}

			if len(streamResp.Choices) > 0 {
				content := streamResp.Choices[0].Delta.Content
				if content != "" {
					select {
					case ch <- content:
					case <-ctx.Done():
						return
					}
				}

				if streamResp.Choices[0].FinishReason != nil {
					return
				}
			}
		}
	}()

	return ch, nil
}
