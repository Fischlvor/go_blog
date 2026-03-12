package persistence

import (
	"context"

	"server-blog-v2/internal/entity"
	"server-blog-v2/internal/repo"
	"server-blog-v2/internal/repo/persistence/gen/model"
	"server-blog-v2/internal/repo/persistence/gen/query"

	"gorm.io/gorm"
)

type commentRepo struct {
	query *query.Query
}

// NewCommentRepo 创建评论仓库。
func NewCommentRepo(db *gorm.DB) repo.CommentRepo {
	return &commentRepo{query: query.Use(db)}
}

func (r *commentRepo) ListByArticleSlug(ctx context.Context, articleSlug string, offset, limit int) ([]*entity.Comment, int64, error) {
	c := r.query.Comment
	do := c.WithContext(ctx).Where(c.ArticleSlug.Eq(articleSlug))

	total, err := do.Count()
	if err != nil {
		return nil, 0, err
	}

	do = do.Order(c.CreatedAt.Desc())
	rows, err := do.Offset(offset).Limit(limit).Find()
	if err != nil {
		return nil, 0, err
	}

	comments := make([]*entity.Comment, len(rows))
	for i, row := range rows {
		comments[i] = toEntityComment(row)
	}
	return comments, total, nil
}

func (r *commentRepo) GetByID(ctx context.Context, id int64) (*entity.Comment, error) {
	c := r.query.Comment
	row, err := c.WithContext(ctx).Where(c.ID.Eq(id)).First()
	if err != nil {
		return nil, err
	}
	return toEntityComment(row), nil
}

func (r *commentRepo) Create(ctx context.Context, comment *entity.Comment) (int64, error) {
	mc := toModelComment(comment)
	if err := r.query.Comment.WithContext(ctx).Create(mc); err != nil {
		return 0, err
	}
	return mc.ID, nil
}

func (r *commentRepo) Delete(ctx context.Context, id int64) error {
	c := r.query.Comment
	_, err := c.WithContext(ctx).Where(c.ID.Eq(id)).Delete()
	return err
}

func (r *commentRepo) CountByArticleSlug(ctx context.Context, articleSlug string) (int64, error) {
	c := r.query.Comment
	return c.WithContext(ctx).Where(c.ArticleSlug.Eq(articleSlug)).Count()
}

func (r *commentRepo) ListAll(ctx context.Context, offset, limit int, keyword, articleSlug, userUUID *string) ([]*entity.Comment, int64, error) {
	c := r.query.Comment
	do := c.WithContext(ctx)

	if keyword != nil && *keyword != "" {
		do = do.Where(c.Content.Like("%" + *keyword + "%"))
	}
	if articleSlug != nil && *articleSlug != "" {
		do = do.Where(c.ArticleSlug.Eq(*articleSlug))
	}
	if userUUID != nil && *userUUID != "" {
		do = do.Where(c.UserUUID.Eq(*userUUID))
	}

	total, err := do.Count()
	if err != nil {
		return nil, 0, err
	}

	do = do.Order(c.CreatedAt.Desc())
	rows, err := do.Offset(offset).Limit(limit).Find()
	if err != nil {
		return nil, 0, err
	}

	comments := make([]*entity.Comment, len(rows))
	for i, row := range rows {
		comments[i] = toEntityComment(row)
	}
	return comments, total, nil
}

func toModelComment(c *entity.Comment) *model.Comment {
	mc := &model.Comment{
		ID:          c.ID,
		ArticleSlug: &c.ArticleSlug,
		UserUUID:    &c.UserUUID,
		Content:     c.Content,
		Status:      &c.Status,
		Likes:       &c.Likes,
	}
	if c.ParentID != nil {
		mc.ParentID = c.ParentID
	}
	if c.IPAddress != nil {
		mc.IPAddress = c.IPAddress
	}
	if c.UserAgent != nil {
		mc.UserAgent = c.UserAgent
	}
	return mc
}

func toEntityComment(mc *model.Comment) *entity.Comment {
	cmt := &entity.Comment{
		ID:        mc.ID,
		ParentID:  mc.ParentID,
		Content:   mc.Content,
		IPAddress: mc.IPAddress,
		UserAgent: mc.UserAgent,
	}
	if mc.ArticleSlug != nil {
		cmt.ArticleSlug = *mc.ArticleSlug
	}
	if mc.UserUUID != nil {
		cmt.UserUUID = *mc.UserUUID
	}
	if mc.Status != nil {
		cmt.Status = *mc.Status
	}
	if mc.Likes != nil {
		cmt.Likes = *mc.Likes
	}
	if mc.CreatedAt != nil {
		cmt.CreatedAt = *mc.CreatedAt
	}
	if mc.UpdatedAt != nil {
		cmt.UpdatedAt = *mc.UpdatedAt
	}
	return cmt
}
