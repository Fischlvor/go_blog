package persistence

import (
	"context"

	"server-blog-v2/internal/entity"
	"server-blog-v2/internal/repo"
	"server-blog-v2/internal/repo/persistence/gen/model"
	"server-blog-v2/internal/repo/persistence/gen/query"

	"gorm.io/gorm"
)

type chatSessionRepo struct {
	query *query.Query
}

// NewChatSessionRepo 创建聊天会话仓库。
func NewChatSessionRepo(db *gorm.DB) repo.ChatSessionRepo {
	return &chatSessionRepo{query: query.Use(db)}
}

func (r *chatSessionRepo) Create(ctx context.Context, session *entity.ChatSession) (int64, error) {
	ms := toModelChatSession(session)
	if err := r.query.AiChatSession.WithContext(ctx).Create(ms); err != nil {
		return 0, err
	}
	return ms.ID, nil
}

func (r *chatSessionRepo) GetByID(ctx context.Context, id int64) (*entity.ChatSession, error) {
	s := r.query.AiChatSession
	row, err := s.WithContext(ctx).Where(s.ID.Eq(id)).First()
	if err != nil {
		return nil, err
	}
	return toEntityChatSession(row), nil
}

func (r *chatSessionRepo) List(ctx context.Context, userUUID string, offset, limit int) ([]*entity.ChatSession, int64, error) {
	s := r.query.AiChatSession
	do := s.WithContext(ctx).Where(s.UserUUID.Eq(userUUID))

	total, err := do.Count()
	if err != nil {
		return nil, 0, err
	}

	do = do.Order(s.UpdatedAt.Desc())
	rows, err := do.Offset(offset).Limit(limit).Find()
	if err != nil {
		return nil, 0, err
	}

	sessions := make([]*entity.ChatSession, len(rows))
	for i, row := range rows {
		sessions[i] = toEntityChatSession(row)
	}
	return sessions, total, nil
}

func (r *chatSessionRepo) Delete(ctx context.Context, id int64) error {
	s := r.query.AiChatSession
	_, err := s.WithContext(ctx).Where(s.ID.Eq(id)).Delete()
	return err
}

func (r *chatSessionRepo) UpdateTitle(ctx context.Context, id int64, title string) error {
	s := r.query.AiChatSession
	_, err := s.WithContext(ctx).Where(s.ID.Eq(id)).UpdateSimple(s.Title.Value(title))
	return err
}

func (r *chatSessionRepo) ListAll(ctx context.Context, offset, limit int, keyword, userUUID *string) ([]*entity.ChatSession, int64, error) {
	s := r.query.AiChatSession
	do := s.WithContext(ctx)

	if keyword != nil && *keyword != "" {
		do = do.Where(s.Title.Like("%" + *keyword + "%"))
	}
	if userUUID != nil && *userUUID != "" {
		do = do.Where(s.UserUUID.Eq(*userUUID))
	}

	total, err := do.Count()
	if err != nil {
		return nil, 0, err
	}

	do = do.Order(s.UpdatedAt.Desc())
	rows, err := do.Offset(offset).Limit(limit).Find()
	if err != nil {
		return nil, 0, err
	}

	sessions := make([]*entity.ChatSession, len(rows))
	for i, row := range rows {
		sessions[i] = toEntityChatSession(row)
	}
	return sessions, total, nil
}

func toModelChatSession(s *entity.ChatSession) *model.AiChatSession {
	return &model.AiChatSession{
		ID:       s.ID,
		UserUUID: &s.UserUUID,
		Title:    &s.Title,
	}
}

func toEntityChatSession(ms *model.AiChatSession) *entity.ChatSession {
	sess := &entity.ChatSession{
		ID: ms.ID,
	}
	if ms.UserUUID != nil {
		sess.UserUUID = *ms.UserUUID
	}
	if ms.Title != nil {
		sess.Title = *ms.Title
	}
	if ms.CreatedAt != nil {
		sess.CreatedAt = *ms.CreatedAt
	}
	if ms.UpdatedAt != nil {
		sess.UpdatedAt = *ms.UpdatedAt
	}
	return sess
}
