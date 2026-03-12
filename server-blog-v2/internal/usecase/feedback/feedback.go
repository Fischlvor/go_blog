package feedback

import (
	"context"
	"errors"
	"fmt"

	"server-blog-v2/internal/entity"
	"server-blog-v2/internal/repo"
	"server-blog-v2/internal/usecase"
	"server-blog-v2/internal/usecase/input"
	"server-blog-v2/internal/usecase/output"
)

var ErrRepo = errors.New("repo")

type useCase struct {
	feedbacks repo.FeedbackRepo
}

// New 创建 Feedback UseCase。
func New(feedbacks repo.FeedbackRepo) usecase.Feedback {
	return &useCase{feedbacks: feedbacks}
}

func (u *useCase) Create(ctx context.Context, params input.CreateFeedback) error {
	feedback := &entity.Feedback{
		Type:    params.Type,
		Content: params.Content,
		Status:  entity.FeedbackStatusPending,
	}
	if params.UserUUID != "" {
		feedback.UserUUID = &params.UserUUID
	}
	if params.Contact != "" {
		feedback.Contact = &params.Contact
	}

	_, err := u.feedbacks.Create(ctx, feedback)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrRepo, err)
	}

	return nil
}

func (u *useCase) List(ctx context.Context, params input.ListFeedback) (*output.ListResult[output.Feedback], error) {
	offset := (params.Page - 1) * params.PageSize

	var status *string
	if params.Status != nil {
		status = params.Status
	}

	feedbacks, total, err := u.feedbacks.List(ctx, offset, params.PageSize, status)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrRepo, err)
	}

	items := make([]output.Feedback, len(feedbacks))
	for i, f := range feedbacks {
		item := output.Feedback{
			ID:        f.ID,
			Type:      f.Type,
			Content:   f.Content,
			Status:    f.Status,
			CreatedAt: f.CreatedAt,
			UpdatedAt: f.UpdatedAt,
		}
		if f.UserUUID != nil {
			item.UserUUID = *f.UserUUID
		}
		if f.Contact != nil {
			item.Contact = *f.Contact
		}
		if f.Reply != nil {
			item.Reply = *f.Reply
		}
		items[i] = item
	}

	return &output.ListResult[output.Feedback]{
		Items:    items,
		Page:     params.Page,
		PageSize: params.PageSize,
		Total:    total,
	}, nil
}

// Delete 删除反馈（管理端）。
func (u *useCase) Delete(ctx context.Context, id int64) error {
	if err := u.feedbacks.Delete(ctx, id); err != nil {
		return fmt.Errorf("%w: %v", ErrRepo, err)
	}
	return nil
}

// Reply 回复反馈（管理端）。
func (u *useCase) Reply(ctx context.Context, id int64, reply string) error {
	if err := u.feedbacks.UpdateReply(ctx, id, reply); err != nil {
		return fmt.Errorf("%w: %v", ErrRepo, err)
	}
	return nil
}
