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

// listCategories 分类列表。
// @Summary 分类列表（管理端）
// @Tags Admin.Category
// @Security BearerAuth
// @Produce json
// @Param page query int false "页码" default(1)
// @Param page_size query int false "分页大小" default(10)
// @Success 200 {object} shared.Envelope
// @Router /admin/category/list [get]
func (a *Admin) listCategories(c fiber.Ctx) error {
	pq := shared.ParsePageQuery(c)

	result, err := a.content.ListCategories(c.Context(), input.ListCategories{
		PageParams: input.PageParams{
			Page:     pq.Page,
			PageSize: pq.PageSize,
		},
	})

	if err != nil {
		a.logger.Error(err, "http - admin - category - listCategories")
		return shared.WriteError(c, http.StatusInternalServerError, bizcode.ErrorDatabase, "failed to list categories")
	}

	return shared.WriteSuccess(c, shared.WithData(shared.NewPage(result.Items, result.Page, result.PageSize, result.Total)))
}

// createCategory 创建分类。
// @Summary 创建分类（管理端）
// @Tags Admin.Category
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param body body request.CreateCategory true "分类信息"
// @Success 200 {object} shared.Envelope
// @Router /admin/category/create [post]
func (a *Admin) createCategory(c fiber.Ctx) error {
	var req request.CreateCategory
	if err := c.Bind().JSON(&req); err != nil {
		return shared.WriteError(c, http.StatusBadRequest, bizcode.ErrorParam, "invalid request body")
	}

	if err := a.validate.Struct(req); err != nil {
		return shared.WriteError(c, http.StatusBadRequest, bizcode.ErrorParamFormat, err.Error())
	}

	id, err := a.content.CreateCategory(c.Context(), input.CreateCategory{
		Name: req.Name,
		Slug: req.Slug,
	})

	if err != nil {
		a.logger.Error(err, "http - admin - category - createCategory")
		return shared.WriteError(c, http.StatusInternalServerError, bizcode.ErrorDatabase, "failed to create category")
	}

	return shared.WriteSuccess(c, shared.WithData(fiber.Map{"id": id}))
}

// updateCategory 更新分类。
// @Summary 更新分类（管理端）
// @Tags Admin.Category
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param body body request.UpdateCategory true "分类信息"
// @Success 200 {object} shared.Envelope
// @Router /admin/category/update [put]
func (a *Admin) updateCategory(c fiber.Ctx) error {
	var req request.UpdateCategory
	if err := c.Bind().JSON(&req); err != nil {
		return shared.WriteError(c, http.StatusBadRequest, bizcode.ErrorParam, "invalid request body")
	}

	if err := a.validate.Struct(req); err != nil {
		return shared.WriteError(c, http.StatusBadRequest, bizcode.ErrorParamFormat, err.Error())
	}

	err := a.content.UpdateCategory(c.Context(), input.UpdateCategory{
		ID:   req.ID,
		Name: req.Name,
		Slug: req.Slug,
	})

	if err != nil {
		a.logger.Error(err, "http - admin - category - updateCategory")
		return shared.WriteError(c, http.StatusInternalServerError, bizcode.ErrorDatabase, "failed to update category")
	}

	return shared.WriteSuccess(c)
}

// deleteCategory 删除分类。
// @Summary 删除分类（管理端）
// @Tags Admin.Category
// @Security BearerAuth
// @Produce json
// @Param id path int true "分类 ID"
// @Success 200 {object} shared.Envelope
// @Router /admin/category/delete/{id} [delete]
func (a *Admin) deleteCategory(c fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return shared.WriteError(c, http.StatusBadRequest, bizcode.ErrorParamFormat, "invalid category id")
	}

	if err := a.content.DeleteCategory(c.Context(), id); err != nil {
		a.logger.Error(err, "http - admin - category - deleteCategory")
		return shared.WriteError(c, http.StatusInternalServerError, bizcode.ErrorDatabase, "failed to delete category")
	}

	return shared.WriteSuccess(c)
}
