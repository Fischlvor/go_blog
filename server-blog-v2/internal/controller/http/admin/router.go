package admin

import (
	"crypto/rsa"

	"github.com/gofiber/fiber/v3"

	"server-blog-v2/config"
	"server-blog-v2/internal/controller/http/middleware"
	"server-blog-v2/internal/repo"
	"server-blog-v2/internal/usecase"
	"server-blog-v2/pkg/logger"
)

// NewRoutes 注册 Admin 路由。
func NewRoutes(
	router fiber.Router,
	cfg *config.Config,
	l logger.Interface,
	publicKey *rsa.PublicKey,
	userRepo repo.UserRepo,
	content usecase.Content,
	comment usecase.Comment,
	feedback usecase.Feedback,
	link usecase.Link,
	file usecase.File,
	resource usecase.Resource,
	user usecase.User,
	emoji usecase.Emoji,
	aiChat usecase.AIChat,
	aiModel usecase.AIModel,
	website usecase.Website,
	advertisement usecase.Advertisement,
) {
	admin := New(cfg, l, content, comment, feedback, link, file, resource, user, emoji, aiChat, aiModel, website, advertisement)

	// 管理员 JWT 中间件（从数据库查询用户角色）
	adminRequired := middleware.NewAdminJWTMiddleware(publicKey, userRepo)

	// ==================== 文章管理 /article ====================
	articleGroup := router.Group("/article", adminRequired)
	{
		articleGroup.Get("/list", admin.listArticles)
		articleGroup.Get("/:id", admin.getArticle)
		articleGroup.Post("/create", admin.createArticle)
		articleGroup.Put("/update", admin.updateArticle)
		articleGroup.Delete("/delete", admin.deleteArticle)
	}

	// ==================== 分类管理 /category ====================
	categoryGroup := router.Group("/category", adminRequired)
	{
		categoryGroup.Get("/list", admin.listCategories)
		categoryGroup.Post("/create", admin.createCategory)
		categoryGroup.Put("/update", admin.updateCategory)
		categoryGroup.Delete("/delete/:id", admin.deleteCategory)
	}

	// ==================== 标签管理 /tag ====================
	tagGroup := router.Group("/tag", adminRequired)
	{
		tagGroup.Get("/list", admin.listTags)
		tagGroup.Post("/create", admin.createTag)
		tagGroup.Put("/update", admin.updateTag)
		tagGroup.Delete("/delete/:id", admin.deleteTag)
	}

	// ==================== 评论管理 /comment ====================
	commentGroup := router.Group("/comment", adminRequired)
	{
		commentGroup.Get("/list", admin.listComments)
		commentGroup.Delete("/delete/:id", admin.deleteComment)
	}

	// ==================== 反馈管理 /feedback ====================
	feedbackGroup := router.Group("/feedback", adminRequired)
	{
		feedbackGroup.Get("/list", admin.listFeedbacks)
		feedbackGroup.Delete("/delete", admin.deleteFeedback)
		feedbackGroup.Put("/reply", admin.replyFeedback)
	}

	// ==================== 友链管理 /friendLink ====================
	linkGroup := router.Group("/friendLink", adminRequired)
	{
		linkGroup.Get("/list", admin.listLinks)
		linkGroup.Post("/create", admin.createLink)
		linkGroup.Put("/update", admin.updateLink)
		linkGroup.Delete("/delete", admin.deleteLink)
	}

	// ==================== 用户管理 /user ====================
	userGroup := router.Group("/user", adminRequired)
	{
		userGroup.Get("/list", admin.listUsers)
		userGroup.Get("/loginList", admin.listLogins)
		userGroup.Put("/freeze", admin.freezeUser)
		userGroup.Put("/unfreeze", admin.unfreezeUser)
	}

	// ==================== 文件管理 /file ====================
	fileGroup := router.Group("/file", adminRequired)
	{
		fileGroup.Post("/upload", admin.uploadFile)
		fileGroup.Delete("/delete", admin.deleteFile)
	}

	// ==================== 图片管理 /image ====================
	imageGroup := router.Group("/image", adminRequired)
	{
		imageGroup.Get("/list", admin.listImages)
		imageGroup.Delete("/delete", admin.deleteImages)
	}

	// ==================== 资源管理 /resources ====================
	resourceGroup := router.Group("/resources", adminRequired)
	{
		resourceGroup.Get("/max-size", admin.getMaxFileSize)
		resourceGroup.Post("/check", admin.checkResource)
		resourceGroup.Post("/init", admin.initResource)
		resourceGroup.Post("/upload-chunk", admin.uploadChunk)
		resourceGroup.Post("/complete", admin.completeResource)
		resourceGroup.Post("/cancel", admin.cancelResource)
		resourceGroup.Get("/progress", admin.progressResource)
		resourceGroup.Get("/list", admin.listResources)
		resourceGroup.Post("/delete", admin.deleteResources)
	}

	// ==================== 表情管理 /emoji ====================
	emojiAdminGroup := router.Group("/emoji", adminRequired)
	{
		emojiAdminGroup.Get("/list", admin.listEmojis)
		emojiAdminGroup.Get("/sprites", admin.listSprites)
		emojiAdminGroup.Post("/regenerate", admin.regenerateSprites)
	}

	// ==================== AI 管理 /ai-management ====================
	aiGroup := router.Group("/ai-management", adminRequired)
	{
		// AI 会话管理
		aiGroup.Get("/session", admin.listAISessions)
		aiGroup.Get("/session/:id", admin.getAISession)
		aiGroup.Delete("/session/:id", admin.deleteAISession)

		// AI 消息管理
		aiGroup.Get("/message", admin.listAIMessages)

		// AI 模型管理
		aiGroup.Get("/model", admin.listAIModels)
		aiGroup.Get("/model/:id", admin.getAIModel)
		aiGroup.Post("/model", admin.createAIModel)
		aiGroup.Put("/model", admin.updateAIModel)
		aiGroup.Delete("/model/:id", admin.deleteAIModel)
	}

	// ==================== 广告管理 /advertisement ====================
	advertisementGroup := router.Group("/advertisement", adminRequired)
	{
		advertisementGroup.Get("/list", admin.listAdvertisements)
	}

	// ==================== 配置管理 /config ====================
	configGroup := router.Group("/config", adminRequired)
	{
		configGroup.Get("/website", admin.getWebsiteConfig)
		configGroup.Get("/system", admin.getSystemConfig)
		configGroup.Get("/email", admin.getEmailConfig)
		configGroup.Get("/qq", admin.getQQConfig)
		configGroup.Get("/qiniu", admin.getQiniuConfig)
		configGroup.Get("/jwt", admin.getJwtConfig)
		configGroup.Get("/gaode", admin.getGaodeConfig)
	}
}
