package content

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"time"

	"server-blog-v2/config"
	"server-blog-v2/internal/entity"
	"server-blog-v2/internal/repo"
	"server-blog-v2/internal/usecase"
	"server-blog-v2/internal/usecase/input"
	"server-blog-v2/internal/usecase/output"
	"server-blog-v2/internal/usecase/urlutil"
)

var (
	ErrRepo     = errors.New("repo")
	ErrNotFound = errors.New("not found")
)

// generateSlug 生成随机 slug（15字节 = 20字符，与 ES _id 格式一致）
func generateSlug() string {
	b := make([]byte, 15)
	rand.Read(b)
	return base64.RawURLEncoding.EncodeToString(b)
}

type useCase struct {
	cfg          *config.Config
	articles     repo.ArticleRepo
	tags         repo.TagRepo
	categories   repo.CategoryRepo
	articleLikes repo.ArticleLikeRepo
	articleViews repo.ArticleViewRepo
	users        repo.UserRepo
}

// New 创建 Content UseCase。
func New(
	cfg *config.Config,
	articles repo.ArticleRepo,
	tags repo.TagRepo,
	categories repo.CategoryRepo,
	articleLikes repo.ArticleLikeRepo,
	articleViews repo.ArticleViewRepo,
	users repo.UserRepo,
) usecase.Content {
	return &useCase{
		cfg:          cfg,
		articles:     articles,
		tags:         tags,
		categories:   categories,
		articleLikes: articleLikes,
		articleViews: articleViews,
		users:        users,
	}
}

// ==================== 文章 - 管理端 ====================

func (u *useCase) ListArticles(ctx context.Context, params input.ListArticles) (*output.ListResult[output.ArticleSummary], error) {
	offset := (params.Page - 1) * params.PageSize

	var keyword *string
	if params.Keyword != nil {
		keyword = &params.Keyword.Keyword
	}
	var sortBy, order *string
	if params.Sort != nil {
		sortBy = &params.Sort.SortBy
		order = &params.Sort.Order
	}
	var categoryID, tagID *int
	if params.CategoryID != nil {
		categoryID = (*int)(params.CategoryID)
	}
	if params.TagID != nil {
		tagID = (*int)(params.TagID)
	}
	var status *string
	if params.Status != nil {
		status = (*string)(params.Status)
	}
	var visibility *string
	if params.Visibility != nil {
		visibility = params.Visibility
	}

	articles, total, err := u.articles.List(ctx, offset, params.PageSize, keyword, sortBy, order, categoryID, tagID, status, visibility)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrRepo, err)
	}

	items, err := u.toArticleSummaries(ctx, articles, nil)
	if err != nil {
		return nil, err
	}

	return &output.ListResult[output.ArticleSummary]{
		Items:    items,
		Page:     params.Page,
		PageSize: params.PageSize,
		Total:    total,
	}, nil
}

func (u *useCase) GetArticleByID(ctx context.Context, id int64) (*output.ArticleDetail, error) {
	article, err := u.articles.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrRepo, err)
	}
	return u.toArticleDetail(ctx, article, nil)
}

func (u *useCase) CreateArticle(ctx context.Context, params input.CreateArticle) (int64, error) {
	// 默认可见性为 public
	visibility := params.Visibility
	if visibility == "" {
		visibility = entity.ArticleVisibilityPublic
	}

	// 如果 slug 为空，自动生成
	slug := params.Slug
	if slug == "" {
		slug = generateSlug()
	}

	article := &entity.Article{
		Title:      params.Title,
		Slug:       slug,
		Excerpt:    params.Excerpt,
		Content:    params.Content,
		AuthorUUID: params.AuthorUUID,
		CategoryID: params.CategoryID,
		TagIDs:     params.TagIDs, // 直接存储标签 ID 数组
		Status:     params.Status,
		Visibility: visibility,
		IsFeatured: params.IsFeatured,
	}
	if params.FeaturedImage != nil {
		article.FeaturedImage = params.FeaturedImage
	}
	if params.Status == entity.ArticleStatusPublished {
		now := time.Now()
		article.PublishedAt = &now
	}

	id, err := u.articles.Create(ctx, article)
	if err != nil {
		return 0, fmt.Errorf("%w: %v", ErrRepo, err)
	}

	return id, nil
}

