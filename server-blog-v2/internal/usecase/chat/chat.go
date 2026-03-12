package chat

import (
	"context"
	"errors"
	"fmt"

	"server-blog-v2/config"
	"server-blog-v2/internal/entity"
	"server-blog-v2/internal/repo"
	"server-blog-v2/internal/usecase"
	"server-blog-v2/internal/usecase/input"
	"server-blog-v2/internal/usecase/output"
	"server-blog-v2/internal/usecase/urlutil"
)

var (
	ErrRepo      = errors.New("repo")
	ErrForbidden = errors.New("forbidden")
)

type useCase struct {
	cfg      *config.Config
	sessions repo.ChatSessionRepo
	messages repo.ChatMessageRepo
	llm      repo.LLMWebAPI
	users    repo.UserRepo
}

// New 创建 AIChat UseCase。
func New(
	cfg *config.Config,
	sessions repo.ChatSessionRepo,
	messages repo.ChatMessageRepo,
	llm repo.LLMWebAPI,
	users repo.UserRepo,
) usecase.AIChat {
	return &useCase{
		cfg:      cfg,
		sessions: sessions,
		messages: messages,
		llm:      llm,
		users:    users,
	}
}

func (u *useCase) CreateSession(ctx context.Context, userUUID string, params input.CreateSession) (*output.Session, error) {
	session := &entity.ChatSession{
		UserUUID: userUUID,
		Title:    params.Title,
	}

	id, err := u.sessions.Create(ctx, session)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrRepo, err)
	}

	created, err := u.sessions.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrRepo, err)
	}

	return &output.Session{
		ID:        created.ID,
		Title:     created.Title,
		CreatedAt: created.CreatedAt,
		UpdatedAt: created.UpdatedAt,
	}, nil
}

func (u *useCase) ListSessions(ctx context.Context, userUUID string, params input.ListSessions) (*output.ListResult[output.Session], error) {
	offset := (params.Page - 1) * params.PageSize

	sessions, total, err := u.sessions.List(ctx, userUUID, offset, params.PageSize)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrRepo, err)
	}

	items := make([]output.Session, len(sessions))
	for i, s := range sessions {
		items[i] = output.Session{
			ID:        s.ID,
			Title:     s.Title,
			CreatedAt: s.CreatedAt,
			UpdatedAt: s.UpdatedAt,
		}
	}

	return &output.ListResult[output.Session]{
		Items:    items,
		Page:     params.Page,
		PageSize: params.PageSize,
		Total:    total,
	}, nil
}

func (u *useCase) GetMessages(ctx context.Context, sessionID int64, userUUID string) ([]*output.Message, error) {
	// 检查会话是否属于该用户
	session, err := u.sessions.GetByID(ctx, sessionID)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrRepo, err)
	}
	if session.UserUUID != userUUID {
		return nil, ErrForbidden
	}

	messages, err := u.messages.ListBySessionID(ctx, sessionID)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrRepo, err)
	}

	items := make([]*output.Message, len(messages))
	for i, m := range messages {
		items[i] = &output.Message{
			ID:        m.ID,
			Role:      m.Role,
			Content:   m.Content,
			CreatedAt: m.CreatedAt,
		}
	}

	return items, nil
}

func (u *useCase) SendMessage(ctx context.Context, sessionID int64, userUUID string, params input.SendMessage) (<-chan string, error) {
	// 检查会话是否属于该用户
	session, err := u.sessions.GetByID(ctx, sessionID)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrRepo, err)
	}
	if session.UserUUID != userUUID {
		return nil, ErrForbidden
	}

	// 保存用户消息
	userMsg := &entity.ChatMessage{
		SessionID: sessionID,
		Role:      "user",
		Content:   params.Content,
	}
	if _, err := u.messages.Create(ctx, userMsg); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrRepo, err)
	}

	// 获取历史消息
	history, err := u.messages.ListBySessionID(ctx, sessionID)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrRepo, err)
	}

	// 构建 LLM 消息
	llmMessages := make([]repo.LLMMessage, 0, len(history)+1)
	llmMessages = append(llmMessages, repo.LLMMessage{
		Role:    "system",
		Content: "你是一个友好的 AI 助手，帮助用户解答问题。",
	})
	for _, m := range history {
		llmMessages = append(llmMessages, repo.LLMMessage{
			Role:    m.Role,
			Content: m.Content,
		})
	}

	// 调用 LLM
	stream, err := u.llm.ChatStream(ctx, llmMessages)
	if err != nil {
		return nil, fmt.Errorf("llm error: %w", err)
	}

	// 创建输出 channel，收集完整响应后保存
	out := make(chan string, 100)
	go func() {
		defer close(out)
		var fullContent string
		for chunk := range stream {
			fullContent += chunk
			out <- chunk
		}

		// 保存助手消息
		if fullContent != "" {
			assistantMsg := &entity.ChatMessage{
				SessionID: sessionID,
				Role:      "assistant",
				Content:   fullContent,
			}
			_, _ = u.messages.Create(ctx, assistantMsg)

			// 如果是第一条消息，更新会话标题
			if len(history) <= 1 && session.Title == "" {
				title := fullContent
				if len(title) > 50 {
					title = title[:50] + "..."
				}
				_ = u.sessions.UpdateTitle(ctx, sessionID, title)
			}
		}
	}()

	return out, nil
}

