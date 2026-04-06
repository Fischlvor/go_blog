package usecase

import (
	"context"

	"server-blog-v2/internal/usecase/input"
	"server-blog-v2/internal/usecase/output"
)

// ==================== 文章 ====================

// Content 内容管理用例（文章、分类、标签）。
type Content interface {
	// 文章 - 管理端
	ListArticles(ctx context.Context, params input.ListArticles) (*output.ListResult[output.ArticleSummary], error)
	GetArticleByID(ctx context.Context, id int64) (*output.ArticleDetail, error)
	GetArticleBySlug(ctx context.Context, slug string) (*output.ArticleDetail, error)
	CreateArticle(ctx context.Context, params input.CreateArticle) (string, error)
	UpdateArticle(ctx context.Context, params input.UpdateArticle) error
	DeleteArticle(ctx context.Context, id int64) error
	DeleteArticleBySlug(ctx context.Context, slug string) error

	// 文章 - 公开端
	ListPublicArticles(ctx context.Context, params input.ListPublicArticles, userUUID *string) (*output.ListResult[output.ArticleSummary], error)
	GetPublicArticleBySlug(ctx context.Context, slug string, userUUID *string) (*output.ArticleDetail, error)
	RecordView(ctx context.Context, articleSlug string, ip, userAgent, referer string)

	// 点赞
	ToggleLikeOnArticle(ctx context.Context, articleSlug, userUUID string) (liked bool, count int32, err error)
	RemoveLikeOnArticle(ctx context.Context, articleSlug, userUUID string) (removed bool, count int32, err error)
	ListUserLikedArticles(ctx context.Context, userUUID string, params input.ListUserLikedArticles) (*output.ListResult[output.ArticleSummary], error)

	// 分类
	ListCategories(ctx context.Context, params input.ListCategories) (*output.ListResult[output.CategoryDetail], error)
	GetAllPublicCategories(ctx context.Context) (*output.AllResult[output.CategoryDetail], error)
	CreateCategory(ctx context.Context, params input.CreateCategory) (int64, error)
	UpdateCategory(ctx context.Context, params input.UpdateCategory) error
	DeleteCategory(ctx context.Context, id int64) error

	// 标签
	ListTags(ctx context.Context, params input.ListTags) (*output.ListResult[output.TagDetail], error)
	GetAllPublicTags(ctx context.Context) (*output.AllResult[output.TagDetail], error)
	CreateTag(ctx context.Context, params input.CreateTag) (int64, error)
	UpdateTag(ctx context.Context, params input.UpdateTag) error
	DeleteTag(ctx context.Context, id int64) error
}

// ==================== 评论 ====================

// Comment 评论用例。
type Comment interface {
	// 公开端
	ListByArticleSlug(ctx context.Context, articleSlug string, params input.ListComments) (*output.ListResult[output.Comment], error)
	Create(ctx context.Context, params input.CreateComment) (int64, error)
	Delete(ctx context.Context, id int64, userUUID string) error

	// 管理端
	ListAll(ctx context.Context, params input.ListAllComments) (*output.ListResult[output.CommentAdmin], error)
	AdminDelete(ctx context.Context, id int64) error
}

// ==================== AI 聊天 ====================

// AIChat AI 聊天用例。
type AIChat interface {
	// 用户端
	CreateSession(ctx context.Context, userUUID string, params input.CreateSession) (*output.Session, error)
	ListSessions(ctx context.Context, userUUID string, params input.ListSessions) (*output.ListResult[output.Session], error)
	GetMessages(ctx context.Context, sessionID int64, userUUID string) ([]*output.Message, error)
	SendMessage(ctx context.Context, sessionID int64, userUUID string, params input.SendMessage) (<-chan string, error)
	DeleteSession(ctx context.Context, sessionID int64, userUUID string) error

	// 管理端
	ListAllSessions(ctx context.Context, params input.ListAllSessions) (*output.ListResult[output.SessionAdmin], error)
	GetSessionByID(ctx context.Context, sessionID int64) (*output.SessionAdmin, error)
	AdminDeleteSession(ctx context.Context, sessionID int64) error
	ListAllMessages(ctx context.Context, params input.ListAllMessages) (*output.ListResult[output.MessageAdmin], error)
}

// AIModel AI 模型管理用例。
type AIModel interface {
	List(ctx context.Context, params input.ListAIModels) (*output.ListResult[output.AIModelInfo], error)
	GetByID(ctx context.Context, id int64) (*output.AIModelInfo, error)
	Create(ctx context.Context, params input.CreateAIModel) (int64, error)
	Update(ctx context.Context, params input.UpdateAIModel) error
	Delete(ctx context.Context, id int64) error
}

// ==================== 反馈 ====================

