//go:build wireinject
// +build wireinject

// Package app 应用装配与生命周期管理。

package app

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
	"strconv"

	"github.com/google/wire"
	"gorm.io/gorm"

	"server-blog-v2/config"
	httpctrl "server-blog-v2/internal/controller/http"
	"server-blog-v2/internal/controller/http/middleware"
	"server-blog-v2/internal/repo"
	"server-blog-v2/internal/repo/persistence"
	"server-blog-v2/internal/repo/storage"
	"server-blog-v2/internal/repo/webapi"
	"server-blog-v2/internal/usecase"
	"server-blog-v2/internal/usecase/advertisement"
	"server-blog-v2/internal/usecase/aimodel"
	"server-blog-v2/internal/usecase/chat"
	"server-blog-v2/internal/usecase/comment"
	"server-blog-v2/internal/usecase/content"
	"server-blog-v2/internal/usecase/emoji"
	"server-blog-v2/internal/usecase/feedback"
	"server-blog-v2/internal/usecase/file"
	"server-blog-v2/internal/usecase/link"
	"server-blog-v2/internal/usecase/resource"
	"server-blog-v2/internal/usecase/user"
	"server-blog-v2/internal/usecase/website"
	"server-blog-v2/pkg/httpserver"
	"server-blog-v2/pkg/logger"
	"server-blog-v2/pkg/postgres"
	pkgRedis "server-blog-v2/pkg/redis"
)

// App 应用容器。
type App struct {
	Info       AppInfo
	Logger     logger.Interface
	HTTPServer *httpserver.Server
}

// AppInfo 应用信息。
type AppInfo struct {
	Name    string
	Version string
}

// NewApp 创建 App。
func NewApp(info AppInfo, l logger.Interface, srv *httpserver.Server) *App {
	return &App{
		Info:       info,
		Logger:     l,
		HTTPServer: srv,
	}
}

// NewAppInfo 创建 AppInfo。
func NewAppInfo(cfg *config.Config) AppInfo {
	return AppInfo{Name: cfg.App.Name, Version: cfg.App.Version}
}

// ==================== Infrastructure ====================

// NewLogger 创建 Logger。
func NewLogger(cfg *config.Config) logger.Interface {
	return logger.New(cfg.Log.Level)
}

// NewPostgres 创建 Postgres 连接并返回 cleanup。
func NewPostgres(cfg *config.Config) (*postgres.Postgres, func(), error) {
	pg, err := postgres.New(
		cfg.Postgres.Host,
		cfg.Postgres.Port,
		cfg.Postgres.User,
		cfg.Postgres.Password,
		cfg.Postgres.DBName,
		cfg.Postgres.SSLMode,
		cfg.Postgres.TimeZone,
		cfg.Postgres.MaxIdleConns,
		cfg.Postgres.MaxOpenConns,
	)
	if err != nil {
		return nil, nil, err
	}
	cleanup := func() { pg.Close() }
	return pg, cleanup, nil
}

// NewGormDB 从 Postgres 提取 *gorm.DB。
func NewGormDB(pg *postgres.Postgres) *gorm.DB {
	return pg.DB
}

// NewRedis 创建 Redis 连接并返回 cleanup。
func NewRedis(cfg *config.Config) (*pkgRedis.Redis, func(), error) {
	rdb, err := pkgRedis.New(cfg.Redis.Host, cfg.Redis.Port, cfg.Redis.Password, cfg.Redis.DB)
	if err != nil {
		return nil, nil, err
	}
	cleanup := func() { rdb.Close() }
	return rdb, cleanup, nil
}

// NewPublicKey 加载 SSO 公钥。
func NewPublicKey(cfg *config.Config) (*rsa.PublicKey, error) {
	keyData, err := os.ReadFile(cfg.SSO.PublicKeyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read public key: %w", err)
	}

	block, _ := pem.Decode(keyData)
	if block == nil {
		return nil, fmt.Errorf("failed to decode PEM block")
	}

	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse public key: %w", err)
	}

	rsaPub, ok := pub.(*rsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("not an RSA public key")
	}

	return rsaPub, nil
}