func (u *useCase) UpdateArticle(ctx context.Context, params input.UpdateArticle) error {
	// 默认可见性为 public
	visibility := params.Visibility
	if visibility == "" {
		visibility = entity.ArticleVisibilityPublic
	}

	// 通过 slug 获取现有文章（用于检查是否首次发布）
	existing, err := u.articles.GetBySlug(ctx, params.Slug)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrRepo, err)
	}

	article := &entity.Article{
		Title:      params.Title,
		Slug:       params.Slug,
		Excerpt:    params.Excerpt,
		AuthorUUID: params.AuthorUUID,
		CategoryID: params.CategoryID,
		TagIDs:     params.TagIDs, // 直接存储标签 ID 数组
		Status:     params.Status,
		Visibility: visibility,
		IsFeatured: params.IsFeatured,
	}
	// 只在 content 非空时更新
	if params.Content != "" {
		article.Content = params.Content
	}
	if params.FeaturedImage != nil {
		article.FeaturedImage = params.FeaturedImage
	}

	// 检查是否首次发布
	if existing.PublishedAt == nil && params.Status == entity.ArticleStatusPublished {
		now := time.Now()
		article.PublishedAt = &now
	}

	// 使用 slug 更新文章，只在 content 非空时更新内容
	includeContent := params.Content != ""
	if err := u.articles.UpdateBySlug(ctx, params.Slug, article, includeContent); err != nil {
		return fmt.Errorf("%w: %v", ErrRepo, err)
	}

	return nil
}

func (u *useCase) DeleteArticle(ctx context.Context, id int64) error {
	if err := u.articles.Delete(ctx, id); err != nil {
		return fmt.Errorf("%w: %v", ErrRepo, err)
	}
	return nil
}

// ==================== 文章 - 公开端 ====================

func (u *useCase) ListPublicArticles(ctx context.Context, params input.ListPublicArticles, userUUID *string) (*output.ListResult[output.ArticleSummary], error) {
	offset := (params.Page - 1) * params.PageSize

	var keyword *string
	if params.Keyword != nil {
		keyword = &params.Keyword.Keyword
	}
	var sortBy, order *string
	if params.Sort != nil {
		sortBy = &params.Sort.SortBy
		order = &params.Sort.Order
	}
	var categoryID, tagID *int
	if params.CategoryID != nil {
		categoryID = (*int)(params.CategoryID)
	}
	if params.TagID != nil {
		tagID = (*int)(params.TagID)
	}

	published := entity.ArticleStatusPublished
	public := entity.ArticleVisibilityPublic
	articles, total, err := u.articles.List(ctx, offset, params.PageSize, keyword, sortBy, order, categoryID, tagID, &published, &public)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrRepo, err)
	}

	items, err := u.toArticleSummaries(ctx, articles, userUUID)
	if err != nil {
		return nil, err
	}

	return &output.ListResult[output.ArticleSummary]{
		Items:    items,
		Page:     params.Page,
		PageSize: params.PageSize,
		Total:    total,
	}, nil
}

func (u *useCase) GetPublicArticleBySlug(ctx context.Context, slug string, userUUID *string) (*output.ArticleDetail, error) {
	article, err := u.articles.GetBySlug(ctx, slug)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrRepo, err)
	}
	return u.toArticleDetail(ctx, article, userUUID)
}

func (u *useCase) RecordView(ctx context.Context, articleSlug string, ip, userAgent, referer string) {
	var ua, ref *string
	if userAgent != "" {
		ua = &userAgent
	}
	if referer != "" {
		ref = &referer
	}
	_ = u.articleViews.Record(ctx, &entity.ArticleView{
		ArticleSlug: articleSlug,
		IPAddress:   ip,
		UserAgent:   ua,
		Referer:     ref,
		ViewedAt:    time.Now(),
	})
	_ = u.articleViews.IncrViews(ctx, articleSlug)
}

// ==================== 点赞 ====================

func (u *useCase) ToggleLikeOnArticle(ctx context.Context, articleSlug, userUUID string) (bool, int32, error) {
	liked, count, err := u.articleLikes.Toggle(ctx, articleSlug, userUUID)
	if err != nil {
		return false, 0, fmt.Errorf("%w: %v", ErrRepo, err)
	}
	return liked, count, nil
}

func (u *useCase) RemoveLikeOnArticle(ctx context.Context, articleSlug, userUUID string) (bool, int32, error) {
	removed, count, err := u.articleLikes.Remove(ctx, articleSlug, userUUID)
	if err != nil {
		return false, 0, fmt.Errorf("%w: %v", ErrRepo, err)
	}
	return removed, count, nil
}

func (u *useCase) ListUserLikedArticles(ctx context.Context, userUUID string, params input.ListUserLikedArticles) (*output.ListResult[output.ArticleSummary], error) {
	offset := (params.Page - 1) * params.PageSize

	articles, total, err := u.articleLikes.ListUserLikedArticles(ctx, userUUID, offset, params.PageSize)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrRepo, err)
	}

	items, err := u.toArticleSummaries(ctx, articles, &userUUID)
	if err != nil {
		return nil, err
	}

	return &output.ListResult[output.ArticleSummary]{
		Items:    items,
		Page:     params.Page,
		PageSize: params.PageSize,
		Total:    total,
	}, nil
}

