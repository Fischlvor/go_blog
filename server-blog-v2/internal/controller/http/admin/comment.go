package admin

import (
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v3"

	"server-blog-v2/internal/controller/http/bizcode"
	"server-blog-v2/internal/controller/http/shared"
	"server-blog-v2/internal/usecase/input"
)

// listComments 评论列表。
// @Summary 评论列表（管理端）
// @Tags Admin.Comment
// @Security BearerAuth
// @Produce json
// @Param page query int false "页码" default(1)
// @Param page_size query int false "分页大小" default(10)
// @Param keyword query string false "关键字"
// @Param filter.article_slug query string false "文章 Slug"
// @Success 200 {object} shared.Envelope
// @Router /admin/comment/list [get]
func (a *Admin) listComments(c fiber.Ctx) error {
	pq := shared.ParsePageQueryWithOptions(c, shared.WithAllowedFilters("article_slug", "user_uuid"))

	var keyword *input.KeywordParams
	if pq.Keyword != "" {
		keyword = &input.KeywordParams{Keyword: pq.Keyword}
	}

	var articleSlug *string
	if slug, ok := pq.Filters["article_slug"]; ok && slug != "" {
		articleSlug = &slug
	}

	var userUUID *string
	if uuid, ok := pq.Filters["user_uuid"]; ok && uuid != "" {
		userUUID = &uuid
	}

	result, err := a.comment.ListAll(c.Context(), input.ListAllComments{
		PageParams:  input.PageParams{Page: pq.Page, PageSize: pq.PageSize},
		Keyword:     keyword,
		ArticleSlug: articleSlug,
		UserUUID:    userUUID,
	})

	if err != nil {
		a.logger.Error(err, "http - admin - comment - listComments")
		return shared.WriteError(c, http.StatusInternalServerError, bizcode.ErrorDatabase, "failed to list comments")
	}

	return shared.WriteSuccess(c, shared.WithData(shared.NewPage(result.Items, result.Page, result.PageSize, result.Total)))
}

// deleteComment 删除评论。
// @Summary 删除评论（管理端）
// @Tags Admin.Comment
// @Security BearerAuth
// @Produce json
// @Param id path int true "评论 ID"
// @Success 200 {object} shared.Envelope
// @Router /admin/comment/delete/{id} [delete]
func (a *Admin) deleteComment(c fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return shared.WriteError(c, http.StatusBadRequest, bizcode.ErrorParamFormat, "invalid comment id")
	}

	if err := a.comment.AdminDelete(c.Context(), id); err != nil {
		a.logger.Error(err, "http - admin - comment - deleteComment")
		return shared.WriteError(c, http.StatusInternalServerError, bizcode.ErrorDatabase, "failed to delete comment")
	}

	return shared.WriteSuccess(c)
}
