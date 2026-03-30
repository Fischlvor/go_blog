package admin

import (
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v3"

	"server-blog-v2/internal/controller/http/admin/request"
	"server-blog-v2/internal/controller/http/bizcode"
	"server-blog-v2/internal/controller/http/middleware"
	"server-blog-v2/internal/controller/http/shared"
	"server-blog-v2/internal/usecase/input"
)

// listArticles 文章列表。
// @Summary 文章列表（管理端）
// @Tags Admin.Article
// @Security BearerAuth
// @Produce json
// @Param page query int false "页码" default(1)
// @Param page_size query int false "分页大小" default(10)
// @Param keyword query string false "关键字"
// @Success 200 {object} shared.Envelope
// @Router /admin/article/list [get]
func (a *Admin) listArticles(c fiber.Ctx) error {
	pq := shared.ParsePageQueryWithOptions(c, shared.WithAllowedFilters("category_id", "tag_id", "status", "visibility"))

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

	var status *string
	if s, ok := pq.Filters["status"]; ok && s != "" {
		status = &s
	}

	var visibility *string
	if v, ok := pq.Filters["visibility"]; ok && v != "" {
		visibility = &v
	}

	result, err := a.content.ListArticles(c.Context(), input.ListArticles{
		PageParams: pageParams,
		Keyword:    keywordParams,
		Sort:       sortParams,
		CategoryID: categoryID,
		TagID:      tagID,
		Status:     status,
		Visibility: visibility,
	})

	if err != nil {
		a.logger.Error(err, "http - admin - article - listArticles")
		return shared.WriteError(c, http.StatusInternalServerError, bizcode.ErrorDatabase, "failed to list articles")
	}

	return shared.WriteSuccess(c, shared.WithData(shared.NewPage(result.Items, result.Page, result.PageSize, result.Total)))
}

// getArticle 获取文章详情。
// @Summary 获取文章详情（管理端）
// @Tags Admin.Article
// @Security BearerAuth
// @Produce json
// @Param slug path string true "文章 Slug"
// @Success 200 {object} shared.Envelope
// @Router /admin/article/{slug} [get]
func (a *Admin) getArticle(c fiber.Ctx) error {
	slug := c.Params("slug")
	if slug == "" {
		return shared.WriteError(c, http.StatusBadRequest, bizcode.ErrorParamFormat, "invalid article slug")
	}

	article, err := a.content.GetArticleBySlug(c.Context(), slug)
	if err != nil {
		a.logger.Error(err, "http - admin - article - getArticle")
		return shared.WriteError(c, http.StatusNotFound, bizcode.ErrorNotFound, "article not found")
	}

	return shared.WriteSuccess(c, shared.WithData(article))
}

// createArticle 创建文章。
// @Summary 创建文章（管理端）
// @Tags Admin.Article
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param body body request.CreateArticle true "文章信息"
// @Success 200 {object} shared.Envelope
// @Router /admin/article/create [post]
func (a *Admin) createArticle(c fiber.Ctx) error {
	var req request.CreateArticle
	if err := c.Bind().JSON(&req); err != nil {
		return shared.WriteError(c, http.StatusBadRequest, bizcode.ErrorParam, "invalid request body")
	}

	if err := a.validate.Struct(req); err != nil {
		return shared.WriteError(c, http.StatusBadRequest, bizcode.ErrorParamFormat, shared.TranslateValidationErrors(err))
	}

	// draft 模式下 categoryId 可选（可以为 0），published 模式下必填
	if req.Status == "published" && req.CategoryID == 0 {
		return shared.WriteError(c, http.StatusBadRequest, bizcode.ErrorParamFormat, "category_id is required for published articles")
	}

	// 获取当前用户 UUID
	userUUID := middleware.GetUserUUID(c)
	if userUUID == "" {
		return shared.WriteError(c, http.StatusUnauthorized, bizcode.ErrorUnauthorized, "unauthorized")
	}

	// 处理可选字段
	var excerpt *string
	if req.Excerpt != "" {
		excerpt = &req.Excerpt
	}
	var featuredImage *string
	if req.FeaturedImage != "" {
		featuredImage = &req.FeaturedImage
	}

	// draft 模式下 categoryId 为 0 时转换为 NULL
	categoryID := req.CategoryID
	if req.Status == "draft" && categoryID == 0 {
		categoryID = 0 // 保持为 0，后端处理时转为 NULL
	}

	slug, err := a.content.CreateArticle(c.Context(), input.CreateArticle{
		Title:         req.Title,
		Slug:          req.Slug,
		Content:       req.Content,
		Excerpt:       excerpt,
		FeaturedImage: featuredImage,
		AuthorUUID:    userUUID,
		CategoryID:    categoryID,
		TagIDs:        req.TagIDs,
		Status:        req.Status,
		Visibility:    req.Visibility,
		IsFeatured:    req.IsFeatured,
	})

	if err != nil {
		a.logger.Error(err, "http - admin - article - createArticle")
		return shared.WriteError(c, http.StatusInternalServerError, bizcode.ErrorDatabase, "failed to create article")
	}

	return shared.WriteSuccess(c, shared.WithData(fiber.Map{"slug": slug}))
}

