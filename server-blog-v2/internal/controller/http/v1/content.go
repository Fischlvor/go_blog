package v1

import (
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v3"

	"server-blog-v2/internal/controller/http/middleware"
	"server-blog-v2/internal/controller/http/shared"
	"server-blog-v2/internal/controller/http/v1/response"
	"server-blog-v2/internal/usecase/input"
	"server-blog-v2/internal/usecase/output"
)

// listArticles 文章列表。
// @Summary 文章列表
// @Tags V1.Content
// @Produce json
// @Param page query int false "页码" default(1)
// @Param page_size query int false "分页大小" default(10)
// @Param keyword query string false "关键字"
// @Param filter.category_id query string false "分类 ID"
// @Param filter.tag_id query string false "标签 ID"
// @Success 200 {object} shared.Envelope{data=response.ArticleSummaryPage}
// @Router /v1/content/posts [get]
func (v *V1) listArticles(c fiber.Ctx) error {
	pq := shared.ParsePageQueryWithOptions(c, shared.WithAllowedFilters("category_id", "tag_id"))

	pageParams := input.PageParams{
		Page:     pq.Page,
		PageSize: pq.PageSize,
	}

	var keywordParams *input.KeywordParams
	if pq.Keyword != "" {
		keywordParams = &input.KeywordParams{Keyword: pq.Keyword}
	}

	var sortParams *input.SortParams
	if pq.SortBy != "" {
		sortParams = &input.SortParams{SortBy: pq.SortBy, Order: pq.Order}
	}

	var categoryID input.IntFilterParam
	if cid, ok := pq.Filters["category_id"]; ok && cid != "" {
		if id, err := strconv.Atoi(cid); err == nil {
			categoryID = &id
		}
	}

	var tagID input.IntFilterParam
	if tid, ok := pq.Filters["tag_id"]; ok && tid != "" {
		if id, err := strconv.Atoi(tid); err == nil {
			tagID = &id
		}
	}

	// 获取可选的用户 UUID
	userUUID := optionalUserUUID(c)

	result, err := v.content.ListPublicArticles(c.Context(), input.ListPublicArticles{
		PageParams: pageParams,
		Keyword:    keywordParams,
		Sort:       sortParams,
		CategoryID: categoryID,
		TagID:      tagID,
	}, userUUID)

	if err != nil {
		v.logger.Error(err, "http - v1 - content - listArticles")
		return shared.WriteError(c, http.StatusInternalServerError, response.ErrorListPostsFailed, "failed to list posts")
	}

	// 转换为响应格式
	list := make([]response.ArticleSummary, 0, len(result.Items))
	for _, p := range result.Items {
		list = append(list, toArticleSummaryResponse(p))
	}

	return shared.WriteSuccess(c, shared.WithData(shared.NewPage(list, result.Page, result.PageSize, result.Total)))
}

// getArticle 文章详情。
// @Summary 文章详情
// @Tags V1.Content
// @Produce json
// @Param slug path string true "文章 Slug"
// @Success 200 {object} shared.Envelope{data=response.ArticleDetail}
// @Router /v1/content/posts/{slug} [get]
func (v *V1) getArticle(c fiber.Ctx) error {
	slug := c.Params("slug")
	if slug == "" {
		return shared.WriteError(c, http.StatusBadRequest, response.ErrorParamMissing, "missing slug")
	}

	userUUID := optionalUserUUID(c)

	post, err := v.content.GetPublicArticleBySlug(c.Context(), slug, userUUID)
	if err != nil {
		v.logger.Error(err, "http - v1 - content - getArticle")
		return shared.WriteError(c, http.StatusNotFound, response.ErrorPostNotFound, "post not found")
	}

	// 记录浏览
	v.content.RecordView(c.Context(), slug, c.IP(), c.Get("User-Agent"), c.Get("Referer"))

	return shared.WriteSuccess(c, shared.WithData(toArticleDetailResponse(post)))
}

