package comment

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
	comments repo.CommentRepo
	users    repo.UserRepo
}

// New 创建 Comment UseCase。
func New(cfg *config.Config, comments repo.CommentRepo, users repo.UserRepo) usecase.Comment {
	return &useCase{
		cfg:      cfg,
		comments: comments,
		users:    users,
	}
}

func (u *useCase) ListByArticleSlug(ctx context.Context, articleSlug string, params input.ListComments) (*output.ListResult[output.Comment], error) {
	offset := (params.Page - 1) * params.PageSize

	comments, total, err := u.comments.ListByArticleSlug(ctx, articleSlug, offset, params.PageSize)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrRepo, err)
	}

	items := make([]output.Comment, len(comments))
	for i, c := range comments {
		user := u.getUserInfo(ctx, c.UserUUID)
		items[i] = output.Comment{
			ID:          c.ID,
			ArticleSlug: c.ArticleSlug,
			ParentID:    c.ParentID,
			Content:     c.Content,
			User:        user,
			CreatedAt:   c.CreatedAt,
		}
	}

	return &output.ListResult[output.Comment]{
		Items:    items,
		Page:     params.Page,
		PageSize: params.PageSize,
		Total:    total,
	}, nil
}

func (u *useCase) Create(ctx context.Context, params input.CreateComment) (int64, error) {
	comment := &entity.Comment{
		ArticleSlug: params.ArticleSlug,
		UserUUID:    params.UserUUID,
		Content:     params.Content,
		ParentID:    params.ParentID,
		Status:      entity.CommentStatusApproved,
	}

	id, err := u.comments.Create(ctx, comment)
	if err != nil {
		return 0, fmt.Errorf("%w: %v", ErrRepo, err)
	}

	return id, nil
}

func (u *useCase) Delete(ctx context.Context, id int64, userUUID string) error {
	// 检查评论是否属于该用户
	comment, err := u.comments.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrRepo, err)
	}

	if comment.UserUUID != userUUID {
		return ErrForbidden
	}

	if err := u.comments.Delete(ctx, id); err != nil {
		return fmt.Errorf("%w: %v", ErrRepo, err)
	}

	return nil
}

func (u *useCase) getUserInfo(ctx context.Context, userUUID string) output.CommentUser {
	user, err := u.users.GetByUUID(ctx, userUUID)
	if err != nil || user == nil {
		return output.CommentUser{UUID: userUUID}
	}
	return output.CommentUser{
		UUID:     user.UUID,
		Nickname: user.Nickname,
		Avatar:   urlutil.ResolveImageURL(u.cfg, user.Avatar),
	}
}

// ListAll 评论列表（管理端）。
func (u *useCase) ListAll(ctx context.Context, params input.ListAllComments) (*output.ListResult[output.CommentAdmin], error) {
	offset := (params.Page - 1) * params.PageSize

	var keyword *string
	if params.Keyword != nil {
		keyword = &params.Keyword.Keyword
	}

	comments, total, err := u.comments.ListAll(ctx, offset, params.PageSize, keyword, params.ArticleSlug, params.UserUUID)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrRepo, err)
	}

	items := make([]output.CommentAdmin, len(comments))
	for i, c := range comments {
		user := u.getUserInfo(ctx, c.UserUUID)
		items[i] = output.CommentAdmin{
			ID:          c.ID,
			ArticleSlug: c.ArticleSlug,
			Content:     c.Content,
			User:        user,
			ParentID:    c.ParentID,
			CreatedAt:   c.CreatedAt,
		}
	}

	return &output.ListResult[output.CommentAdmin]{
		Items:    items,
		Page:     params.Page,
		PageSize: params.PageSize,
		Total:    total,
	}, nil
}

// AdminDelete 删除评论（管理端）。
func (u *useCase) AdminDelete(ctx context.Context, id int64) error {
	if err := u.comments.Delete(ctx, id); err != nil {
		return fmt.Errorf("%w: %v", ErrRepo, err)
	}
	return nil
}
