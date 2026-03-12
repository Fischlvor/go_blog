package v1

import (
	"net/http"

	"github.com/gofiber/fiber/v3"

	"server-blog-v2/internal/controller/http/bizcode"
	"server-blog-v2/internal/controller/http/shared"
)

// listLinks 友链列表。
func (v *V1) listLinks(c fiber.Ctx) error {
	links, err := v.link.List(c.Context())
	if err != nil {
		v.logger.Error(err, "http - v1 - link - listLinks")
		return shared.WriteError(c, http.StatusInternalServerError, bizcode.ErrorDatabase, "failed to list links")
	}

	return shared.WriteSuccess(c, shared.WithData(links))
}