// ==================== 分类 ====================

func (u *useCase) ListCategories(ctx context.Context, params input.ListCategories) (*output.ListResult[output.CategoryDetail], error) {
	offset := (params.Page - 1) * params.PageSize

	var keyword *string
	if params.Keyword != nil {
		keyword = &params.Keyword.Keyword
	}
	var sortBy, order *string
	if params.Sort != nil {
		sortBy = &params.Sort.SortBy
		order = &params.Sort.Order
	}

	categories, total, err := u.categories.List(ctx, offset, params.PageSize, keyword, sortBy, order)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrRepo, err)
	}

	items := make([]output.CategoryDetail, len(categories))
	for i, c := range categories {
		items[i] = toCategoryDetail(c)
	}

	return &output.ListResult[output.CategoryDetail]{
		Items:    items,
		Page:     params.Page,
		PageSize: params.PageSize,
		Total:    total,
	}, nil
}

func (u *useCase) GetAllPublicCategories(ctx context.Context) (*output.AllResult[output.CategoryDetail], error) {
	categories, err := u.categories.ListAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrRepo, err)
	}
	items := make([]output.CategoryDetail, len(categories))
	for i, c := range categories {
		items[i] = toCategoryDetail(c)
	}
	return &output.AllResult[output.CategoryDetail]{
		Items: items,
		Total: int64(len(items)),
	}, nil
}

func (u *useCase) CreateCategory(ctx context.Context, params input.CreateCategory) (int64, error) {
	id, err := u.categories.Create(ctx, entity.Category{Name: params.Name, Slug: params.Slug})
	if err != nil {
		return 0, fmt.Errorf("%w: %v", ErrRepo, err)
	}
	return id, nil
}

func (u *useCase) UpdateCategory(ctx context.Context, params input.UpdateCategory) error {
	if err := u.categories.Update(ctx, entity.Category{ID: params.ID, Name: params.Name, Slug: params.Slug}); err != nil {
		return fmt.Errorf("%w: %v", ErrRepo, err)
	}
	return nil
}

func (u *useCase) DeleteCategory(ctx context.Context, id int64) error {
	if err := u.categories.Delete(ctx, id); err != nil {
		return fmt.Errorf("%w: %v", ErrRepo, err)
	}
	return nil
}

// ==================== 标签 ====================

func (u *useCase) ListTags(ctx context.Context, params input.ListTags) (*output.ListResult[output.TagDetail], error) {
	offset := (params.Page - 1) * params.PageSize

	var keyword *string
	if params.Keyword != nil {
		keyword = &params.Keyword.Keyword
	}
	var sortBy, order *string
	if params.Sort != nil {
		sortBy = &params.Sort.SortBy
		order = &params.Sort.Order
	}

	tags, total, err := u.tags.List(ctx, offset, params.PageSize, keyword, sortBy, order)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrRepo, err)
	}

	items := make([]output.TagDetail, len(tags))
	for i, t := range tags {
		items[i] = toTagDetail(t)
	}

	return &output.ListResult[output.TagDetail]{
		Items:    items,
		Page:     params.Page,
		PageSize: params.PageSize,
		Total:    total,
	}, nil
}

func (u *useCase) GetAllPublicTags(ctx context.Context) (*output.AllResult[output.TagDetail], error) {
	tags, err := u.tags.ListAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrRepo, err)
	}
	items := make([]output.TagDetail, len(tags))
	for i, t := range tags {
		items[i] = toTagDetail(t)
	}
	return &output.AllResult[output.TagDetail]{
		Items: items,
		Total: int64(len(items)),
	}, nil
}

func (u *useCase) CreateTag(ctx context.Context, params input.CreateTag) (int64, error) {
	id, err := u.tags.Create(ctx, entity.Tag{Name: params.Name, Slug: params.Slug})
	if err != nil {
		return 0, fmt.Errorf("%w: %v", ErrRepo, err)
	}
	return id, nil
}

func (u *useCase) UpdateTag(ctx context.Context, params input.UpdateTag) error {
	if err := u.tags.Update(ctx, entity.Tag{ID: params.ID, Name: params.Name, Slug: params.Slug}); err != nil {
		return fmt.Errorf("%w: %v", ErrRepo, err)
	}
	return nil
}

func (u *useCase) DeleteTag(ctx context.Context, id int64) error {
	if err := u.tags.Delete(ctx, id); err != nil {
		return fmt.Errorf("%w: %v", ErrRepo, err)
	}
	return nil
}

// ==================== 辅助函数 ====================

