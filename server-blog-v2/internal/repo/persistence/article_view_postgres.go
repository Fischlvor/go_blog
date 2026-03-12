package persistence

import (
	"context"

	"server-blog-v2/internal/entity"
	"server-blog-v2/internal/repo"
	"server-blog-v2/internal/repo/persistence/gen/model"
	"server-blog-v2/internal/repo/persistence/gen/query"

	"gorm.io/gorm"
)

type articleViewRepo struct {
	query *query.Query
}

// NewArticleViewRepo 创建文章浏览记录仓库。
func NewArticleViewRepo(db *gorm.DB) repo.ArticleViewRepo {
	return &articleViewRepo{query: query.Use(db)}
}

func (r *articleViewRepo) Record(ctx context.Context, view *entity.ArticleView) error {
	mv := toModelArticleView(view)
	return r.query.ArticleView.WithContext(ctx).Create(mv)
}

func (r *articleViewRepo) IncrViews(ctx context.Context, articleSlug string) error {
	a := r.query.Article
	_, err := a.WithContext(ctx).Where(a.Slug.Eq(articleSlug)).UpdateSimple(a.Views.Add(1))
	return err
}

func toModelArticleView(v *entity.ArticleView) *model.ArticleView {
	mv := &model.ArticleView{
		ID:          v.ID,
		ArticleSlug: &v.ArticleSlug,
		IPAddress:   v.IPAddress,
		ViewedAt:    &v.ViewedAt,
	}
	if v.UserAgent != nil {
		mv.UserAgent = v.UserAgent
	}
	if v.Referer != nil {
		mv.Referer = v.Referer
	}
	return mv
}
