package v1

import (
	"crypto/rsa"

	"github.com/gofiber/fiber/v3"

	"server-blog-v2/config"
	"server-blog-v2/internal/controller/http/middleware"
	"server-blog-v2/internal/repo"
	"server-blog-v2/internal/repo/webapi"
	"server-blog-v2/internal/usecase"
	"server-blog-v2/pkg/logger"
)

// NewRoutes 注册 V1 路由（兼容原 server-blog 路由结构）。
func NewRoutes(
	router fiber.Router,
	cfg *config.Config,
	l logger.Interface,
	publicKey *rsa.PublicKey,
	content usecase.Content,
	comment usecase.Comment,
	aiChat usecase.AIChat,
	feedback usecase.Feedback,
	link usecase.Link,
	file usecase.File,
	user usecase.User,
	website usecase.Website,
	emoji usecase.Emoji,
	advertisement usecase.Advertisement,
	sessionManager *middleware.SessionManager,
	ssoClient *webapi.SSOClient,
	userRepo repo.UserRepo,
) {
	v1 := New(cfg, l, content, comment, aiChat, feedback, link, file, user, website, emoji, advertisement, sessionManager)

	// SSO JWT 中间件配置（支持自动刷新 token）
	ssoJWTConfig := middleware.SSOJWTConfig{
		PublicKey:          publicKey,
		SSOClient:          ssoClient,
		UserRoleGetter:     userRepo, // UserRepo 实现了 GetRoleByUUID
		UserCreator:        userRepo, // UserRepo 实现了 CreateFromSSO
		RefreshTokenGetter: sessionManager,
		Logger:             l,
	}

	// JWT 中间件（支持自动刷新）
	jwtRequired := middleware.NewSSOUserJWTMiddleware(ssoJWTConfig)
	jwtOptional := middleware.NewOptionalUserJWTMiddleware(publicKey)

	// ==================== 认证 /auth ====================
	authGroup := router.Group("/auth")
	{
		authGroup.Get("/sso_login_url", v1.getSSOLoginURL)
		authGroup.Get("/callback", v1.ssoCallback)
	}

	// ==================== 文章 /article ====================
	articleGroup := router.Group("/article")
	{
		// 公开接口
		articleGroup.Get("/search", v1.listArticles, jwtOptional)
		articleGroup.Get("/category", v1.listCategories)
		articleGroup.Get("/tags", v1.listTags)
		// 需要登录（放在 :slug 之前避免被匹配）
		articleGroup.Get("/likes", v1.listUserLikedArticles, jwtRequired)
		// 动态路由放最后
		articleGroup.Get("/:slug", v1.getArticle, jwtOptional)
		articleGroup.Post("/:slug/like", v1.toggleArticleLike, jwtRequired)
		articleGroup.Delete("/:slug/like", v1.removeArticleLike, jwtRequired)
	}

	// ==================== 评论 /comment ====================
	commentGroup := router.Group("/comment")
	{
		// 公开接口
		commentGroup.Get("/new", v1.listNewComments)
		commentGroup.Get("/:article_slug", v1.listComments)
		// 需要登录
		commentGroup.Post("/create", v1.createComment, jwtRequired)
		commentGroup.Delete("/delete", v1.deleteComment, jwtRequired)
	}

	// ==================== AI 聊天 /ai-chat ====================
	aiChatGroup := router.Group("/ai-chat")
	{
		// 公开接口
		aiChatGroup.Get("/models", v1.getAvailableModels)
		// 需要登录
		aiChatGroup.Post("/session", v1.createSession, jwtRequired)
		aiChatGroup.Get("/sessions", v1.listSessions, jwtRequired)
		aiChatGroup.Get("/session/:id", v1.getSessionDetail, jwtRequired)
		aiChatGroup.Get("/messages", v1.getMessages, jwtRequired)
		aiChatGroup.Post("/message", v1.sendMessage, jwtRequired)
		aiChatGroup.Post("/message/stream", v1.sendMessageStream, jwtRequired)
		aiChatGroup.Delete("/session", v1.deleteSession, jwtRequired)
		aiChatGroup.Put("/session", v1.updateSession, jwtRequired)
	}

	// ==================== 反馈 /feedback ====================
	feedbackGroup := router.Group("/feedback")
	{
		// 公开接口
		feedbackGroup.Get("/new", v1.listNewFeedback)
		// 需要登录
		feedbackGroup.Post("/create", v1.createFeedback, jwtRequired)
		feedbackGroup.Get("/info", v1.getFeedbackInfo, jwtRequired)
	}

	// ==================== 友链 /friendLink ====================
	friendLinkGroup := router.Group("/friendLink")
	{
		// 公开接口
		friendLinkGroup.Get("/info", v1.listLinks)
	}

	// ==================== 网站 /website ====================
	websiteGroup := router.Group("/website")
	{
		// 公开接口
		websiteGroup.Get("/logo", v1.getWebsiteLogo)
		websiteGroup.Get("/title", v1.getWebsiteTitle)
		websiteGroup.Get("/info", v1.getWebsiteInfo)
		websiteGroup.Get("/carousel", v1.getWebsiteCarousel)
		websiteGroup.Get("/news", v1.getWebsiteNews)
		websiteGroup.Get("/calendar", v1.getWebsiteCalendar)
		websiteGroup.Get("/footerLink", v1.getWebsiteFooterLink)
	}

	// ==================== 配置 /config ====================
	configGroup := router.Group("/config")
	{
		configGroup.Get("/website", v1.getWebsiteInfo)
	}

	// ==================== 用户 /user ====================
	userGroup := router.Group("/user")
	{
		// 公开接口
		userGroup.Get("/card", v1.getUserCard)
		// 需要登录
		userGroup.Get("/info", v1.getProfile, jwtRequired)
		userGroup.Put("/changeInfo", v1.updateProfile, jwtRequired)
		userGroup.Post("/logout", v1.logout, jwtRequired)
		userGroup.Get("/weather", v1.getUserWeather, jwtRequired)
		userGroup.Get("/chart", v1.getUserChart, jwtRequired)
	}

	// ==================== 文件上传 /image ====================
	imageGroup := router.Group("/image", jwtRequired)
	{
		imageGroup.Post("/upload", v1.uploadFile)
	}

	// ==================== 表情 /emoji ====================
	emojiGroup := router.Group("/emoji")
	{
		// 公开接口
		emojiGroup.Get("/config", v1.getEmojiConfig)
		emojiGroup.Get("/groups", v1.getEmojiGroups)
	}

	// ==================== 广告 /advertisement ====================
	advertisementGroup := router.Group("/advertisement")
	{
		// 公开接口
		advertisementGroup.Get("/info", v1.getAdvertisementInfo)
	}
}