// updateArticle 更新文章。
// @Summary 更新文章（管理端）
// @Tags Admin.Article
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param body body request.UpdateArticle true "文章信息"
// @Success 200 {object} shared.Envelope
// @Router /admin/article/update [put]
func (a *Admin) updateArticle(c fiber.Ctx) error {
	var req request.UpdateArticle
	if err := c.Bind().JSON(&req); err != nil {
		return shared.WriteError(c, http.StatusBadRequest, bizcode.ErrorParam, "invalid request body")
	}

	if err := a.validate.Struct(req); err != nil {
		return shared.WriteError(c, http.StatusBadRequest, bizcode.ErrorParamFormat, shared.TranslateValidationErrors(err))
	}

	// draft 模式下 categoryId 可选（可以为 0），published 模式下必填
	if req.Status == "published" && req.CategoryID == 0 {
		return shared.WriteError(c, http.StatusBadRequest, bizcode.ErrorParamFormat, "category_id is required for published articles")
	}

	// 处理可选字段
	var excerpt *string
	if req.Excerpt != "" {
		excerpt = &req.Excerpt
	}
	var featuredImage *string
	if req.FeaturedImage != "" {
		featuredImage = &req.FeaturedImage
	}

	// draft 模式下 categoryId 为 0 时转换为 NULL
	categoryID := req.CategoryID
	if req.Status == "draft" && categoryID == 0 {
		categoryID = 0 // 保持为 0，后端处理时转为 NULL
	}

	err := a.content.UpdateArticle(c.Context(), input.UpdateArticle{
		Slug:          req.Slug,
		Title:         req.Title,
		Content:       req.Content,
		Excerpt:       excerpt,
		FeaturedImage: featuredImage,
		CategoryID:    categoryID,
		TagIDs:        req.TagIDs,
		Status:        req.Status,
		Visibility:    req.Visibility,
		IsFeatured:    req.IsFeatured,
	})

	if err != nil {
		a.logger.Error(err, "http - admin - article - updateArticle")
		return shared.WriteError(c, http.StatusInternalServerError, bizcode.ErrorDatabase, "failed to update article")
	}

	return shared.WriteSuccess(c)
}

// deleteArticle 批量删除文章。
// @Summary 批量删除文章（管理端）
// @Tags Admin.Article
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param body body request.BatchDeleteStr true "文章 Slug 列表"
// @Success 200 {object} shared.Envelope
// @Router /admin/article/delete [delete]
func (a *Admin) deleteArticle(c fiber.Ctx) error {
	var req request.BatchDeleteStr
	if err := c.Bind().JSON(&req); err != nil {
		return shared.WriteError(c, http.StatusBadRequest, bizcode.ErrorParam, "invalid request body")
	}

	if err := a.validate.Struct(req); err != nil {
		return shared.WriteError(c, http.StatusBadRequest, bizcode.ErrorParamFormat, err.Error())
	}

	for _, slug := range req.IDs {
		if err := a.content.DeleteArticleBySlug(c.Context(), slug); err != nil {
			a.logger.Error(err, "http - admin - article - deleteArticle", "slug", slug)
		}
	}

	return shared.WriteSuccess(c)
}
