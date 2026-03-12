package middleware

import (
	"fmt"
	"runtime/debug"

	"github.com/gofiber/fiber/v3"

	"server-blog-v2/pkg/logger"
)

// Recovery panic 恢复中间件。
func Recovery(l logger.Interface) fiber.Handler {
	return func(c fiber.Ctx) (err error) {
		defer func() {
			if r := recover(); r != nil {
				l.Error(fmt.Sprintf("panic recovered: %v\n%s", r, debug.Stack()))
				err = fiber.ErrInternalServerError
			}
		}()
		return c.Next()
	}
}
