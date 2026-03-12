package v1

import (
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v3"

	"server-blog-v2/internal/controller/http/bizcode"
	"server-blog-v2/internal/controller/http/middleware"
	"server-blog-v2/internal/controller/http/shared"
	"server-blog-v2/internal/usecase/input"
)

// listComments 评论列表。
func (v *V1) listComments(c fiber.Ctx) error {
	articleSlug := c.Params("article_slug")
	if articleSlug == "" {
		return shared.WriteError(c, http.StatusBadRequest, bizcode.ErrorParamMissing, "article_slug is required")
	}

	pq := shared.ParsePageQuery(c)

	result, err := v.comment.ListByArticleSlug(c.Context(), articleSlug, input.ListComments{
		PageParams: input.PageParams{Page: pq.Page, PageSize: pq.PageSize},
	})
	if err != nil {
		v.logger.Error(err, "http - v1 - comment - listComments")
		return shared.WriteError(c, http.StatusInternalServerError, bizcode.ErrorDatabase, "failed to list comments")
	}

	return shared.WriteSuccess(c, shared.WithData(shared.NewPage(result.Items, result.Page, result.PageSize, result.Total)))
}

// createComment 创建评论。
func (v *V1) createComment(c fiber.Ctx) error {
	userUUID := middleware.GetUserUUID(c)
	if userUUID == "" {
		return shared.WriteError(c, http.StatusUnauthorized, bizcode.ErrorLoginRequired, "login required")
	}

	var req struct {
		ArticleSlug string `json:"article_slug"`
		Content     string `json:"content"`
		ParentID    *int64 `json:"parent_id"`
	}

	if err := c.Bind().JSON(&req); err != nil {
		return shared.WriteError(c, http.StatusBadRequest, bizcode.ErrorParam, "invalid request body")
	}

	if req.Content == "" {
		return shared.WriteError(c, http.StatusBadRequest, bizcode.ErrorParamMissing, "content is required")
	}

	id, err := v.comment.Create(c.Context(), input.CreateComment{
		ArticleSlug: req.ArticleSlug,
		UserUUID:    userUUID,
		Content:     req.Content,
		ParentID:    req.ParentID,
	})
	if err != nil {
		v.logger.Error(err, "http - v1 - comment - createComment")
		return shared.WriteError(c, http.StatusInternalServerError, bizcode.ErrorDatabase, "failed to create comment")
	}

	return shared.WriteSuccess(c, shared.WithData(fiber.Map{"id": id}))
}

// deleteComment 删除评论。
func (v *V1) deleteComment(c fiber.Ctx) error {
	userUUID := middleware.GetUserUUID(c)
	if userUUID == "" {
		return shared.WriteError(c, http.StatusUnauthorized, bizcode.ErrorLoginRequired, "login required")
	}

	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return shared.WriteError(c, http.StatusBadRequest, bizcode.ErrorParamFormat, "invalid id")
	}

	if err := v.comment.Delete(c.Context(), id, userUUID); err != nil {
		v.logger.Error(err, "http - v1 - comment - deleteComment")
		return shared.WriteError(c, http.StatusInternalServerError, bizcode.ErrorDatabase, "failed to delete comment")
	}

	return shared.WriteSuccess(c)
}

// listNewComments 获取最新评论（公开，用于首页展示）。
func (v *V1) listNewComments(c fiber.Ctx) error {
	result, err := v.comment.ListAll(c.Context(), input.ListAllComments{
		PageParams: input.PageParams{Page: 1, PageSize: 10},
	})
	if err != nil {
		v.logger.Error(err, "http - v1 - comment - listNewComments")
		return shared.WriteError(c, http.StatusInternalServerError, bizcode.ErrorDatabase, "failed to list new comments")
	}

	return shared.WriteSuccess(c, shared.WithData(result.Items))
}
