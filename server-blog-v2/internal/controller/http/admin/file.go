package admin

import (
	"net/http"

	"github.com/gofiber/fiber/v3"

	"server-blog-v2/internal/controller/http/bizcode"
	"server-blog-v2/internal/controller/http/shared"
	"server-blog-v2/internal/usecase/input"
)

// uploadFile 上传文件。
// @Summary 上传文件（管理端）
// @Tags Admin.File
// @Security BearerAuth
// @Accept multipart/form-data
// @Produce json
// @Param file formance file true "文件"
// @Success 200 {object} shared.Envelope
// @Router /admin/file/upload [post]
func (a *Admin) uploadFile(c fiber.Ctx) error {
	file, err := c.FormFile("file")
	if err != nil {
		return shared.WriteError(c, http.StatusBadRequest, bizcode.ErrorParam, "file is required")
	}

	// 打开文件
	src, err := file.Open()
	if err != nil {
		return shared.WriteError(c, http.StatusInternalServerError, bizcode.ErrorSystem, "failed to open file")
	}
	defer src.Close()

	result, err := a.file.Upload(c.Context(), input.UploadFile{
		File:     src,
		Filename: file.Filename,
		Size:     file.Size,
	})

	if err != nil {
		a.logger.Error(err, "http - admin - file - uploadFile")
		return shared.WriteError(c, http.StatusInternalServerError, bizcode.ErrorSystem, "failed to upload file")
	}

	return shared.WriteSuccess(c, shared.WithData(result))
}

// deleteFile 删除文件。
// @Summary 删除文件（管理端）
// @Tags Admin.File
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param key query string true "文件 Key"
// @Success 200 {object} shared.Envelope
// @Router /admin/file/delete [delete]
func (a *Admin) deleteFile(c fiber.Ctx) error {
	key := c.Query("key")
	if key == "" {
		return shared.WriteError(c, http.StatusBadRequest, bizcode.ErrorParamMissing, "key is required")
	}

	if err := a.file.Delete(c.Context(), key); err != nil {
		a.logger.Error(err, "http - admin - file - deleteFile")
		return shared.WriteError(c, http.StatusInternalServerError, bizcode.ErrorSystem, "failed to delete file")
	}

	return shared.WriteSuccess(c)
}

// listImages 图片列表。
// @Summary 图片列表（管理端）
// @Tags Admin.Image
// @Security BearerAuth
// @Produce json
// @Param page query int false "页码"
// @Param page_size query int false "每页数量"
// @Param name query string false "文件名"
// @Success 200 {object} shared.Envelope
// @Router /admin/image/list [get]
func (a *Admin) listImages(c fiber.Ctx) error {
	pq := shared.ParsePageQueryWithOptions(c, shared.WithAllowedFilters("name", "mime_type"))

	result, err := a.file.List(c.Context(), input.ListFiles{
		PageParams: input.PageParams{Page: pq.Page, PageSize: pq.PageSize},
		Filename:   pq.Filters["name"],
		MimeType:   pq.Filters["mime_type"],
	})
	if err != nil {
		a.logger.Error(err, "http - admin - file - listImages")
		return shared.WriteError(c, http.StatusInternalServerError, bizcode.ErrorDatabase, "failed to list images")
	}

	return shared.WriteSuccess(c, shared.WithData(shared.NewPage(result.Items, result.Page, result.PageSize, result.Total)))
}

// deleteImages 批量删除图片。
// @Summary 批量删除图片（管理端）
// @Tags Admin.Image
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param ids body []int64 true "图片 ID 列表"
// @Success 200 {object} shared.Envelope
// @Router /admin/image/delete [delete]
func (a *Admin) deleteImages(c fiber.Ctx) error {
	var req struct {
		IDs []int64 `json:"ids"`
	}
	if err := c.Bind().JSON(&req); err != nil {
		return shared.WriteError(c, http.StatusBadRequest, bizcode.ErrorParam, "invalid request body")
	}

	if len(req.IDs) == 0 {
		return shared.WriteError(c, http.StatusBadRequest, bizcode.ErrorParamMissing, "ids is required")
	}

	if err := a.file.DeleteByIDs(c.Context(), req.IDs); err != nil {
		a.logger.Error(err, "http - admin - file - deleteImages")
		return shared.WriteError(c, http.StatusInternalServerError, bizcode.ErrorSystem, "failed to delete images")
	}

	return shared.WriteSuccess(c)
}

// listResources 资源列表。
// @Summary 资源列表（管理端）
// @Tags Admin.Resource
// @Security BearerAuth
// @Produce json
// @Param page query int false "页码"
// @Param page_size query int false "每页数量"
// @Param file_name query string false "文件名"
// @Param mime_type query string false "MIME类型"
// @Success 200 {object} shared.Envelope
// @Router /admin/resources/list [get]
func (a *Admin) listResources(c fiber.Ctx) error {
	pq := shared.ParsePageQueryWithOptions(c, shared.WithAllowedFilters("file_name", "mime_type"))

	result, err := a.resource.List(c.Context(), nil, input.ListResources{
		PageParams: input.PageParams{Page: pq.Page, PageSize: pq.PageSize},
		Filename:   pq.Filters["file_name"],
		MimeType:   pq.Filters["mime_type"],
	})
	if err != nil {
		a.logger.Error(err, "http - admin - resource - listResources")
		return shared.WriteError(c, http.StatusInternalServerError, bizcode.ErrorDatabase, "failed to list resources")
	}

	return shared.WriteSuccess(c, shared.WithData(shared.NewPage(result.Items, result.Page, result.PageSize, result.Total)))
}

// deleteResources 批量删除资源。
// @Summary 批量删除资源（管理端）
// @Tags Admin.Resource
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param ids body []int64 true "资源 ID 列表"
// @Success 200 {object} shared.Envelope
// @Router /admin/resources/delete [post]
func (a *Admin) deleteResources(c fiber.Ctx) error {
	var req struct {
		IDs []int64 `json:"ids"`
	}
	if err := c.Bind().JSON(&req); err != nil {
		return shared.WriteError(c, http.StatusBadRequest, bizcode.ErrorParam, "invalid request body")
	}

	if len(req.IDs) == 0 {
		return shared.WriteError(c, http.StatusBadRequest, bizcode.ErrorParamMissing, "ids is required")
	}

	if err := a.resource.DeleteByIDs(c.Context(), req.IDs); err != nil {
		a.logger.Error(err, "http - admin - resource - deleteResources")
		return shared.WriteError(c, http.StatusInternalServerError, bizcode.ErrorSystem, "failed to delete resources")
	}

	return shared.WriteSuccess(c)
}
