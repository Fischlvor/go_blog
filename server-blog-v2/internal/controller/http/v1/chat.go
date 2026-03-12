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

// createSession 创建聊天会话。
func (v *V1) createSession(c fiber.Ctx) error {
	userUUID := middleware.GetUserUUID(c)
	if userUUID == "" {
		return shared.WriteError(c, http.StatusUnauthorized, bizcode.ErrorLoginRequired, "login required")
	}

	var req struct {
		Title string `json:"title"`
	}
	if err := c.Bind().JSON(&req); err != nil {
		return shared.WriteError(c, http.StatusBadRequest, bizcode.ErrorParam, "invalid request body")
	}

	session, err := v.aiChat.CreateSession(c.Context(), userUUID, input.CreateSession{
		Title: req.Title,
	})
	if err != nil {
		v.logger.Error(err, "http - v1 - chat - createSession")
		return shared.WriteError(c, http.StatusInternalServerError, bizcode.ErrorSessionCreateFail, "failed to create session")
	}

	return shared.WriteSuccess(c, shared.WithData(session))
}

// listSessions 会话列表。
func (v *V1) listSessions(c fiber.Ctx) error {
	userUUID := middleware.GetUserUUID(c)
	if userUUID == "" {
		return shared.WriteError(c, http.StatusUnauthorized, bizcode.ErrorLoginRequired, "login required")
	}

	pq := shared.ParsePageQuery(c)

	result, err := v.aiChat.ListSessions(c.Context(), userUUID, input.ListSessions{
		PageParams: input.PageParams{Page: pq.Page, PageSize: pq.PageSize},
	})
	if err != nil {
		v.logger.Error(err, "http - v1 - chat - listSessions")
		return shared.WriteError(c, http.StatusInternalServerError, bizcode.ErrorDatabase, "failed to list sessions")
	}

	return shared.WriteSuccess(c, shared.WithData(shared.NewPage(result.Items, result.Page, result.PageSize, result.Total)))
}

// getMessages 获取会话消息。
func (v *V1) getMessages(c fiber.Ctx) error {
	userUUID := middleware.GetUserUUID(c)
	if userUUID == "" {
		return shared.WriteError(c, http.StatusUnauthorized, bizcode.ErrorLoginRequired, "login required")
	}

	sessionID, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return shared.WriteError(c, http.StatusBadRequest, bizcode.ErrorParamFormat, "invalid session id")
	}

	messages, err := v.aiChat.GetMessages(c.Context(), sessionID, userUUID)
	if err != nil {
		v.logger.Error(err, "http - v1 - chat - getMessages")
		return shared.WriteError(c, http.StatusInternalServerError, bizcode.ErrorDatabase, "failed to get messages")
	}

	return shared.WriteSuccess(c, shared.WithData(messages))
}

// sendMessage 发送消息（SSE 流式响应）。
func (v *V1) sendMessage(c fiber.Ctx) error {
	userUUID := middleware.GetUserUUID(c)
	if userUUID == "" {
		return shared.WriteError(c, http.StatusUnauthorized, bizcode.ErrorLoginRequired, "login required")
	}

	sessionID, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return shared.WriteError(c, http.StatusBadRequest, bizcode.ErrorParamFormat, "invalid session id")
	}

	var req struct {
		Content string `json:"content"`
	}
	if err := c.Bind().JSON(&req); err != nil {
		return shared.WriteError(c, http.StatusBadRequest, bizcode.ErrorParam, "invalid request body")
	}

	if req.Content == "" {
		return shared.WriteError(c, http.StatusBadRequest, bizcode.ErrorParamMissing, "content is required")
	}

	// 获取流式响应通道
	stream, err := v.aiChat.SendMessage(c.Context(), sessionID, userUUID, input.SendMessage{
		Content: req.Content,
	})
	if err != nil {
		v.logger.Error(err, "http - v1 - chat - sendMessage")
		return shared.WriteError(c, http.StatusInternalServerError, bizcode.ErrorMessageSendFail, "failed to send message")
	}

	// 设置 SSE 响应头
	c.Set("Content-Type", "text/event-stream")
	c.Set("Cache-Control", "no-cache")
	c.Set("Connection", "keep-alive")

	// 流式输出
	for chunk := range stream {
		if _, err := c.Write([]byte("data: " + chunk + "\n\n")); err != nil {
			break
		}
	}

	return nil
}

// deleteSession 删除会话。
func (v *V1) deleteSession(c fiber.Ctx) error {
	userUUID := middleware.GetUserUUID(c)
	if userUUID == "" {
		return shared.WriteError(c, http.StatusUnauthorized, bizcode.ErrorLoginRequired, "login required")
	}

	sessionID, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return shared.WriteError(c, http.StatusBadRequest, bizcode.ErrorParamFormat, "invalid session id")
	}

	if err := v.aiChat.DeleteSession(c.Context(), sessionID, userUUID); err != nil {
		v.logger.Error(err, "http - v1 - chat - deleteSession")
		return shared.WriteError(c, http.StatusInternalServerError, bizcode.ErrorDatabase, "failed to delete session")
	}

	return shared.WriteSuccess(c)
}

// getAvailableModels 获取可用模型列表。
func (v *V1) getAvailableModels(c fiber.Ctx) error {
	models := []map[string]interface{}{
		{
			"name":         v.cfg.AI.Model,
			"display_name": "DeepSeek R1 (七牛云)",
			"provider":     "qiniu",
			"is_active":    true,
		},
	}
	return shared.WriteSuccess(c, shared.WithData(models))
}

// getSessionDetail 获取会话详情。
func (v *V1) getSessionDetail(c fiber.Ctx) error {
	userUUID := middleware.GetUserUUID(c)
	if userUUID == "" {
		return shared.WriteError(c, http.StatusUnauthorized, bizcode.ErrorLoginRequired, "login required")
	}

	sessionID, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return shared.WriteError(c, http.StatusBadRequest, bizcode.ErrorParamFormat, "invalid session id")
	}

	messages, err := v.aiChat.GetMessages(c.Context(), sessionID, userUUID)
	if err != nil {
		v.logger.Error(err, "http - v1 - chat - getSessionDetail")
		return shared.WriteError(c, http.StatusInternalServerError, bizcode.ErrorDatabase, "failed to get session detail")
	}

	return shared.WriteSuccess(c, shared.WithData(messages))
}

// sendMessageStream 流式发送消息（SSE）。
func (v *V1) sendMessageStream(c fiber.Ctx) error {
	// 复用 sendMessage 的实现
	return v.sendMessage(c)
}

// updateSession 更新会话。
func (v *V1) updateSession(c fiber.Ctx) error {
	userUUID := middleware.GetUserUUID(c)
	if userUUID == "" {
		return shared.WriteError(c, http.StatusUnauthorized, bizcode.ErrorLoginRequired, "login required")
	}

	var req struct {
		ID    int64  `json:"id"`
		Title string `json:"title"`
	}
	if err := c.Bind().JSON(&req); err != nil {
		return shared.WriteError(c, http.StatusBadRequest, bizcode.ErrorParam, "invalid request body")
	}

	// TODO: 实现会话更新逻辑
	return shared.WriteSuccess(c, shared.WithMsg("session updated"))
}
