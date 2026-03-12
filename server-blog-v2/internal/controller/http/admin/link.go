package admin

import (
	"net/http"

	"github.com/gofiber/fiber/v3"

	"server-blog-v2/internal/controller/http/admin/request"
	"server-blog-v2/internal/controller/http/bizcode"
	"server-blog-v2/internal/controller/http/shared"
	"server-blog-v2/internal/usecase/input"
)

// listLinks 友链列表。
// @Summary 友链列表（管理端）
// @Tags Admin.Link
// @Security BearerAuth
// @Produce json
// @Success 200 {object} shared.Envelope
// @Router /admin/friendLink/list [get]
func (a *Admin) listLinks(c fiber.Ctx) error {
	result, err := a.link.List(c.Context())
	if err != nil {
		a.logger.Error(err, "http - admin - link - listLinks")
		return shared.WriteError(c, http.StatusInternalServerError, bizcode.ErrorDatabase, "failed to list links")
	}

	return shared.WriteSuccess(c, shared.WithData(result))
}

// createLink 创建友链。
// @Summary 创建友链（管理端）
// @Tags Admin.Link
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param body body request.CreateLink true "友链信息"
// @Success 200 {object} shared.Envelope
// @Router /admin/friendLink/create [post]
func (a *Admin) createLink(c fiber.Ctx) error {
	var req request.CreateLink
	if err := c.Bind().JSON(&req); err != nil {
		return shared.WriteError(c, http.StatusBadRequest, bizcode.ErrorParam, "invalid request body")
	}

	if err := a.validate.Struct(req); err != nil {
		return shared.WriteError(c, http.StatusBadRequest, bizcode.ErrorParamFormat, err.Error())
	}

	id, err := a.link.Create(c.Context(), input.CreateLink{
		Name:        req.Name,
		URL:         req.URL,
		Logo:        req.Logo,
		Description: req.Description,
		Sort:        req.Sort,
		IsVisible:   true,
	})

	if err != nil {
		a.logger.Error(err, "http - admin - link - createLink")
		return shared.WriteError(c, http.StatusInternalServerError, bizcode.ErrorDatabase, "failed to create link")
	}

	return shared.WriteSuccess(c, shared.WithData(fiber.Map{"id": id}))
}

// updateLink 更新友链。
// @Summary 更新友链（管理端）
// @Tags Admin.Link
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param body body request.UpdateLink true "友链信息"
// @Success 200 {object} shared.Envelope
// @Router /admin/friendLink/update [put]
func (a *Admin) updateLink(c fiber.Ctx) error {
	var req request.UpdateLink
	if err := c.Bind().JSON(&req); err != nil {
		return shared.WriteError(c, http.StatusBadRequest, bizcode.ErrorParam, "invalid request body")
	}

	if err := a.validate.Struct(req); err != nil {
		return shared.WriteError(c, http.StatusBadRequest, bizcode.ErrorParamFormat, err.Error())
	}

	err := a.link.Update(c.Context(), input.UpdateLink{
		ID:          req.ID,
		Name:        req.Name,
		URL:         req.URL,
		Logo:        req.Logo,
		Description: req.Description,
		Sort:        req.Sort,
		IsVisible:   true,
	})

	if err != nil {
		a.logger.Error(err, "http - admin - link - updateLink")
		return shared.WriteError(c, http.StatusInternalServerError, bizcode.ErrorDatabase, "failed to update link")
	}

	return shared.WriteSuccess(c)
}

// deleteLink 批量删除友链。
// @Summary 批量删除友链（管理端）
// @Tags Admin.Link
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param body body request.BatchDelete true "友链 ID 列表"
// @Success 200 {object} shared.Envelope
// @Router /admin/friendLink/delete [delete]
func (a *Admin) deleteLink(c fiber.Ctx) error {
	var req request.BatchDelete
	if err := c.Bind().JSON(&req); err != nil {
		return shared.WriteError(c, http.StatusBadRequest, bizcode.ErrorParam, "invalid request body")
	}

	if err := a.validate.Struct(req); err != nil {
		return shared.WriteError(c, http.StatusBadRequest, bizcode.ErrorParamFormat, err.Error())
	}

	for _, id := range req.IDs {
		if err := a.link.Delete(c.Context(), id); err != nil {
			a.logger.Error(err, "http - admin - link - deleteLink", "id", id)
		}
	}

	return shared.WriteSuccess(c)
}
