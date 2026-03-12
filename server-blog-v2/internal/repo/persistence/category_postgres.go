package persistence

import (
	"context"

	"server-blog-v2/internal/entity"
	"server-blog-v2/internal/repo"
	"server-blog-v2/internal/repo/persistence/gen/model"
	"server-blog-v2/internal/repo/persistence/gen/query"

	"gorm.io/gorm"
)

type categoryRepo struct {
	query *query.Query
}

// NewCategoryRepo 创建分类仓库。
func NewCategoryRepo(db *gorm.DB) repo.CategoryRepo {
	return &categoryRepo{query: query.Use(db)}
}

func (r *categoryRepo) List(ctx context.Context, offset, limit int, keyword *string, sortBy, order *string) ([]*entity.Category, int64, error) {
	c := r.query.Category
	do := c.WithContext(ctx)

	if keyword != nil && *keyword != "" {
		kw := "%" + *keyword + "%"
		do = do.Where(c.Name.Like(kw))
	}

	total, err := do.Count()
	if err != nil {
		return nil, 0, err
	}

	do = do.Order(c.ID.Asc())
	rows, err := do.Offset(offset).Limit(limit).Find()
	if err != nil {
		return nil, 0, err
	}

	categories := make([]*entity.Category, len(rows))
	for i, row := range rows {
		categories[i] = toEntityCategory(row)
	}
	return categories, total, nil
}

func (r *categoryRepo) ListAll(ctx context.Context) ([]*entity.Category, error) {
	c := r.query.Category
	rows, err := c.WithContext(ctx).Order(c.ID.Asc()).Find()
	if err != nil {
		return nil, err
	}

	categories := make([]*entity.Category, len(rows))
	for i, row := range rows {
		categories[i] = toEntityCategory(row)
	}
	return categories, nil
}

func (r *categoryRepo) GetByID(ctx context.Context, id int64) (*entity.Category, error) {
	c := r.query.Category
	row, err := c.WithContext(ctx).Where(c.ID.Eq(id)).First()
	if err != nil {
		return nil, err
	}
	return toEntityCategory(row), nil
}

func (r *categoryRepo) Create(ctx context.Context, category entity.Category) (int64, error) {
	mc := toModelCategory(&category)
	if err := r.query.Category.WithContext(ctx).Create(mc); err != nil {
		return 0, err
	}
	return mc.ID, nil
}

func (r *categoryRepo) Update(ctx context.Context, category entity.Category) error {
	c := r.query.Category
	mc := toModelCategory(&category)
	_, err := c.WithContext(ctx).Where(c.ID.Eq(category.ID)).Updates(mc)
	return err
}

func (r *categoryRepo) Delete(ctx context.Context, id int64) error {
	c := r.query.Category
	_, err := c.WithContext(ctx).Where(c.ID.Eq(id)).Delete()
	return err
}

func toModelCategory(c *entity.Category) *model.Category {
	return &model.Category{
		ID:        c.ID,
		Name:      c.Name,
		Slug:      &c.Slug,
		ArticleCount: &c.ArticleCount,
	}
}

func toEntityCategory(mc *model.Category) *entity.Category {
	cat := &entity.Category{
		ID:   mc.ID,
		Name: mc.Name,
	}
	if mc.Slug != nil {
		cat.Slug = *mc.Slug
	}
	if mc.ArticleCount != nil {
		cat.ArticleCount = *mc.ArticleCount
	}
	if mc.CreatedAt != nil {
		cat.CreatedAt = *mc.CreatedAt
	}
	if mc.UpdatedAt != nil {
		cat.UpdatedAt = *mc.UpdatedAt
	}
	return cat
}