// ==================== Repo ====================

// NewObjectStore 创建七牛云存储。
func NewObjectStore(cfg *config.Config) repo.ObjectStore {
	return storage.NewQiniuStore(
		cfg.Qiniu.AccessKey,
		cfg.Qiniu.SecretKey,
		cfg.Qiniu.Bucket,
		cfg.Qiniu.Domain,
		cfg.Qiniu.Zone,
		cfg.Qiniu.UseHTTPS,
		"",
	)
}

// NewLLMWebAPI 创建 LLM API 客户端。
func NewLLMWebAPI(cfg *config.Config) repo.LLMWebAPI {
	return webapi.NewLLMWebAPI(
		cfg.AI.APIKey,
		cfg.AI.BaseURL,
		cfg.AI.Model,
	)
}

// NewSSOClient 创建 SSO 客户端。
func NewSSOClient(cfg *config.Config) *webapi.SSOClient {
	return webapi.NewSSOClient(
		cfg.SSO.ServiceURL,
		cfg.SSO.ClientID,
		cfg.SSO.ClientSecret,
	)
}

// ==================== UseCase ====================

// NewContentUseCase 创建 Content UseCase。
func NewContentUseCase(
	cfg *config.Config,
	articles repo.ArticleRepo,
	tags repo.TagRepo,
	categories repo.CategoryRepo,
	articleLikes repo.ArticleLikeRepo,
	articleViews repo.ArticleViewRepo,
	users repo.UserRepo,
) usecase.Content {
	return content.New(cfg, articles, tags, categories, articleLikes, articleViews, users)
}

// NewCommentUseCase 创建 Comment UseCase。
func NewCommentUseCase(cfg *config.Config, comments repo.CommentRepo, users repo.UserRepo) usecase.Comment {
	return comment.New(cfg, comments, users)
}

// NewAIChatUseCase 创建 AIChat UseCase。
func NewAIChatUseCase(
	cfg *config.Config,
	sessions repo.ChatSessionRepo,
	messages repo.ChatMessageRepo,
	llm repo.LLMWebAPI,
	users repo.UserRepo,
) usecase.AIChat {
	return chat.New(cfg, sessions, messages, llm, users)
}

// NewFeedbackUseCase 创建 Feedback UseCase。
func NewFeedbackUseCase(feedbacks repo.FeedbackRepo) usecase.Feedback {
	return feedback.New(feedbacks)
}

// NewLinkUseCase 创建 Link UseCase。
func NewLinkUseCase(cfg *config.Config, links repo.LinkRepo) usecase.Link {
	return link.New(cfg, links)
}

// NewFileUseCase 创建 File UseCase。
func NewFileUseCase(files repo.FileRepo, objectStore repo.ObjectStore) usecase.File {
	return file.New(files, objectStore)
}

// NewResourceUseCase 创建 Resource UseCase。
func NewResourceUseCase(resources repo.ResourceRepo, tasks repo.ResourceUploadTaskRepo, objectStore repo.ObjectStore, rdb *pkgRedis.Redis) usecase.Resource {
	return resource.New(resources, tasks, objectStore, rdb)
}

// NewAIModelUseCase 创建 AIModel UseCase。
func NewAIModelUseCase(models repo.AIModelRepo) usecase.AIModel {
	return aimodel.New(models)
}

// NewUserUseCase 创建 User UseCase。
func NewUserUseCase(cfg *config.Config, users repo.UserRepo) usecase.User {
	return user.New(cfg, users)
}

// NewWebsiteUseCase 创建 Website UseCase。
func NewWebsiteUseCase(cfg *config.Config, redis *pkgRedis.Redis, footerLinks repo.FooterLinkRepo) usecase.Website {
	return website.New(cfg, redis, footerLinks)
}

