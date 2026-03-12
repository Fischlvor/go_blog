package persistence

import (
	"context"

	"server-blog-v2/internal/entity"
	"server-blog-v2/internal/repo"
	"server-blog-v2/internal/repo/persistence/gen/model"
	"server-blog-v2/internal/repo/persistence/gen/query"

	"gorm.io/gorm"
)

type chatMessageRepo struct {
	query *query.Query
}

// NewChatMessageRepo 创建聊天消息仓库。
func NewChatMessageRepo(db *gorm.DB) repo.ChatMessageRepo {
	return &chatMessageRepo{query: query.Use(db)}
}

func (r *chatMessageRepo) Create(ctx context.Context, msg *entity.ChatMessage) (int64, error) {
	mm := toModelChatMessage(msg)
	if err := r.query.AiChatMessage.WithContext(ctx).Create(mm); err != nil {
		return 0, err
	}
	return mm.ID, nil
}

func (r *chatMessageRepo) ListBySessionID(ctx context.Context, sessionID int64) ([]*entity.ChatMessage, error) {
	m := r.query.AiChatMessage
	rows, err := m.WithContext(ctx).Where(m.SessionID.Eq(sessionID)).Order(m.CreatedAt.Asc()).Find()
	if err != nil {
		return nil, err
	}

	messages := make([]*entity.ChatMessage, len(rows))
	for i, row := range rows {
		messages[i] = toEntityChatMessage(row)
	}
	return messages, nil
}

func (r *chatMessageRepo) ListAll(ctx context.Context, offset, limit int, sessionID *int64, role *string) ([]*entity.ChatMessage, int64, error) {
	m := r.query.AiChatMessage
	q := m.WithContext(ctx)

	if sessionID != nil {
		q = q.Where(m.SessionID.Eq(*sessionID))
	}
	if role != nil && *role != "" {
		q = q.Where(m.Role.Eq(*role))
	}

	total, err := q.Count()
	if err != nil {
		return nil, 0, err
	}

	rows, err := q.Order(m.CreatedAt.Desc()).Offset(offset).Limit(limit).Find()
	if err != nil {
		return nil, 0, err
	}

	messages := make([]*entity.ChatMessage, len(rows))
	for i, row := range rows {
		messages[i] = toEntityChatMessage(row)
	}
	return messages, total, nil
}

func toModelChatMessage(m *entity.ChatMessage) *model.AiChatMessage {
	return &model.AiChatMessage{
		ID:        m.ID,
		SessionID: m.SessionID,
		Role:      m.Role,
		Content:   m.Content,
	}
}

func toEntityChatMessage(mm *model.AiChatMessage) *entity.ChatMessage {
	msg := &entity.ChatMessage{
		ID:        mm.ID,
		SessionID: mm.SessionID,
		Role:      mm.Role,
		Content:   mm.Content,
	}
	if mm.Tokens != nil {
		msg.Tokens = int(*mm.Tokens)
	}
	if mm.CreatedAt != nil {
		msg.CreatedAt = *mm.CreatedAt
	}
	return msg
}
