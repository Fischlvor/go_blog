package shared

import (
	"net/http"

	"github.com/gofiber/fiber/v3"

	"server-blog-v2/internal/controller/http/bizcode"
)

// Option 响应选项。
type Option func(*Envelope)

// WithMsg 设置消息。
func WithMsg(msg string) Option { return func(e *Envelope) { e.Message = msg } }

// WithData 设置数据。
func WithData(data interface{}) Option { return func(e *Envelope) { e.Data = data } }

// WithCode 设置业务码。
func WithCode(code string) Option { return func(e *Envelope) { e.Code = code } }

// WriteSuccess 写入成功响应。
func WriteSuccess(ctx fiber.Ctx, opts ...Option) error {
	env := Envelope{Code: bizcode.Success, Message: "ok"}
	for _, opt := range opts {
		opt(&env)
	}
	return ctx.Status(http.StatusOK).JSON(env)
}

// WriteError 写入错误响应。
func WriteError(ctx fiber.Ctx, httpCode int, bizCode string, msg string) error {
	return ctx.Status(httpCode).JSON(Envelope{Code: bizCode, Message: msg})
}
