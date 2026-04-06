package http

import (
	"crypto/rsa"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"

	"server-blog-v2/config"
	"server-blog-v2/internal/controller/http/admin"
	"server-blog-v2/internal/controller/http/middleware"
	"server-blog-v2/internal/controller/http/qiniu"
	v1 "server-blog-v2/internal/controller/http/v1"
	"server-blog-v2/internal/repo"
	"server-blog-v2/internal/repo/webapi"
	"server-blog-v2/internal/usecase"
	"server-blog-v2/pkg/logger"
)

// NewRouter 创建路由。
func NewRouter(
	app *fiber.App,
	cfg *config.Config,
	l logger.Interface,
	publicKey *rsa.PublicKey,
	userRepo repo.UserRepo,
	content usecase.Content,
	comment usecase.Comment,
	aiChat usecase.AIChat,
	aiModel usecase.AIModel,
	feedback usecase.Feedback,
	link usecase.Link,
	file usecase.File,
	resource usecase.Resource,
	user usecase.User,
	setting usecase.Setting,
	website usecase.Website,
	emoji usecase.Emoji,
	advertisement usecase.Advertisement,
	sessionManager *middleware.SessionManager,
	ssoClient *webapi.SSOClient,
) {
	// 全局中间件
	app.Use(middleware.Logger(l))
	app.Use(middleware.Recovery(l))
	app.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "http://localhost:5173", "http://127.0.0.1:3000", "http://127.0.0.1:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"X-New-Access-Token", "X-Token-Expires-In"}, // 暴露给前端的响应头
		AllowCredentials: true,
	}))

	// 健康检查
	app.Get("/healthz", func(c fiber.Ctx) error {
		return c.SendString("ok")
	})

	// API 路由组
	api := app.Group("/api")

	// V1 公开 API
	v1Group := api.Group("/v1")
	v1.NewRoutes(v1Group, cfg, l, publicKey, content, comment, aiChat, feedback, link, file, user, setting, website, emoji, advertisement, sessionManager, ssoClient, userRepo)

	// Admin API
	adminGroup := api.Group("/admin")
	admin.NewRoutes(adminGroup, cfg, l, publicKey, userRepo, sessionManager, ssoClient, content, comment, feedback, link, file, resource, user, setting, emoji, aiChat, aiModel, website, advertisement)

	// 第三方回调（无需认证）
	callbackGroup := api.Group("/callback")
	qiniu.NewRoutes(callbackGroup.Group("/qiniu"), l, resource)
}
