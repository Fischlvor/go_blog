package persistence

import (
	"context"

	"server-blog-v2/internal/entity"
	"server-blog-v2/internal/repo"
	"server-blog-v2/internal/repo/persistence/gen/model"
	"server-blog-v2/internal/repo/persistence/gen/query"

	"gorm.io/gorm"
)

type linkRepo struct {
	query *query.Query
}

// NewLinkRepo 创建友链仓库。
func NewLinkRepo(db *gorm.DB) repo.LinkRepo {
	return &linkRepo{query: query.Use(db)}
}

func (r *linkRepo) List(ctx context.Context) ([]*entity.Link, error) {
	l := r.query.Link
	rows, err := l.WithContext(ctx).Where(l.IsVisible.Is(true)).Order(l.Sort.Asc()).Find()
	if err != nil {
		return nil, err
	}

	links := make([]*entity.Link, len(rows))
	for i, row := range rows {
		links[i] = toEntityLink(row)
	}
	return links, nil
}

func (r *linkRepo) GetByID(ctx context.Context, id int64) (*entity.Link, error) {
	l := r.query.Link
	row, err := l.WithContext(ctx).Where(l.ID.Eq(id)).First()
	if err != nil {
		return nil, err
	}
	return toEntityLink(row), nil
}

func (r *linkRepo) Create(ctx context.Context, link *entity.Link) (int64, error) {
	ml := toModelLink(link)
	if err := r.query.Link.WithContext(ctx).Create(ml); err != nil {
		return 0, err
	}
	return ml.ID, nil
}

func (r *linkRepo) Update(ctx context.Context, link *entity.Link) error {
	l := r.query.Link
	ml := toModelLink(link)
	_, err := l.WithContext(ctx).Where(l.ID.Eq(link.ID)).Updates(ml)
	return err
}

func (r *linkRepo) Delete(ctx context.Context, id int64) error {
	l := r.query.Link
	_, err := l.WithContext(ctx).Where(l.ID.Eq(id)).Delete()
	return err
}

func toModelLink(l *entity.Link) *model.Link {
	sort := int32(l.Sort)
	ml := &model.Link{
		ID:        l.ID,
		Name:      l.Name,
		URL:       l.URL,
		Sort:      &sort,
		IsVisible: &l.IsVisible,
	}
	if l.Logo != nil {
		ml.Logo = l.Logo
	}
	if l.Description != nil {
		ml.Description = l.Description
	}
	return ml
}

func toEntityLink(ml *model.Link) *entity.Link {
	link := &entity.Link{
		ID:          ml.ID,
		Name:        ml.Name,
		URL:         ml.URL,
		Logo:        ml.Logo,
		Description: ml.Description,
	}
	if ml.Sort != nil {
		link.Sort = int(*ml.Sort)
	}
	if ml.IsVisible != nil {
		link.IsVisible = *ml.IsVisible
	}
	if ml.CreatedAt != nil {
		link.CreatedAt = *ml.CreatedAt
	}
	if ml.UpdatedAt != nil {
		link.UpdatedAt = *ml.UpdatedAt
	}
	return link
}
