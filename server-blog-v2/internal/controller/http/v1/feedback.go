package v1

import (
	"net/http"

	"github.com/gofiber/fiber/v3"

	"server-blog-v2/internal/controller/http/bizcode"
	"server-blog-v2/internal/controller/http/middleware"
	"server-blog-v2/internal/controller/http/shared"
	"server-blog-v2/internal/usecase/input"
)

// createFeedback 创建反馈。
func (v *V1) createFeedback(c fiber.Ctx) error {
	var req struct {
		Type    string `json:"type"`
		Content string `json:"content"`
		Contact string `json:"contact"`
	}

	if err := c.Bind().JSON(&req); err != nil {
		return shared.WriteError(c, http.StatusBadRequest, bizcode.ErrorParam, "invalid request body")
	}

	if req.Content == "" {
		return shared.WriteError(c, http.StatusBadRequest, bizcode.ErrorParamMissing, "content is required")
	}

	if err := v.feedback.Create(c.Context(), input.CreateFeedback{
		Type:    req.Type,
		Content: req.Content,
		Contact: req.Contact,
	}); err != nil {
		v.logger.Error(err, "http - v1 - feedback - createFeedback")
		return shared.WriteError(c, http.StatusInternalServerError, bizcode.ErrorDatabase, "failed to create feedback")
	}

	return shared.WriteSuccess(c, shared.WithMsg("feedback submitted"))
}

// listNewFeedback 获取最新反馈（公开，用于首页展示已回复的反馈）。
func (v *V1) listNewFeedback(c fiber.Ctx) error {
	status := "replied"
	result, err := v.feedback.List(c.Context(), input.ListFeedback{
		PageParams: input.PageParams{Page: 1, PageSize: 10},
		Status:     &status,
	})
	if err != nil {
		v.logger.Error(err, "http - v1 - feedback - listNewFeedback")
		return shared.WriteError(c, http.StatusInternalServerError, bizcode.ErrorDatabase, "failed to list feedbacks")
	}

	return shared.WriteSuccess(c, shared.WithData(result.Items))
}

// getFeedbackInfo 获取反馈详情（需要登录）。
func (v *V1) getFeedbackInfo(c fiber.Ctx) error {
	// 用户查看自己的反馈列表
	userUUID := middleware.GetUserUUID(c)
	if userUUID == "" {
		return shared.WriteError(c, http.StatusUnauthorized, bizcode.ErrorLoginRequired, "login required")
	}

	// TODO: 实现用户反馈列表查询
	return shared.WriteSuccess(c, shared.WithData([]interface{}{}))
}
