package admin

import (
	"net/http"

	"github.com/gofiber/fiber/v3"

	"server-blog-v2/internal/controller/http/bizcode"
	"server-blog-v2/internal/controller/http/shared"
	"server-blog-v2/internal/usecase/input"
)

// listAdvertisements 广告列表。
// @Summary 广告列表（管理端）
// @Tags Admin.Advertisement
// @Security BearerAuth
// @Produce json
// @Param page query int false "页码" default(1)
// @Param page_size query int false "分页大小" default(10)
// @Param title query string false "标题"
// @Success 200 {object} shared.Envelope
// @Router /admin/advertisement/list [get]
func (a *Admin) listAdvertisements(c fiber.Ctx) error {
	pq := shared.ParsePageQueryWithOptions(c, shared.WithAllowedFilters("title"))

	var title *string
	if t, ok := pq.Filters["title"]; ok && t != "" {
		title = &t
	}

	result, err := a.advertisement.List(c.Context(), input.ListAdvertisements{
		PageParams: input.PageParams{Page: pq.Page, PageSize: pq.PageSize},
		Title:      title,
	})

	if err != nil {
		a.logger.Error(err, "http - admin - advertisement - listAdvertisements")
		return shared.WriteError(c, http.StatusInternalServerError, bizcode.ErrorDatabase, "failed to list advertisements")
	}

	return shared.WriteSuccess(c, shared.WithData(shared.NewPage(result.Items, result.Page, result.PageSize, result.Total)))
}
