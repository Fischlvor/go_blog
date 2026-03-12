package admin

import (
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v3"

	"server-blog-v2/internal/controller/http/bizcode"
	"server-blog-v2/internal/controller/http/middleware"
	"server-blog-v2/internal/controller/http/shared"
	"server-blog-v2/internal/usecase/input"
)

// getMaxFileSize 获取最大文件大小。
// @Summary 获取最大文件大小（管理端）
// @Tags Admin.Resource
// @Security BearerAuth
// @Produce json
// @Success 200 {object} shared.Envelope
// @Router /admin/resources/max-size [get]
func (a *Admin) getMaxFileSize(c fiber.Ctx) error {
	maxSize := a.resource.GetMaxFileSize(c.Context())
	return shared.WriteSuccess(c, shared.WithData(map[string]interface{}{
		"max_size": maxSize,
	}))
}

// checkResource 检查文件（秒传/续传检测）。
// @Summary 检查文件（管理端）
// @Tags Admin.Resource
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param body body input.ResourceCheck true "检查参数"
// @Success 200 {object} shared.Envelope
// @Router /admin/resources/check [post]
func (a *Admin) checkResource(c fiber.Ctx) error {
	userUUID := middleware.GetUserUUID(c)
	if userUUID == "" {
		return shared.WriteError(c, http.StatusUnauthorized, bizcode.ErrorUnauthorized, "unauthorized")
	}

	var req input.ResourceCheck
	if err := c.Bind().JSON(&req); err != nil {
		return shared.WriteError(c, http.StatusBadRequest, bizcode.ErrorParam, "invalid request body")
	}

	result, err := a.resource.Check(c.Context(), userUUID, req)
	if err != nil {
		a.logger.Error(err, "http - admin - resource - checkResource")
		return shared.WriteError(c, http.StatusInternalServerError, bizcode.ErrorSystem, err.Error())
	}

	return shared.WriteSuccess(c, shared.WithData(result))
}

// initResource 初始化上传任务。
// @Summary 初始化上传任务（管理端）
// @Tags Admin.Resource
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param body body input.ResourceInit true "初始化参数"
// @Success 200 {object} shared.Envelope
// @Router /admin/resources/init [post]
func (a *Admin) initResource(c fiber.Ctx) error {
	userUUID := middleware.GetUserUUID(c)
	if userUUID == "" {
		return shared.WriteError(c, http.StatusUnauthorized, bizcode.ErrorUnauthorized, "unauthorized")
	}

	var req input.ResourceInit
	if err := c.Bind().JSON(&req); err != nil {
		return shared.WriteError(c, http.StatusBadRequest, bizcode.ErrorParam, "invalid request body")
	}

	result, err := a.resource.Init(c.Context(), userUUID, req)
	if err != nil {
		a.logger.Error(err, "http - admin - resource - initResource")
		return shared.WriteError(c, http.StatusInternalServerError, bizcode.ErrorSystem, err.Error())
	}

	return shared.WriteSuccess(c, shared.WithData(result))
}

// uploadChunk 上传分片。
// @Summary 上传分片（管理端）
// @Tags Admin.Resource
// @Security BearerAuth
// @Accept multipart/form-data
// @Produce json
// @Param task_id formData string true "任务ID"
// @Param chunk_number formData int true "块号"
// @Param chunk_data formData file true "块数据"
// @Success 200 {object} shared.Envelope
// @Router /admin/resources/upload-chunk [post]
func (a *Admin) uploadChunk(c fiber.Ctx) error {
	userUUID := middleware.GetUserUUID(c)
	if userUUID == "" {
		return shared.WriteError(c, http.StatusUnauthorized, bizcode.ErrorUnauthorized, "unauthorized")
	}

	taskID := c.FormValue("task_id")
	if taskID == "" {
		return shared.WriteError(c, http.StatusBadRequest, bizcode.ErrorParamMissing, "task_id is required")
	}

	chunkNumberStr := c.FormValue("chunk_number")
	chunkNumber, err := strconv.Atoi(chunkNumberStr)
	if err != nil {
		return shared.WriteError(c, http.StatusBadRequest, bizcode.ErrorParam, "invalid chunk_number")
	}

	file, err := c.FormFile("chunk_data")
	if err != nil {
		return shared.WriteError(c, http.StatusBadRequest, bizcode.ErrorParamMissing, "chunk_data is required")
	}

	src, err := file.Open()
	if err != nil {
		return shared.WriteError(c, http.StatusInternalServerError, bizcode.ErrorSystem, "failed to open file")
	}
	defer src.Close()

	result, err := a.resource.UploadChunk(c.Context(), userUUID, input.ResourceUploadChunk{
		TaskID:      taskID,
		ChunkNumber: chunkNumber,
		ChunkData:   src,
		ChunkSize:   file.Size,
	})
	if err != nil {
		a.logger.Error(err, "http - admin - resource - uploadChunk")
		return shared.WriteError(c, http.StatusInternalServerError, bizcode.ErrorSystem, err.Error())
	}

	return shared.WriteSuccess(c, shared.WithData(result))
}

