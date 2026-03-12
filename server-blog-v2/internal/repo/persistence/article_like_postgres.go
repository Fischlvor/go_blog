package persistence

import (
	"context"

	"server-blog-v2/internal/entity"
	"server-blog-v2/internal/repo"
	"server-blog-v2/internal/repo/persistence/gen/model"
	"server-blog-v2/internal/repo/persistence/gen/query"

	"gorm.io/gorm"
)

type articleLikeRepo struct {
	query *query.Query
}

// NewArticleLikeRepo 创建文章点赞仓库。
func NewArticleLikeRepo(db *gorm.DB) repo.ArticleLikeRepo {
	return &articleLikeRepo{query: query.Use(db)}
}

func (r *articleLikeRepo) HasLiked(ctx context.Context, articleSlug, userUUID string) (bool, error) {
	al := r.query.ArticleLike
	count, err := al.WithContext(ctx).Where(al.ArticleSlug.Eq(articleSlug), al.UserUUID.Eq(userUUID)).Count()
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *articleLikeRepo) Toggle(ctx context.Context, articleSlug, userUUID string) (liked bool, count int32, err error) {
	err = r.query.Transaction(func(tx *query.Query) error {
		al := tx.ArticleLike
		a := tx.Article

		// 检查是否已点赞
		existing, findErr := al.WithContext(ctx).Where(al.ArticleSlug.Eq(articleSlug), al.UserUUID.Eq(userUUID)).First()
		if findErr != nil && findErr != gorm.ErrRecordNotFound {
			return findErr
		}

		if existing != nil {
			// 已点赞，取消
			if _, delErr := al.WithContext(ctx).Where(al.ID.Eq(existing.ID)).Delete(); delErr != nil {
				return delErr
			}
			// 减少点赞数
			if _, updateErr := a.WithContext(ctx).Where(a.Slug.Eq(articleSlug)).UpdateSimple(a.Likes.Sub(1)); updateErr != nil {
				return updateErr
			}
			liked = false
		} else {
			// 未点赞，添加
			newLike := &model.ArticleLike{ArticleSlug: &articleSlug, UserUUID: &userUUID}
			if createErr := al.WithContext(ctx).Create(newLike); createErr != nil {
				return createErr
			}
			// 增加点赞数
			if _, updateErr := a.WithContext(ctx).Where(a.Slug.Eq(articleSlug)).UpdateSimple(a.Likes.Add(1)); updateErr != nil {
				return updateErr
			}
			liked = true
		}

		// 获取最新点赞数
		article, getErr := a.WithContext(ctx).Where(a.Slug.Eq(articleSlug)).First()
		if getErr != nil {
			return getErr
		}
		if article.Likes != nil {
			count = *article.Likes
		}
		return nil
	})
	return
}

func (r *articleLikeRepo) Remove(ctx context.Context, articleSlug, userUUID string) (removed bool, count int32, err error) {
	err = r.query.Transaction(func(tx *query.Query) error {
		al := tx.ArticleLike
		a := tx.Article

		// 检查是否已点赞
		existing, findErr := al.WithContext(ctx).Where(al.ArticleSlug.Eq(articleSlug), al.UserUUID.Eq(userUUID)).First()
		if findErr != nil {
			if findErr == gorm.ErrRecordNotFound {
				removed = false
				// 获取当前点赞数
				article, getErr := a.WithContext(ctx).Where(a.Slug.Eq(articleSlug)).First()
				if getErr != nil {
					return getErr
				}
				if article.Likes != nil {
					count = *article.Likes
				}
				return nil
			}
			return findErr
		}

		// 删除点赞
		if _, delErr := al.WithContext(ctx).Where(al.ID.Eq(existing.ID)).Delete(); delErr != nil {
			return delErr
		}
		// 减少点赞数
		if _, updateErr := a.WithContext(ctx).Where(a.Slug.Eq(articleSlug)).UpdateSimple(a.Likes.Sub(1)); updateErr != nil {
			return updateErr
		}
		removed = true

		// 获取最新点赞数
		article, getErr := a.WithContext(ctx).Where(a.Slug.Eq(articleSlug)).First()
		if getErr != nil {
			return getErr
		}
		if article.Likes != nil {
			count = *article.Likes
		}
		return nil
	})
	return
}

func (r *articleLikeRepo) ListUserLikedArticles(ctx context.Context, userUUID string, offset, limit int) ([]*entity.Article, int64, error) {
	al := r.query.ArticleLike
	a := r.query.Article

	// 获取用户点赞的文章 slug 列表
	likes, err := al.WithContext(ctx).Where(al.UserUUID.Eq(userUUID)).Order(al.CreatedAt.Desc()).Find()
	if err != nil {
		return nil, 0, err
	}

	if len(likes) == 0 {
		return []*entity.Article{}, 0, nil
	}

	slugs := make([]string, len(likes))
	for i, like := range likes {
		if like.ArticleSlug != nil {
			slugs[i] = *like.ArticleSlug
		}
	}

	total := int64(len(slugs))

	// 分页获取文章
	start := offset
	end := offset + limit
	if start >= len(slugs) {
		return []*entity.Article{}, total, nil
	}
	if end > len(slugs) {
		end = len(slugs)
	}
	pagedSlugs := slugs[start:end]

	articles, err := a.WithContext(ctx).Where(a.Slug.In(pagedSlugs...)).Find()
	if err != nil {
		return nil, 0, err
	}

	result := make([]*entity.Article, len(articles))
	for i, ma := range articles {
		result[i] = toEntityArticle(ma)
	}

	return result, total, nil
}