// NewEmojiUseCase 创建 Emoji UseCase。
func NewEmojiUseCase(cfg *config.Config, emojis repo.EmojiRepo, sprites repo.EmojiSpriteRepo) usecase.Emoji {
	return emoji.New(cfg, emojis, sprites)
}

// NewAdvertisementUseCase 创建 Advertisement UseCase。
func NewAdvertisementUseCase(cfg *config.Config, ads repo.AdvertisementRepo) usecase.Advertisement {
	return advertisement.New(cfg, ads)
}

// NewSessionManager 创建 Session 管理器。
func NewSessionManager(redis *pkgRedis.Redis, l logger.Interface) *middleware.SessionManager {
	return middleware.NewSessionManager(redis, l)
}

// ==================== HTTP Server ====================

func SetupHTTPServer(
	cfg *config.Config,
	l logger.Interface,
	publicKey *rsa.PublicKey,
	userRepo repo.UserRepo,
	contentUC usecase.Content,
	commentUC usecase.Comment,
	aiChatUC usecase.AIChat,
	aiModelUC usecase.AIModel,
	feedbackUC usecase.Feedback,
	linkUC usecase.Link,
	fileUC usecase.File,
	resourceUC usecase.Resource,
	userUC usecase.User,
	websiteUC usecase.Website,
	emojiUC usecase.Emoji,
	advertisementUC usecase.Advertisement,
	sessionManager *middleware.SessionManager,
	ssoClient *webapi.SSOClient,
) *httpserver.Server {
	srv := httpserver.New(l, httpserver.WithPort(strconv.Itoa(cfg.HTTP.Port)), httpserver.WithPrefork(cfg.HTTP.UsePreforkMode))
	httpctrl.NewRouter(srv.App, cfg, l, publicKey, userRepo, contentUC, commentUC, aiChatUC, aiModelUC, feedbackUC, linkUC, fileUC, resourceUC, userUC, websiteUC, emojiUC, advertisementUC, sessionManager, ssoClient)
	return srv
}

// ==================== Wire ProviderSet ====================

var ProviderSet = wire.NewSet(
	// App
	NewAppInfo,
	NewLogger,
	NewApp,

	// Infrastructure
	NewPostgres,
	NewGormDB,
	NewRedis,
	NewPublicKey,

	// Repo - Persistence (PostgreSQL)
	persistence.NewArticleRepo,
	persistence.NewArticleLikeRepo,
	persistence.NewArticleViewRepo,
	persistence.NewCategoryRepo,
	persistence.NewTagRepo,
	persistence.NewCommentRepo,
	persistence.NewUserRepo,
	persistence.NewChatSessionRepo,
	persistence.NewChatMessageRepo,
	persistence.NewFeedbackRepo,
	persistence.NewLinkRepo,
	persistence.NewFileRepo,
	persistence.NewResourceRepo,
	persistence.NewResourceUploadTaskRepo,
	persistence.NewAIModelRepo,
	persistence.NewEmojiRepo,
	persistence.NewEmojiSpriteRepo,
	persistence.NewAdvertisementRepo,
	persistence.NewFooterLinkRepo,

	// Repo - Storage & WebAPI
	NewObjectStore,
	NewLLMWebAPI,
	NewSSOClient,

	// UseCase
	NewContentUseCase,
	NewCommentUseCase,
	NewAIChatUseCase,
	NewFeedbackUseCase,
	NewLinkUseCase,
	NewFileUseCase,
	NewResourceUseCase,
	NewAIModelUseCase,
	NewUserUseCase,
	NewWebsiteUseCase,
	NewEmojiUseCase,
	NewAdvertisementUseCase,

	// Session
	NewSessionManager,

	// HTTP Server
	SetupHTTPServer,
)

// InitializeApp 初始化 App 并返回 cleanup。
func InitializeApp(cfg *config.Config) (*App, func(), error) {
	wire.Build(ProviderSet)
	return nil, nil, nil
}
