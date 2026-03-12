package admin

import (
	"net/http"

	"github.com/gofiber/fiber/v3"

	"server-blog-v2/internal/controller/http/admin/request"
	"server-blog-v2/internal/controller/http/bizcode"
	"server-blog-v2/internal/controller/http/shared"
	"server-blog-v2/internal/usecase/input"
)

// listFeedbacks 反馈列表。
// @Summary 反馈列表（管理端）
// @Tags Admin.Feedback
// @Security BearerAuth
// @Produce json
// @Param page query int false "页码" default(1)
// @Param page_size query int false "分页大小" default(10)
// @Param filter.status query string false "状态" Enums(pending, replied)
// @Success 200 {object} shared.Envelope
// @Router /admin/feedback/list [get]
func (a *Admin) listFeedbacks(c fiber.Ctx) error {
	pq := shared.ParsePageQueryWithOptions(c, shared.WithAllowedFilters("status"))

	var status *string
	if s, ok := pq.Filters["status"]; ok && s != "" {
		status = &s
	}

	result, err := a.feedback.List(c.Context(), input.ListFeedback{
		PageParams: input.PageParams{Page: pq.Page, PageSize: pq.PageSize},
		Status:     status,
	})

	if err != nil {
		a.logger.Error(err, "http - admin - feedback - listFeedbacks")
		return shared.WriteError(c, http.StatusInternalServerError, bizcode.ErrorDatabase, "failed to list feedbacks")
	}

	return shared.WriteSuccess(c, shared.WithData(shared.NewPage(result.Items, result.Page, result.PageSize, result.Total)))
}

// deleteFeedback 批量删除反馈。
// @Summary 批量删除反馈（管理端）
// @Tags Admin.Feedback
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param body body request.BatchDelete true "反馈 ID 列表"
// @Success 200 {object} shared.Envelope
// @Router /admin/feedback/delete [delete]
func (a *Admin) deleteFeedback(c fiber.Ctx) error {
	var req request.BatchDelete
	if err := c.Bind().JSON(&req); err != nil {
		return shared.WriteError(c, http.StatusBadRequest, bizcode.ErrorParam, "invalid request body")
	}

	if err := a.validate.Struct(req); err != nil {
		return shared.WriteError(c, http.StatusBadRequest, bizcode.ErrorParamFormat, err.Error())
	}

	for _, id := range req.IDs {
		if err := a.feedback.Delete(c.Context(), id); err != nil {
			a.logger.Error(err, "http - admin - feedback - deleteFeedback", "id", id)
		}
	}

	return shared.WriteSuccess(c)
}

// replyFeedback 回复反馈。
// @Summary 回复反馈（管理端）
// @Tags Admin.Feedback
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param body body request.ReplyFeedback true "回复信息"
// @Success 200 {object} shared.Envelope
// @Router /admin/feedback/reply [put]
func (a *Admin) replyFeedback(c fiber.Ctx) error {
	var req request.ReplyFeedback
	if err := c.Bind().JSON(&req); err != nil {
		return shared.WriteError(c, http.StatusBadRequest, bizcode.ErrorParam, "invalid request body")
	}

	if err := a.validate.Struct(req); err != nil {
		return shared.WriteError(c, http.StatusBadRequest, bizcode.ErrorParamFormat, err.Error())
	}

	if err := a.feedback.Reply(c.Context(), req.ID, req.Reply); err != nil {
		a.logger.Error(err, "http - admin - feedback - replyFeedback")
		return shared.WriteError(c, http.StatusInternalServerError, bizcode.ErrorDatabase, "failed to reply feedback")
	}

	return shared.WriteSuccess(c)
}