// Feedback 反馈用例。
type Feedback interface {
	// 用户端
	Create(ctx context.Context, params input.CreateFeedback) error

	// 管理端
	List(ctx context.Context, params input.ListFeedback) (*output.ListResult[output.Feedback], error)
	Delete(ctx context.Context, id int64) error
	Reply(ctx context.Context, id int64, reply string) error
}

// ==================== 友链 ====================

// Link 友链用例。
type Link interface {
	List(ctx context.Context) ([]*output.FriendLink, error)
	Create(ctx context.Context, params input.CreateLink) (int64, error)
	Update(ctx context.Context, params input.UpdateLink) error
	Delete(ctx context.Context, id int64) error
}

// ==================== 文件 ====================

// File 文件上传用例。
type File interface {
	Upload(ctx context.Context, params input.UploadFile) (*output.UploadResult, error)
	Delete(ctx context.Context, key string) error
	List(ctx context.Context, params input.ListFiles) (*output.ListResult[output.FileInfo], error)
	DeleteByIDs(ctx context.Context, ids []int64) error
}

// Resource 资源管理用例。
type Resource interface {
	List(ctx context.Context, userUUID *string, params input.ListResources) (*output.ListResult[output.ResourceInfo], error)
	DeleteByIDs(ctx context.Context, ids []int64) error
	// 分片上传相关
	GetMaxFileSize(ctx context.Context) int64
	Check(ctx context.Context, userUUID string, params input.ResourceCheck) (*output.ResourceCheckResponse, error)
	Init(ctx context.Context, userUUID string, params input.ResourceInit) (*output.ResourceInitResponse, error)
	UploadChunk(ctx context.Context, userUUID string, params input.ResourceUploadChunk) (*output.ResourceUploadChunkResponse, error)
	Complete(ctx context.Context, userUUID string, params input.ResourceComplete) (*output.ResourceCompleteResponse, error)
	Cancel(ctx context.Context, userUUID string, params input.ResourceCancel) error
	Progress(ctx context.Context, userUUID string, params input.ResourceProgress) (*output.ResourceProgressResponse, error)
	// 七牛云回调
	HandleQiniuCallback(ctx context.Context, inputKey string, code int, items []input.QiniuCallbackItem) error
}

// ==================== 用户 ====================

// User 用户用例。
type User interface {
	// 用户端
	GetProfile(ctx context.Context, userUUID string) (*output.UserProfile, error)
	UpdateProfile(ctx context.Context, userUUID string, params input.UpdateProfile) error
	GetChart(ctx context.Context, days int) (*output.UserChart, error)

	// 管理端
	List(ctx context.Context, params input.ListUsers) (*output.ListResult[output.UserAdmin], error)
	ListLogins(ctx context.Context, params input.ListLogins) (*output.ListResult[output.LoginInfo], error)
	Freeze(ctx context.Context, userUUID string) error
	Unfreeze(ctx context.Context, userUUID string) error
}

// ==================== 网站 ====================

// Website 网站信息用例。
type Website interface {
	GetInfo(ctx context.Context) *output.WebsiteInfo
	GetCarousel(ctx context.Context) ([]string, error)
	GetNews(ctx context.Context, source string) (*output.HotSearchData, error)
	GetCalendar(ctx context.Context) (*output.CalendarData, error)
	GetFooterLinks(ctx context.Context) ([]output.FooterLink, error)
}

// Setting 站点配置用例。
type Setting interface {
	GetAllSiteSettings(ctx context.Context) (*output.AllResult[output.SiteSettingDetail], error)
	GetPublicSiteSettings(ctx context.Context) (*output.AllResult[output.SiteSettingDetail], error)
	GetSiteSettingByKey(ctx context.Context, key string) (*output.SiteSettingDetail, error)
	GetWebsiteSettingsMap(ctx context.Context) (map[string]string, error)
	UpdateSiteSettings(ctx context.Context, settings []input.UpsertSiteSetting) error
}

// ==================== 表情 ====================

// SSEWriter SSE 响应写入器接口。
type SSEWriter interface {
	WriteSSEEvent(event string, data interface{}) error
	Flush()
}

// Emoji 表情用例。
type Emoji interface {
	GetConfig(ctx context.Context) (*output.EmojiConfig, error)
	ListGroups(ctx context.Context) ([]*output.EmojiGroup, error)
	List(ctx context.Context, params input.ListEmojis) (*output.ListResult[output.EmojiInfo], error)
	ListSprites(ctx context.Context) ([]*output.SpriteInfo, error)
	// 雪碧图生成
	RegenerateSprites(ctx context.Context, groupKeys []string, sseWriter SSEWriter) error
}

// ==================== 广告 ====================

// Advertisement 广告用例。
type Advertisement interface {
	GetInfo(ctx context.Context) (*output.AdvertisementInfo, error)
	List(ctx context.Context, params input.ListAdvertisements) (*output.ListResult[output.AdvertisementItem], error)
}
