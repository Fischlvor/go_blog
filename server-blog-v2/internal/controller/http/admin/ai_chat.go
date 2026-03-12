package admin

import (
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v3"

	"server-blog-v2/internal/controller/http/bizcode"
	"server-blog-v2/internal/controller/http/shared"
	"server-blog-v2/internal/usecase/input"
)

// listAISessions AI 会话列表。
// @Summary AI 会话列表（管理端）
// @Tags Admin.AIManagement
// @Security BearerAuth
// @Produce json
// @Param page query int false "页码" default(1)
// @Param page_size query int false "分页大小" default(10)
// @Param keyword query string false "关键字"
// @Param filter.user_uuid query string false "用户 UUID"
// @Success 200 {object} shared.Envelope
// @Router /admin/ai-management/session [get]
func (a *Admin) listAISessions(c fiber.Ctx) error {
	pq := shared.ParsePageQueryWithOptions(c, shared.WithAllowedFilters("user_uuid"))

	var keyword *input.KeywordParams
	if pq.Keyword != "" {
		keyword = &input.KeywordParams{Keyword: pq.Keyword}
	}

	var userUUID *string
	if uuid, ok := pq.Filters["user_uuid"]; ok && uuid != "" {
		userUUID = &uuid
	}

	result, err := a.aiChat.ListAllSessions(c.Context(), input.ListAllSessions{
		PageParams: input.PageParams{Page: pq.Page, PageSize: pq.PageSize},
		Keyword:    keyword,
		UserUUID:   userUUID,
	})

	if err != nil {
		a.logger.Error(err, "http - admin - ai_chat - listAISessions")
		return shared.WriteError(c, http.StatusInternalServerError, bizcode.ErrorDatabase, "failed to list sessions")
	}

	return shared.WriteSuccess(c, shared.WithData(shared.NewPage(result.Items, result.Page, result.PageSize, result.Total)))
}

// getAISession 获取 AI 会话详情。
// @Summary 获取 AI 会话详情（管理端）
// @Tags Admin.AIManagement
// @Security BearerAuth
// @Produce json
// @Param id path int true "会话 ID"
// @Success 200 {object} shared.Envelope
// @Router /admin/ai-management/session/{id} [get]
func (a *Admin) getAISession(c fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return shared.WriteError(c, http.StatusBadRequest, bizcode.ErrorParamFormat, "invalid session id")
	}

	session, err := a.aiChat.GetSessionByID(c.Context(), id)
	if err != nil {
		a.logger.Error(err, "http - admin - ai_chat - getAISession")
		return shared.WriteError(c, http.StatusNotFound, bizcode.ErrorNotFound, "session not found")
	}

	return shared.WriteSuccess(c, shared.WithData(session))
}

// deleteAISession 删除 AI 会话。
// @Summary 删除 AI 会话（管理端）
// @Tags Admin.AIManagement
// @Security BearerAuth
// @Produce json
// @Param id path int true "会话 ID"
// @Success 200 {object} shared.Envelope
// @Router /admin/ai-management/session/{id} [delete]
func (a *Admin) deleteAISession(c fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return shared.WriteError(c, http.StatusBadRequest, bizcode.ErrorParamFormat, "invalid session id")
	}

	if err := a.aiChat.AdminDeleteSession(c.Context(), id); err != nil {
		a.logger.Error(err, "http - admin - ai_chat - deleteAISession")
		return shared.WriteError(c, http.StatusInternalServerError, bizcode.ErrorDatabase, "failed to delete session")
	}

	return shared.WriteSuccess(c)
}

// listAIMessages AI 消息列表。
// @Summary AI 消息列表（管理端）
// @Tags Admin.AIManagement
// @Security BearerAuth
// @Produce json
// @Param page query int false "页码" default(1)
// @Param page_size query int false "分页大小" default(10)
// @Param session_id query int false "会话 ID"
// @Param role query string false "角色"
// @Success 200 {object} shared.Envelope
// @Router /admin/ai-management/message [get]
func (a *Admin) listAIMessages(c fiber.Ctx) error {
	pq := shared.ParsePageQueryWithOptions(c, shared.WithAllowedFilters("session_id", "role"))

	var sessionID *int64
	if sid, ok := pq.Filters["session_id"]; ok && sid != "" {
		if id, err := strconv.ParseInt(sid, 10, 64); err == nil {
			sessionID = &id
		}
	}

	var role *string
	if r, ok := pq.Filters["role"]; ok && r != "" {
		role = &r
	}

	result, err := a.aiChat.ListAllMessages(c.Context(), input.ListAllMessages{
		PageParams: input.PageParams{Page: pq.Page, PageSize: pq.PageSize},
		SessionID:  sessionID,
		Role:       role,
	})

	if err != nil {
		a.logger.Error(err, "http - admin - ai_chat - listAIMessages")
		return shared.WriteError(c, http.StatusInternalServerError, bizcode.ErrorDatabase, "failed to list messages")
	}

	return shared.WriteSuccess(c, shared.WithData(shared.NewPage(result.Items, result.Page, result.PageSize, result.Total)))
}
