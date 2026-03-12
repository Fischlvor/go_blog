package admin

import (
	"net/http"

	"github.com/gofiber/fiber/v3"

	"server-blog-v2/internal/controller/http/admin/request"
	"server-blog-v2/internal/controller/http/bizcode"
	"server-blog-v2/internal/controller/http/shared"
	"server-blog-v2/internal/usecase/input"
)

// listUsers 用户列表。
// @Summary 用户列表（管理端）
// @Tags Admin.User
// @Security BearerAuth
// @Produce json
// @Param page query int false "页码" default(1)
// @Param page_size query int false "分页大小" default(10)
// @Param keyword query string false "关键字"
// @Param filter.status query string false "状态" Enums(active, frozen)
// @Success 200 {object} shared.Envelope
// @Router /admin/user/list [get]
func (a *Admin) listUsers(c fiber.Ctx) error {
	pq := shared.ParsePageQueryWithOptions(c, shared.WithAllowedFilters("status"))

	var keyword *input.KeywordParams
	if pq.Keyword != "" {
		keyword = &input.KeywordParams{Keyword: pq.Keyword}
	}

	var status *string
	if s, ok := pq.Filters["status"]; ok && s != "" {
		status = &s
	}

	result, err := a.user.List(c.Context(), input.ListUsers{
		PageParams: input.PageParams{Page: pq.Page, PageSize: pq.PageSize},
		Keyword:    keyword,
		Status:     status,
	})

	if err != nil {
		a.logger.Error(err, "http - admin - user - listUsers")
		return shared.WriteError(c, http.StatusInternalServerError, bizcode.ErrorDatabase, "failed to list users")
	}

	return shared.WriteSuccess(c, shared.WithData(shared.NewPage(result.Items, result.Page, result.PageSize, result.Total)))
}

// freezeUser 冻结用户。
// @Summary 冻结用户（管理端）
// @Tags Admin.User
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param body body request.FreezeUser true "用户信息"
// @Success 200 {object} shared.Envelope
// @Router /admin/user/freeze [put]
func (a *Admin) freezeUser(c fiber.Ctx) error {
	var req request.FreezeUser
	if err := c.Bind().JSON(&req); err != nil {
		return shared.WriteError(c, http.StatusBadRequest, bizcode.ErrorParam, "invalid request body")
	}

	if err := a.validate.Struct(req); err != nil {
		return shared.WriteError(c, http.StatusBadRequest, bizcode.ErrorParamFormat, err.Error())
	}

	if err := a.user.Freeze(c.Context(), req.UserUUID); err != nil {
		a.logger.Error(err, "http - admin - user - freezeUser")
		return shared.WriteError(c, http.StatusInternalServerError, bizcode.ErrorDatabase, "failed to freeze user")
	}

	return shared.WriteSuccess(c)
}

// unfreezeUser 解冻用户。
// @Summary 解冻用户（管理端）
// @Tags Admin.User
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param body body request.UnfreezeUser true "用户信息"
// @Success 200 {object} shared.Envelope
// @Router /admin/user/unfreeze [put]
func (a *Admin) unfreezeUser(c fiber.Ctx) error {
	var req request.UnfreezeUser
	if err := c.Bind().JSON(&req); err != nil {
		return shared.WriteError(c, http.StatusBadRequest, bizcode.ErrorParam, "invalid request body")
	}

	if err := a.validate.Struct(req); err != nil {
		return shared.WriteError(c, http.StatusBadRequest, bizcode.ErrorParamFormat, err.Error())
	}

	if err := a.user.Unfreeze(c.Context(), req.UserUUID); err != nil {
		a.logger.Error(err, "http - admin - user - unfreezeUser")
		return shared.WriteError(c, http.StatusInternalServerError, bizcode.ErrorDatabase, "failed to unfreeze user")
	}

	return shared.WriteSuccess(c)
}

// listLogins 登录记录列表。
// @Summary 登录记录列表（管理端）
// @Tags Admin.User
// @Security BearerAuth
// @Produce json
// @Param page query int false "页码" default(1)
// @Param page_size query int false "分页大小" default(10)
// @Param keyword query string false "关键字"
// @Success 200 {object} shared.Envelope
// @Router /admin/user/loginList [get]
func (a *Admin) listLogins(c fiber.Ctx) error {
	pq := shared.ParsePageQueryWithOptions(c)

	var keyword *input.KeywordParams
	if pq.Keyword != "" {
		keyword = &input.KeywordParams{Keyword: pq.Keyword}
	}

	result, err := a.user.ListLogins(c.Context(), input.ListLogins{
		PageParams: input.PageParams{Page: pq.Page, PageSize: pq.PageSize},
		Keyword:    keyword,
	})

	if err != nil {
		a.logger.Error(err, "http - admin - user - listLogins")
		return shared.WriteError(c, http.StatusInternalServerError, bizcode.ErrorDatabase, "failed to list logins")
	}

	return shared.WriteSuccess(c, shared.WithData(shared.NewPage(result.Items, result.Page, result.PageSize, result.Total)))
}
