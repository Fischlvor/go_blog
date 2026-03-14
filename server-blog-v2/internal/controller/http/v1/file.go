package v1

import (
	"net/http"

	"github.com/gofiber/fiber/v3"

	"server-blog-v2/internal/controller/http/bizcode"
	"server-blog-v2/internal/controller/http/middleware"
	"server-blog-v2/internal/controller/http/shared"
	"server-blog-v2/internal/usecase/input"
)

// uploadFile 上传文件。
func (v *V1) uploadFile(c fiber.Ctx) error {
	file, err := c.FormFile("file")
	if err != nil {
		return shared.WriteError(c, http.StatusBadRequest, bizcode.ErrorParamMissing, "file is required")
	}

	f, err := file.Open()
	if err != nil {
		return shared.WriteError(c, http.StatusInternalServerError, bizcode.ErrorSystem, "failed to open file")
	}
	defer f.Close()

	usage := c.FormValue("usage", "article_content")
	userUUID := middleware.GetUserUUID(c)

	result, err := v.file.Upload(c.Context(), input.UploadFile{
		File:        f,
		Filename:    file.Filename,
		Size:        file.Size,
		ContentType: file.Header.Get("Content-Type"),
		Usage:       usage,
		UserUUID:    userUUID,
	})
	if err != nil {
		v.logger.Error(err, "http - v1 - file - uploadFile")
		return shared.WriteError(c, http.StatusInternalServerError, bizcode.ErrorThirdParty, "failed to upload file")
	}

	return shared.WriteSuccess(c, shared.WithData(result))
}
