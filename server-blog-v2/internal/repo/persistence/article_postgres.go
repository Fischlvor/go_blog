package persistence

import (
	"context"

	"server-blog-v2/internal/entity"
	"server-blog-v2/internal/repo"
	"server-blog-v2/internal/repo/persistence/gen/model"
	"server-blog-v2/internal/repo/persistence/gen/query"

	"gorm.io/gen/field"
	"gorm.io/gorm"
)

type articleRepo struct {
	query *query.Query
}

// NewArticleRepo 创建文章仓库。
func NewArticleRepo(db *gorm.DB) repo.ArticleRepo {
	return &articleRepo{query: query.Use(db)}
}

func (r *articleRepo) List(ctx context.Context, offset, limit int, keyword *string, sortBy, order *string, categoryID, tagID *int, status, visibility *string) ([]*entity.Article, int64, error) {
	a := r.query.Article
	do := a.WithContext(ctx)

	if keyword != nil && *keyword != "" {
		kw := "%" + *keyword + "%"
		do = do.Where(a.Title.Like(kw))
	}
	if categoryID != nil {
		do = do.Where(a.CategoryID.Eq(int64(*categoryID)))
	}
	if status != nil && *status != "" {
		do = do.Where(a.Status.Eq(*status))
	}
	if visibility != nil && *visibility != "" {
		do = do.Where(a.Visibility.Eq(*visibility))
	}

	total, err := do.Count()
	if err != nil {
		return nil, 0, err
	}

	do = do.Order(a.CreatedAt.Desc())
	rows, err := do.Offset(offset).Limit(limit).Find()
	if err != nil {
		return nil, 0, err
	}

	articles := make([]*entity.Article, len(rows))
	for i, row := range rows {
		articles[i] = toEntityArticle(row)
	}
	return articles, total, nil
}

func (r *articleRepo) GetByID(ctx context.Context, id int64) (*entity.Article, error) {
	a := r.query.Article
	row, err := a.WithContext(ctx).Where(a.ID.Eq(id)).First()
	if err != nil {
		return nil, err
	}
	return toEntityArticle(row), nil
}

func (r *articleRepo) Create(ctx context.Context, article *entity.Article) (int64, error) {
	ma := toModelArticle(article)
	if err := r.query.Article.WithContext(ctx).Create(ma); err != nil {
		return 0, err
	}
	return ma.ID, nil
}

func (r *articleRepo) Update(ctx context.Context, article *entity.Article) error {
	a := r.query.Article
	ma := toModelArticle(article)
	_, err := a.WithContext(ctx).Where(a.ID.Eq(article.ID)).Updates(ma)
	return err
}

func (r *articleRepo) UpdateBySlug(ctx context.Context, slug string, article *entity.Article, includeContent bool) error {
	a := r.query.Article
	ma := toModelArticle(article)

	// 基础更新字段
	fields := []field.Expr{a.Title, a.CategoryID, a.TagIds, a.Status, a.Visibility, a.IsFeatured}

	// 可选字段
	if includeContent {
		fields = append(fields, a.Content)
	}
	if article.Excerpt != nil {
		fields = append(fields, a.Excerpt)
	}
	if article.FeaturedImage != nil {
		fields = append(fields, a.FeaturedImage)
	}
	if article.PublishedAt != nil {
		fields = append(fields, a.PublishedAt)
	}

	_, err := a.WithContext(ctx).Where(a.Slug.Eq(slug)).Select(fields...).Updates(ma)
	return err
}

func (r *articleRepo) Delete(ctx context.Context, id int64) error {
	a := r.query.Article
	_, err := a.WithContext(ctx).Where(a.ID.Eq(id)).Delete()
	return err
}

func (r *articleRepo) GetBySlug(ctx context.Context, slug string) (*entity.Article, error) {
	a := r.query.Article
	ma, err := a.WithContext(ctx).Where(a.Slug.Eq(slug)).First()
	if err != nil {
		return nil, err
	}
	return toEntityArticle(ma), nil
}

func toModelArticle(a *entity.Article) *model.Article {
	ma := &model.Article{
		ID:              a.ID,
		Title:           a.Title,
		Slug:            a.Slug,
		Content:         a.Content,
		AuthorUUID:      &a.AuthorUUID,
		CategoryID:      a.CategoryID,
		TagIDs:          a.TagIDs, // 标签 ID 数组
		Status:          &a.Status,
		Visibility:      a.Visibility, // 非指针类型
		Views:           &a.Views,
		Likes:           &a.Likes,
		IsFeatured:      &a.IsFeatured,
		PublishedAt:     a.PublishedAt,
		MetaTitle:       a.MetaTitle,
		MetaDescription: a.MetaDescription,
	}
	if a.Excerpt != nil {
		ma.Excerpt = a.Excerpt
	}
	if a.FeaturedImage != nil {
		ma.FeaturedImage = a.FeaturedImage
	}
	if a.ReadTime != nil {
		ma.ReadTime = a.ReadTime
	}
	return ma
}

func toEntityArticle(ma *model.Article) *entity.Article {
	a := &entity.Article{
		ID:              ma.ID,
		Title:           ma.Title,
		Slug:            ma.Slug,
		Content:         ma.Content,
		FeaturedImage:   ma.FeaturedImage,
		CategoryID:      ma.CategoryID,
		MetaTitle:       ma.MetaTitle,
		MetaDescription: ma.MetaDescription,
		PublishedAt:     ma.PublishedAt,
	}
	// 处理指针类型字段
	if ma.Status != nil {
		a.Status = *ma.Status
	}
	a.Visibility = ma.Visibility // 非指针类型
	if ma.Views != nil {
		a.Views = *ma.Views
	}
	if ma.Likes != nil {
		a.Likes = *ma.Likes
	}
	if ma.IsFeatured != nil {
		a.IsFeatured = *ma.IsFeatured
	}
	if ma.CreatedAt != nil {
		a.CreatedAt = *ma.CreatedAt
	}
	if ma.UpdatedAt != nil {
		a.UpdatedAt = *ma.UpdatedAt
	}
	if ma.Excerpt != nil {
		a.Excerpt = ma.Excerpt
	}
	if ma.ReadTime != nil {
		a.ReadTime = ma.ReadTime
	}
	if ma.AuthorUUID != nil {
		a.AuthorUUID = *ma.AuthorUUID
	}
	// 标签 ID 数组
	if ma.TagIDs != nil {
		a.TagIDs = ma.TagIDs
	}
	return a
}
