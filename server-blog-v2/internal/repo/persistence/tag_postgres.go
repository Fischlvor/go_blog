package persistence

import (
	"context"

	"server-blog-v2/internal/entity"
	"server-blog-v2/internal/repo"
	"server-blog-v2/internal/repo/persistence/gen/model"
	"server-blog-v2/internal/repo/persistence/gen/query"

	"gorm.io/gorm"
)

type tagRepo struct {
	query *query.Query
}

// NewTagRepo 创建标签仓库。
func NewTagRepo(db *gorm.DB) repo.TagRepo {
	return &tagRepo{query: query.Use(db)}
}

func (r *tagRepo) List(ctx context.Context, offset, limit int, keyword *string, sortBy, order *string) ([]*entity.Tag, int64, error) {
	t := r.query.Tag
	do := t.WithContext(ctx)

	if keyword != nil && *keyword != "" {
		kw := "%" + *keyword + "%"
		do = do.Where(t.Name.Like(kw))
	}

	total, err := do.Count()
	if err != nil {
		return nil, 0, err
	}

	do = do.Order(t.ID.Asc())
	rows, err := do.Offset(offset).Limit(limit).Find()
	if err != nil {
		return nil, 0, err
	}

	tags := make([]*entity.Tag, len(rows))
	for i, row := range rows {
		tags[i] = toEntityTag(row)
	}
	return tags, total, nil
}

func (r *tagRepo) ListAll(ctx context.Context) ([]*entity.Tag, error) {
	t := r.query.Tag
	rows, err := t.WithContext(ctx).Order(t.ID.Asc()).Find()
	if err != nil {
		return nil, err
	}

	tags := make([]*entity.Tag, len(rows))
	for i, row := range rows {
		tags[i] = toEntityTag(row)
	}
	return tags, nil
}

func (r *tagRepo) ListByArticleSlug(ctx context.Context, articleSlug string) ([]*entity.Tag, error) {
	t := r.query.Tag
	at := r.query.ArticleTag

	tagIDs, err := at.WithContext(ctx).Where(at.ArticleSlug.Eq(articleSlug)).Select(at.TagID).Find()
	if err != nil {
		return nil, err
	}
	if len(tagIDs) == 0 {
		return []*entity.Tag{}, nil
	}

	ids := make([]int64, len(tagIDs))
	for i, tid := range tagIDs {
		ids[i] = tid.TagID
	}

	rows, err := t.WithContext(ctx).Where(t.ID.In(ids...)).Find()
	if err != nil {
		return nil, err
	}

	tags := make([]*entity.Tag, len(rows))
	for i, row := range rows {
		tags[i] = toEntityTag(row)
	}
	return tags, nil
}

func (r *tagRepo) GetByID(ctx context.Context, id int64) (*entity.Tag, error) {
	t := r.query.Tag
	row, err := t.WithContext(ctx).Where(t.ID.Eq(id)).First()
	if err != nil {
		return nil, err
	}
	return toEntityTag(row), nil
}

func (r *tagRepo) Create(ctx context.Context, tag entity.Tag) (int64, error) {
	mt := toModelTag(&tag)
	if err := r.query.Tag.WithContext(ctx).Create(mt); err != nil {
		return 0, err
	}
	return mt.ID, nil
}

func (r *tagRepo) Update(ctx context.Context, tag entity.Tag) error {
	t := r.query.Tag
	mt := toModelTag(&tag)
	_, err := t.WithContext(ctx).Where(t.ID.Eq(tag.ID)).Updates(mt)
	return err
}

func (r *tagRepo) Delete(ctx context.Context, id int64) error {
	t := r.query.Tag
	_, err := t.WithContext(ctx).Where(t.ID.Eq(id)).Delete()
	return err
}

func toModelTag(t *entity.Tag) *model.Tag {
	return &model.Tag{
		ID:        t.ID,
		Name:      t.Name,
		Slug:      &t.Slug,
		ArticleCount: &t.ArticleCount,
	}
}

func toEntityTag(mt *model.Tag) *entity.Tag {
	tag := &entity.Tag{
		ID:   mt.ID,
		Name: mt.Name,
	}
	if mt.Slug != nil {
		tag.Slug = *mt.Slug
	}
	if mt.ArticleCount != nil {
		tag.ArticleCount = *mt.ArticleCount
	}
	if mt.CreatedAt != nil {
		tag.CreatedAt = *mt.CreatedAt
	}
	if mt.UpdatedAt != nil {
		tag.UpdatedAt = *mt.UpdatedAt
	}
	return tag
}