func (u *useCase) toArticleSummaries(ctx context.Context, articles []*entity.Article, userUUID *string) ([]output.ArticleSummary, error) {
	items := make([]output.ArticleSummary, len(articles))
	for i, a := range articles {
		author := u.getAuthorInfo(ctx, a.AuthorUUID)
		like := u.getLikeInfo(ctx, a.Slug, a.Likes, userUUID)
		cat, _ := u.categories.GetByID(ctx, a.CategoryID)
		tags, _ := u.tags.ListByIDs(ctx, a.TagIDs) // 使用 tag_ids 数组
		items[i] = u.toArticleSummary(a, author, like, cat, tags)
	}
	return items, nil
}

func (u *useCase) toArticleDetail(ctx context.Context, a *entity.Article, userUUID *string) (*output.ArticleDetail, error) {
	author := u.getAuthorInfo(ctx, a.AuthorUUID)
	like := u.getLikeInfo(ctx, a.Slug, a.Likes, userUUID)
	cat, _ := u.categories.GetByID(ctx, a.CategoryID)
	tags, _ := u.tags.ListByIDs(ctx, a.TagIDs) // 使用 tag_ids 数组

	detail := &output.ArticleDetail{
		BaseArticle: u.toBaseArticle(a),
		Author:      author,
		Like:        like,
		Category:    output.BaseCategory{ID: cat.ID, Name: cat.Name, Slug: cat.Slug},
		Tags:        toBaseTags(tags),
		Content:     a.Content,
	}
	if a.MetaTitle != nil {
		detail.MetaTitle = *a.MetaTitle
	}
	if a.MetaDescription != nil {
		detail.MetaDescription = *a.MetaDescription
	}
	return detail, nil
}

func (u *useCase) getAuthorInfo(ctx context.Context, authorUUID string) output.AuthorInfo {
	user, err := u.users.GetByUUID(ctx, authorUUID)
	if err != nil || user == nil {
		return output.AuthorInfo{UUID: authorUUID}
	}
	return output.AuthorInfo{
		UUID:     user.UUID,
		Nickname: user.Nickname,
		Avatar:   urlutil.ResolveImageURL(u.cfg, user.Avatar),
	}
}

func (u *useCase) getLikeInfo(ctx context.Context, articleSlug string, likes int32, userUUID *string) output.LikeInfo {
	if userUUID == nil {
		return output.LikeInfo{Likes: likes}
	}
	liked, err := u.articleLikes.HasLiked(ctx, articleSlug, *userUUID)
	if err != nil {
		return output.LikeInfo{Likes: likes}
	}
	return output.LikeInfo{Liked: &liked, Likes: likes}
}

func (u *useCase) toBaseArticle(a *entity.Article) output.BaseArticle {
	bp := output.BaseArticle{
		ID:          a.ID,
		Title:       a.Title,
		Slug:        a.Slug,
		AuthorUUID:  a.AuthorUUID,
		Status:      a.Status,
		Visibility:  a.Visibility,
		Views:       a.Views,
		IsFeatured:  a.IsFeatured,
		PublishedAt: a.PublishedAt,
		CreatedAt:   a.CreatedAt,
		UpdatedAt:   a.UpdatedAt,
	}
	if a.Excerpt != nil {
		bp.Excerpt = *a.Excerpt
	}
	if a.FeaturedImage != nil {
		bp.FeaturedImage = urlutil.ResolveImageURL(u.cfg, *a.FeaturedImage)
	}
	if a.ReadTime != nil {
		bp.ReadTime = *a.ReadTime
	}
	return bp
}

func (u *useCase) toArticleSummary(a *entity.Article, author output.AuthorInfo, like output.LikeInfo, cat *entity.Category, tags []*entity.Tag) output.ArticleSummary {
	return output.ArticleSummary{
		BaseArticle: u.toBaseArticle(a),
		Author:      author,
		Like:        like,
		Category:    output.BaseCategory{ID: cat.ID, Name: cat.Name, Slug: cat.Slug},
		Tags:        toBaseTags(tags),
	}
}

func toBaseTags(tags []*entity.Tag) []output.BaseTag {
	bt := make([]output.BaseTag, len(tags))
	for i, t := range tags {
		bt[i] = output.BaseTag{ID: t.ID, Name: t.Name, Slug: t.Slug}
	}
	return bt
}

func toCategoryDetail(c *entity.Category) output.CategoryDetail {
	return output.CategoryDetail{
		BaseCategory: output.BaseCategory{ID: c.ID, Name: c.Name, Slug: c.Slug},
		ArticleCount: c.ArticleCount,
		CreatedAt:    c.CreatedAt,
		UpdatedAt:    c.UpdatedAt,
	}
}

func toTagDetail(t *entity.Tag) output.TagDetail {
	return output.TagDetail{
		BaseTag:      output.BaseTag{ID: t.ID, Name: t.Name, Slug: t.Slug},
		ArticleCount: t.ArticleCount,
		CreatedAt:    t.CreatedAt,
		UpdatedAt:    t.UpdatedAt,
	}
}