// toggleArticleLike 点赞/取消点赞。
// @Summary 点赞/取消点赞
// @Tags V1.Content
// @Security BearerAuth
// @Param slug path string true "文章 Slug"
// @Success 200 {object} shared.Envelope{data=response.LikeInfo}
// @Router /v1/content/posts/{slug}/likes [post]
func (v *V1) toggleArticleLike(c fiber.Ctx) error {
	userUUID := middleware.GetUserUUID(c)
	if userUUID == "" {
		return shared.WriteError(c, http.StatusUnauthorized, response.ErrorLoginRequired, "login required")
	}

	slug := c.Params("slug")
	if slug == "" {
		return shared.WriteError(c, http.StatusBadRequest, response.ErrorParamMissing, "missing slug")
	}

	liked, likes, err := v.content.ToggleLikeOnArticle(c.Context(), slug, userUUID)
	if err != nil {
		v.logger.Error(err, "http - v1 - content - toggleArticleLike")
		return shared.WriteError(c, http.StatusInternalServerError, response.ErrorLikePostFailed, "like post failed")
	}

	return shared.WriteSuccess(c, shared.WithData(response.LikeInfo{Liked: &liked, Likes: likes}))
}

// removeArticleLike 取消点赞。
// @Summary 取消点赞
// @Tags V1.Content
// @Security BearerAuth
// @Param slug path string true "文章 Slug"
// @Success 200 {object} shared.Envelope{data=response.LikeInfo}
// @Router /v1/content/posts/{slug}/likes [delete]
func (v *V1) removeArticleLike(c fiber.Ctx) error {
	userUUID := middleware.GetUserUUID(c)
	if userUUID == "" {
		return shared.WriteError(c, http.StatusUnauthorized, response.ErrorLoginRequired, "login required")
	}

	slug := c.Params("slug")
	if slug == "" {
		return shared.WriteError(c, http.StatusBadRequest, response.ErrorParamMissing, "missing slug")
	}

	_, likes, err := v.content.RemoveLikeOnArticle(c.Context(), slug, userUUID)
	if err != nil {
		v.logger.Error(err, "http - v1 - content - removeArticleLike")
		return shared.WriteError(c, http.StatusInternalServerError, response.ErrorUnlikePostFailed, "unlike post failed")
	}

	return shared.WriteSuccess(c, shared.WithData(response.LikeInfo{Liked: nil, Likes: likes}))
}

// listCategories 分类列表。
func (v *V1) listCategories(c fiber.Ctx) error {
	result, err := v.content.GetAllPublicCategories(c.Context())
	if err != nil {
		v.logger.Error(err, "http - v1 - content - listCategories")
		return shared.WriteError(c, http.StatusInternalServerError, response.ErrorListCategoriesFailed, "failed to list categories")
	}

	list := make([]response.CategoryDetail, 0, len(result.Items))
	for _, cat := range result.Items {
		list = append(list, response.CategoryDetail{
			ID:        cat.ID,
			Name:      cat.Name,
			Slug:      cat.Slug,
			ArticleCount: cat.ArticleCount,
			CreatedAt: cat.CreatedAt,
			UpdatedAt: cat.UpdatedAt,
		})
	}

	return shared.WriteSuccess(c, shared.WithData(list))
}

// listTags 标签列表。
func (v *V1) listTags(c fiber.Ctx) error {
	result, err := v.content.GetAllPublicTags(c.Context())
	if err != nil {
		v.logger.Error(err, "http - v1 - content - listTags")
		return shared.WriteError(c, http.StatusInternalServerError, response.ErrorListTagsFailed, "failed to list tags")
	}

	list := make([]response.TagDetail, 0, len(result.Items))
	for _, tag := range result.Items {
		list = append(list, response.TagDetail{
			ID:        tag.ID,
			Name:      tag.Name,
			Slug:      tag.Slug,
			ArticleCount: tag.ArticleCount,
			CreatedAt: tag.CreatedAt,
			UpdatedAt: tag.UpdatedAt,
		})
	}

	return shared.WriteSuccess(c, shared.WithData(list))
}

