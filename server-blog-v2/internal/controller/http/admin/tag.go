package admin

import (
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v3"

	"server-blog-v2/internal/controller/http/admin/request"
	"server-blog-v2/internal/controller/http/bizcode"
	"server-blog-v2/internal/controller/http/shared"
	"server-blog-v2/internal/usecase/input"
)

// listTags 标签列表。
// @Summary 标签列表（管理端）
// @Tags Admin.Tag
// @Security BearerAuth
// @Produce json
// @Param page query int false "页码" default(1)
// @Param page_size query int false "分页大小" default(10)
// @Success 200 {object} shared.Envelope
// @Router /admin/tag/list [get]
func (a *Admin) listTags(c fiber.Ctx) error {
	pq := shared.ParsePageQuery(c)

	result, err := a.content.ListTags(c.Context(), input.ListTags{
		PageParams: input.PageParams{
			Page:     pq.Page,
			PageSize: pq.PageSize,
		},
	})

	if err != nil {
		a.logger.Error(err, "http - admin - tag - listTags")
		return shared.WriteError(c, http.StatusInternalServerError, bizcode.ErrorDatabase, "failed to list tags")
	}

	return shared.WriteSuccess(c, shared.WithData(shared.NewPage(result.Items, result.Page, result.PageSize, result.Total)))
}

// createTag 创建标签。
// @Summary 创建标签（管理端）
// @Tags Admin.Tag
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param body body request.CreateTag true "标签信息"
// @Success 200 {object} shared.Envelope
// @Router /admin/tag/create [post]
func (a *Admin) createTag(c fiber.Ctx) error {
	var req request.CreateTag
	if err := c.Bind().JSON(&req); err != nil {
		return shared.WriteError(c, http.StatusBadRequest, bizcode.ErrorParam, "invalid request body")
	}

	if err := a.validate.Struct(req); err != nil {
		return shared.WriteError(c, http.StatusBadRequest, bizcode.ErrorParamFormat, err.Error())
	}

	id, err := a.content.CreateTag(c.Context(), input.CreateTag{
		Name: req.Name,
		Slug: req.Slug,
	})

	if err != nil {
		a.logger.Error(err, "http - admin - tag - createTag")
		return shared.WriteError(c, http.StatusInternalServerError, bizcode.ErrorDatabase, "failed to create tag")
	}

	return shared.WriteSuccess(c, shared.WithData(fiber.Map{"id": id}))
}

// updateTag 更新标签。
// @Summary 更新标签（管理端）
// @Tags Admin.Tag
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param body body request.UpdateTag true "标签信息"
// @Success 200 {object} shared.Envelope
// @Router /admin/tag/update [put]
func (a *Admin) updateTag(c fiber.Ctx) error {
	var req request.UpdateTag
	if err := c.Bind().JSON(&req); err != nil {
		return shared.WriteError(c, http.StatusBadRequest, bizcode.ErrorParam, "invalid request body")
	}

	if err := a.validate.Struct(req); err != nil {
		return shared.WriteError(c, http.StatusBadRequest, bizcode.ErrorParamFormat, err.Error())
	}

	err := a.content.UpdateTag(c.Context(), input.UpdateTag{
		ID:   req.ID,
		Name: req.Name,
		Slug: req.Slug,
	})

	if err != nil {
		a.logger.Error(err, "http - admin - tag - updateTag")
		return shared.WriteError(c, http.StatusInternalServerError, bizcode.ErrorDatabase, "failed to update tag")
	}

	return shared.WriteSuccess(c)
}

// deleteTag 删除标签。
// @Summary 删除标签（管理端）
// @Tags Admin.Tag
// @Security BearerAuth
// @Produce json
// @Param id path int true "标签 ID"
// @Success 200 {object} shared.Envelope
// @Router /admin/tag/delete/{id} [delete]
func (a *Admin) deleteTag(c fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return shared.WriteError(c, http.StatusBadRequest, bizcode.ErrorParamFormat, "invalid tag id")
	}

	if err := a.content.DeleteTag(c.Context(), id); err != nil {
		a.logger.Error(err, "http - admin - tag - deleteTag")
		return shared.WriteError(c, http.StatusInternalServerError, bizcode.ErrorDatabase, "failed to delete tag")
	}

	return shared.WriteSuccess(c)
}
