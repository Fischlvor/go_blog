package qiniu

import (
	"github.com/gofiber/fiber/v3"

	"server-blog-v2/internal/controller/http/shared"
	"server-blog-v2/internal/usecase"
	"server-blog-v2/internal/usecase/input"
	"server-blog-v2/pkg/logger"
)

// Qiniu 七牛云回调控制器。
type Qiniu struct {
	logger   logger.Interface
	resource usecase.Resource
}

// NewQiniu 创建七牛云回调控制器。
func NewQiniu(l logger.Interface, resource usecase.Resource) *Qiniu {
	return &Qiniu{
		logger:   l,
		resource: resource,
	}
}

// NewRoutes 注册七牛云回调路由。
func NewRoutes(router fiber.Router, l logger.Interface, resource usecase.Resource) {
	q := NewQiniu(l, resource)

	router.Post("/callback", q.callback)
}

// callback 七牛云转码回调。
// POST /api/callback/qiniu/callback
// 此接口不需要认证，由七牛云服务器调用。
func (q *Qiniu) callback(c fiber.Ctx) error {
	var req input.QiniuCallback
	if err := c.Bind().JSON(&req); err != nil {
		q.logger.Error(err, "http - qiniu - callback - bind")
		// 返回成功，避免七牛重试（解析失败重试也没用）
		return shared.WriteSuccess(c)
	}

	q.logger.Info("七牛云转码回调",
		"id", req.ID,
		"code", req.Code,
		"inputKey", req.InputKey,
		"items_count", len(req.Items),
	)

	// 处理回调
	if err := q.resource.HandleQiniuCallback(c.Context(), req.InputKey, req.Code, req.Items); err != nil {
		q.logger.Error(err, "http - qiniu - callback - HandleQiniuCallback",
			"inputKey", req.InputKey,
		)
		// 仍然返回成功，避免无限重试
		// 错误已记录日志，可以后续手动处理
	}

	// 返回成功，否则七牛会重试
	return shared.WriteSuccess(c)
}