// ==================== 辅助函数 ====================

func toArticleSummaryResponse(p output.ArticleSummary) response.ArticleSummary {
	tags := make([]response.BaseTag, len(p.Tags))
	for i, t := range p.Tags {
		tags[i] = response.BaseTag{ID: t.ID, Name: t.Name, Slug: t.Slug}
	}
	return response.ArticleSummary{
		ID:            p.ID,
		Title:         p.Title,
		Slug:          p.Slug,
		Excerpt:       p.Excerpt,
		FeaturedImage: p.FeaturedImage,
		AuthorUUID:    p.AuthorUUID,
		Author:        response.AuthorInfo{UUID: p.Author.UUID, Nickname: p.Author.Nickname, Avatar: p.Author.Avatar},
		Status:        p.Status,
		Views:         p.Views,
		Like:          response.LikeInfo{Liked: p.Like.Liked, Likes: p.Like.Likes},
		IsFeatured:    p.IsFeatured,
		PublishedAt:   p.PublishedAt,
		CreatedAt:     p.CreatedAt,
		UpdatedAt:     p.UpdatedAt,
		Category:      response.BaseCategory{ID: p.Category.ID, Name: p.Category.Name, Slug: p.Category.Slug},
		Tags:          tags,
	}
}

func toArticleDetailResponse(p *output.ArticleDetail) response.ArticleDetail {
	tags := make([]response.BaseTag, len(p.Tags))
	for i, t := range p.Tags {
		tags[i] = response.BaseTag{ID: t.ID, Name: t.Name, Slug: t.Slug}
	}
	return response.ArticleDetail{
		ID:              p.ID,
		Title:           p.Title,
		Slug:            p.Slug,
		Excerpt:         p.Excerpt,
		FeaturedImage:   p.FeaturedImage,
		AuthorUUID:      p.AuthorUUID,
		Author:          response.AuthorInfo{UUID: p.Author.UUID, Nickname: p.Author.Nickname, Avatar: p.Author.Avatar},
		Status:          p.Status,
		Views:           p.Views,
		Like:            response.LikeInfo{Liked: p.Like.Liked, Likes: p.Like.Likes},
		IsFeatured:      p.IsFeatured,
		PublishedAt:     p.PublishedAt,
		CreatedAt:       p.CreatedAt,
		UpdatedAt:       p.UpdatedAt,
		Category:        response.BaseCategory{ID: p.Category.ID, Name: p.Category.Name, Slug: p.Category.Slug},
		Tags:            tags,
		Content:         p.Content,
		MetaTitle:       p.MetaTitle,
		MetaDescription: p.MetaDescription,
	}
}

// listUserLikedArticles 获取用户点赞的文章列表。
func (v *V1) listUserLikedArticles(c fiber.Ctx) error {
	userUUID := middleware.GetUserUUID(c)
	if userUUID == "" {
		return shared.WriteError(c, http.StatusUnauthorized, response.ErrorLoginRequired, "login required")
	}

	pq := shared.ParsePageQuery(c)

	result, err := v.content.ListUserLikedArticles(c.Context(), userUUID, input.ListUserLikedArticles{
		PageParams: input.PageParams{Page: pq.Page, PageSize: pq.PageSize},
	})
	if err != nil {
		v.logger.Error(err, "http - v1 - content - listUserLikedArticles")
		return shared.WriteError(c, http.StatusInternalServerError, response.ErrorListPostsFailed, "failed to list liked articles")
	}

	list := make([]response.ArticleSummary, 0, len(result.Items))
	for _, p := range result.Items {
		list = append(list, toArticleSummaryResponse(p))
	}

	return shared.WriteSuccess(c, shared.WithData(shared.NewPage(list, result.Page, result.PageSize, result.Total)))
}