func (u *useCase) DeleteSession(ctx context.Context, sessionID int64, userUUID string) error {
	// 检查会话是否属于该用户
	session, err := u.sessions.GetByID(ctx, sessionID)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrRepo, err)
	}
	if session.UserUUID != userUUID {
		return ErrForbidden
	}

	if err := u.sessions.Delete(ctx, sessionID); err != nil {
		return fmt.Errorf("%w: %v", ErrRepo, err)
	}

	return nil
}

// ListAllSessions 会话列表（管理端）。
func (u *useCase) ListAllSessions(ctx context.Context, params input.ListAllSessions) (*output.ListResult[output.SessionAdmin], error) {
	offset := (params.Page - 1) * params.PageSize

	var keyword *string
	if params.Keyword != nil {
		keyword = &params.Keyword.Keyword
	}

	sessions, total, err := u.sessions.ListAll(ctx, offset, params.PageSize, keyword, params.UserUUID)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrRepo, err)
	}

	items := make([]output.SessionAdmin, len(sessions))
	for i, s := range sessions {
		// 获取用户信息
		user, _ := u.users.GetByUUID(ctx, s.UserUUID)

		// 获取消息数量
		messages, _ := u.messages.ListBySessionID(ctx, s.ID)

		item := output.SessionAdmin{
			ID:           s.ID,
			Title:        s.Title,
			UserUUID:     s.UserUUID,
			MessageCount: len(messages),
			CreatedAt:    s.CreatedAt,
			UpdatedAt:    s.UpdatedAt,
		}
		if user != nil {
			item.User = &output.SessionUser{
				UUID:     user.UUID,
				Username: user.Nickname,
				Avatar:   urlutil.ResolveImageURL(u.cfg, user.Avatar),
			}
		}
		items[i] = item
	}

	return &output.ListResult[output.SessionAdmin]{
		Items:    items,
		Page:     params.Page,
		PageSize: params.PageSize,
		Total:    total,
	}, nil
}

// GetSessionByID 获取会话详情（管理端）。
func (u *useCase) GetSessionByID(ctx context.Context, sessionID int64) (*output.SessionAdmin, error) {
	session, err := u.sessions.GetByID(ctx, sessionID)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrRepo, err)
	}

	// 获取用户信息
	user, _ := u.users.GetByUUID(ctx, session.UserUUID)

	// 获取消息数量
	messages, _ := u.messages.ListBySessionID(ctx, sessionID)

	result := &output.SessionAdmin{
		ID:           session.ID,
		Title:        session.Title,
		UserUUID:     session.UserUUID,
		MessageCount: len(messages),
		CreatedAt:    session.CreatedAt,
		UpdatedAt:    session.UpdatedAt,
	}
	if user != nil {
		result.User = &output.SessionUser{
			UUID:     user.UUID,
			Username: user.Nickname,
			Avatar:   urlutil.ResolveImageURL(u.cfg, user.Avatar),
		}
	}
	return result, nil
}

// AdminDeleteSession 删除会话（管理端）。
func (u *useCase) AdminDeleteSession(ctx context.Context, sessionID int64) error {
	if err := u.sessions.Delete(ctx, sessionID); err != nil {
		return fmt.Errorf("%w: %v", ErrRepo, err)
	}
	return nil
}

// ListAllMessages 消息列表（管理端）。
func (u *useCase) ListAllMessages(ctx context.Context, params input.ListAllMessages) (*output.ListResult[output.MessageAdmin], error) {
	offset := (params.Page - 1) * params.PageSize

	messages, total, err := u.messages.ListAll(ctx, offset, params.PageSize, params.SessionID, params.Role)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrRepo, err)
	}

	items := make([]output.MessageAdmin, len(messages))
	for i, m := range messages {
		items[i] = output.MessageAdmin{
			ID:        m.ID,
			SessionID: m.SessionID,
			Role:      m.Role,
			Content:   m.Content,
			Tokens:    m.Tokens,
			CreatedAt: m.CreatedAt,
		}
	}

	return &output.ListResult[output.MessageAdmin]{
		Items:    items,
		Page:     params.Page,
		PageSize: params.PageSize,
		Total:    total,
	}, nil
}
