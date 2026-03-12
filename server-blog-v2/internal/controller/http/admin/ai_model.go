package admin

import (
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v3"

	"server-blog-v2/internal/controller/http/bizcode"
	"server-blog-v2/internal/controller/http/shared"
	"server-blog-v2/internal/usecase/input"
)

// listAIModels AI 模型列表。
// @Summary AI 模型列表（管理端）
// @Tags Admin.AIManagement
// @Security BearerAuth
// @Produce json
// @Param page query int false "页码"
// @Param page_size query int false "每页数量"
// @Param name query string false "模型名称"
// @Param provider query string false "提供商"
// @Success 200 {object} shared.Envelope
// @Router /admin/ai-management/model [get]
func (a *Admin) listAIModels(c fiber.Ctx) error {
	pq := shared.ParsePageQueryWithOptions(c, shared.WithAllowedFilters("name", "provider"))

	result, err := a.aiModel.List(c.Context(), input.ListAIModels{
		PageParams: input.PageParams{Page: pq.Page, PageSize: pq.PageSize},
		Name:       pq.Filters["name"],
		Provider:   pq.Filters["provider"],
	})
	if err != nil {
		a.logger.Error(err, "http - admin - ai_model - listAIModels")
		return shared.WriteError(c, http.StatusInternalServerError, bizcode.ErrorDatabase, "failed to list ai models")
	}

	return shared.WriteSuccess(c, shared.WithData(shared.NewPage(result.Items, result.Page, result.PageSize, result.Total)))
}

// getAIModel 获取 AI 模型详情。
// @Summary 获取 AI 模型详情（管理端）
// @Tags Admin.AIManagement
// @Security BearerAuth
// @Produce json
// @Param id path int true "模型 ID"
// @Success 200 {object} shared.Envelope
// @Router /admin/ai-management/model/{id} [get]
func (a *Admin) getAIModel(c fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return shared.WriteError(c, http.StatusBadRequest, bizcode.ErrorParam, "invalid id")
	}

	result, err := a.aiModel.GetByID(c.Context(), id)
	if err != nil {
		a.logger.Error(err, "http - admin - ai_model - getAIModel")
		return shared.WriteError(c, http.StatusInternalServerError, bizcode.ErrorDatabase, "failed to get ai model")
	}

	return shared.WriteSuccess(c, shared.WithData(result))
}

// createAIModel 创建 AI 模型。
// @Summary 创建 AI 模型（管理端）
// @Tags Admin.AIManagement
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param body body input.CreateAIModel true "模型信息"
// @Success 200 {object} shared.Envelope
// @Router /admin/ai-management/model [post]
func (a *Admin) createAIModel(c fiber.Ctx) error {
	var req input.CreateAIModel
	if err := c.Bind().JSON(&req); err != nil {
		return shared.WriteError(c, http.StatusBadRequest, bizcode.ErrorParam, "invalid request body")
	}

	id, err := a.aiModel.Create(c.Context(), req)
	if err != nil {
		a.logger.Error(err, "http - admin - ai_model - createAIModel")
		return shared.WriteError(c, http.StatusInternalServerError, bizcode.ErrorDatabase, "failed to create ai model")
	}

	return shared.WriteSuccess(c, shared.WithData(map[string]int64{"id": id}))
}

// updateAIModel 更新 AI 模型。
// @Summary 更新 AI 模型（管理端）
// @Tags Admin.AIManagement
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param body body input.UpdateAIModel true "模型信息"
// @Success 200 {object} shared.Envelope
// @Router /admin/ai-management/model [put]
func (a *Admin) updateAIModel(c fiber.Ctx) error {
	var req input.UpdateAIModel
	if err := c.Bind().JSON(&req); err != nil {
		return shared.WriteError(c, http.StatusBadRequest, bizcode.ErrorParam, "invalid request body")
	}

	if err := a.aiModel.Update(c.Context(), req); err != nil {
		a.logger.Error(err, "http - admin - ai_model - updateAIModel")
		return shared.WriteError(c, http.StatusInternalServerError, bizcode.ErrorDatabase, "failed to update ai model")
	}

	return shared.WriteSuccess(c)
}

// deleteAIModel 删除 AI 模型。
// @Summary 删除 AI 模型（管理端）
// @Tags Admin.AIManagement
// @Security BearerAuth
// @Produce json
// @Param id path int true "模型 ID"
// @Success 200 {object} shared.Envelope
// @Router /admin/ai-management/model/{id} [delete]
func (a *Admin) deleteAIModel(c fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return shared.WriteError(c, http.StatusBadRequest, bizcode.ErrorParam, "invalid id")
	}

	if err := a.aiModel.Delete(c.Context(), id); err != nil {
		a.logger.Error(err, "http - admin - ai_model - deleteAIModel")
		return shared.WriteError(c, http.StatusInternalServerError, bizcode.ErrorDatabase, "failed to delete ai model")
	}

	return shared.WriteSuccess(c)
}