// completeResource 完成上传。
// @Summary 完成上传（管理端）
// @Tags Admin.Resource
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param body body input.ResourceComplete true "完成参数"
// @Success 200 {object} shared.Envelope
// @Router /admin/resources/complete [post]
func (a *Admin) completeResource(c fiber.Ctx) error {
	userUUID := middleware.GetUserUUID(c)
	if userUUID == "" {
		return shared.WriteError(c, http.StatusUnauthorized, bizcode.ErrorUnauthorized, "unauthorized")
	}

	var req input.ResourceComplete
	if err := c.Bind().JSON(&req); err != nil {
		return shared.WriteError(c, http.StatusBadRequest, bizcode.ErrorParam, "invalid request body")
	}

	result, err := a.resource.Complete(c.Context(), userUUID, req)
	if err != nil {
		a.logger.Error(err, "http - admin - resource - completeResource")
		return shared.WriteError(c, http.StatusInternalServerError, bizcode.ErrorSystem, err.Error())
	}

	return shared.WriteSuccess(c, shared.WithData(result))
}

// cancelResource 取消上传。
// @Summary 取消上传（管理端）
// @Tags Admin.Resource
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param body body input.ResourceCancel true "取消参数"
// @Success 200 {object} shared.Envelope
// @Router /admin/resources/cancel [post]
func (a *Admin) cancelResource(c fiber.Ctx) error {
	userUUID := middleware.GetUserUUID(c)
	if userUUID == "" {
		return shared.WriteError(c, http.StatusUnauthorized, bizcode.ErrorUnauthorized, "unauthorized")
	}

	var req input.ResourceCancel
	if err := c.Bind().JSON(&req); err != nil {
		return shared.WriteError(c, http.StatusBadRequest, bizcode.ErrorParam, "invalid request body")
	}

	if err := a.resource.Cancel(c.Context(), userUUID, req); err != nil {
		a.logger.Error(err, "http - admin - resource - cancelResource")
		return shared.WriteError(c, http.StatusInternalServerError, bizcode.ErrorSystem, err.Error())
	}

	return shared.WriteSuccess(c)
}

// progressResource 查询上传进度。
// @Summary 查询上传进度（管理端）
// @Tags Admin.Resource
// @Security BearerAuth
// @Produce json
// @Param task_id query string true "任务ID"
// @Success 200 {object} shared.Envelope
// @Router /admin/resources/progress [get]
func (a *Admin) progressResource(c fiber.Ctx) error {
	userUUID := middleware.GetUserUUID(c)
	if userUUID == "" {
		return shared.WriteError(c, http.StatusUnauthorized, bizcode.ErrorUnauthorized, "unauthorized")
	}

	taskID := c.Query("task_id")
	if taskID == "" {
		return shared.WriteError(c, http.StatusBadRequest, bizcode.ErrorParamMissing, "task_id is required")
	}

	result, err := a.resource.Progress(c.Context(), userUUID, input.ResourceProgress{TaskID: taskID})
	if err != nil {
		a.logger.Error(err, "http - admin - resource - progressResource")
		return shared.WriteError(c, http.StatusInternalServerError, bizcode.ErrorSystem, err.Error())
	}

	return shared.WriteSuccess(c, shared.WithData(result))
}
