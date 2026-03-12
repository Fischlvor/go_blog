package persistence

import (
	"context"

	"server-blog-v2/internal/entity"
	"server-blog-v2/internal/repo"
	"server-blog-v2/internal/repo/persistence/gen/model"
	"server-blog-v2/internal/repo/persistence/gen/query"

	"gorm.io/gorm"
)

type feedbackRepo struct {
	query *query.Query
}

// NewFeedbackRepo 创建反馈仓库。
func NewFeedbackRepo(db *gorm.DB) repo.FeedbackRepo {
	return &feedbackRepo{query: query.Use(db)}
}

func (r *feedbackRepo) Create(ctx context.Context, feedback *entity.Feedback) (int64, error) {
	mf := toModelFeedback(feedback)
	if err := r.query.Feedback.WithContext(ctx).Create(mf); err != nil {
		return 0, err
	}
	return mf.ID, nil
}

func (r *feedbackRepo) List(ctx context.Context, offset, limit int, status *string) ([]*entity.Feedback, int64, error) {
	f := r.query.Feedback
	do := f.WithContext(ctx)

	if status != nil && *status != "" {
		do = do.Where(f.Status.Eq(*status))
	}

	total, err := do.Count()
	if err != nil {
		return nil, 0, err
	}

	do = do.Order(f.CreatedAt.Desc())
	rows, err := do.Offset(offset).Limit(limit).Find()
	if err != nil {
		return nil, 0, err
	}

	feedbacks := make([]*entity.Feedback, len(rows))
	for i, row := range rows {
		feedbacks[i] = toEntityFeedback(row)
	}
	return feedbacks, total, nil
}

func (r *feedbackRepo) GetByID(ctx context.Context, id int64) (*entity.Feedback, error) {
	f := r.query.Feedback
	row, err := f.WithContext(ctx).Where(f.ID.Eq(id)).First()
	if err != nil {
		return nil, err
	}
	return toEntityFeedback(row), nil
}

func (r *feedbackRepo) Delete(ctx context.Context, id int64) error {
	f := r.query.Feedback
	_, err := f.WithContext(ctx).Where(f.ID.Eq(id)).Delete()
	return err
}

func (r *feedbackRepo) UpdateStatus(ctx context.Context, id int64, status string) error {
	f := r.query.Feedback
	_, err := f.WithContext(ctx).Where(f.ID.Eq(id)).UpdateSimple(f.Status.Value(status))
	return err
}

func (r *feedbackRepo) UpdateReply(ctx context.Context, id int64, reply string) error {
	f := r.query.Feedback
	_, err := f.WithContext(ctx).Where(f.ID.Eq(id)).UpdateSimple(
		f.Reply.Value(reply),
		f.Status.Value("replied"),
	)
	return err
}

func toModelFeedback(f *entity.Feedback) *model.Feedback {
	mf := &model.Feedback{
		ID:       f.ID,
		UserUUID: f.UserUUID,
		Type:     &f.Type,
		Content:  f.Content,
		Status:   &f.Status,
	}
	if f.Contact != nil {
		mf.Contact = f.Contact
	}
	if f.Reply != nil {
		mf.Reply = f.Reply
	}
	return mf
}

func toEntityFeedback(mf *model.Feedback) *entity.Feedback {
	fb := &entity.Feedback{
		ID:       mf.ID,
		UserUUID: mf.UserUUID,
		Content:  mf.Content,
		Contact:  mf.Contact,
		Reply:    mf.Reply,
	}
	if mf.Type != nil {
		fb.Type = *mf.Type
	}
	if mf.Status != nil {
		fb.Status = *mf.Status
	}
	if mf.CreatedAt != nil {
		fb.CreatedAt = *mf.CreatedAt
	}
	if mf.UpdatedAt != nil {
		fb.UpdatedAt = *mf.UpdatedAt
	}
	return fb
}
