package repo

import (
	"context"
	"io"
	"time"

	"server-blog-v2/internal/entity"
)

// ==================== 文章相关 ====================

// ArticleRepo 文章数据仓库 (PostgreSQL)。
type ArticleRepo interface {
	List(ctx context.Context, offset, limit int, keyword *string, sortBy, order *string, categoryID, tagID *int, status, visibility *string) ([]*entity.Article, int64, error)
	GetByID(ctx context.Context, id int64) (*entity.Article, error)
	GetBySlug(ctx context.Context, slug string) (*entity.Article, error)
	Create(ctx context.Context, article *entity.Article) (int64, error)
	Update(ctx context.Context, article *entity.Article) error
	UpdateBySlug(ctx context.Context, slug string, article *entity.Article, includeContent bool) error // 用 slug 更新文章
	Delete(ctx context.Context, id int64) error
}

// ArticleSearchRepo 文章搜索仓库 (Elasticsearch)。
type ArticleSearchRepo interface {
	Index(ctx context.Context, article *entity.Article) error
	Search(ctx context.Context, query string, offset, limit int, categoryID, tagID *int) ([]*entity.Article, int64, error)
	Get(ctx context.Context, id string) (*entity.Article, error)
	Delete(ctx context.Context, id string) error
	BulkIndex(ctx context.Context, articles []*entity.Article) error
	UpdateField(ctx context.Context, id string, field string, value interface{}) error
}

// ArticleCacheRepo 文章缓存仓库 (Redis)。
type ArticleCacheRepo interface {
	GetViews(ctx context.Context, id string) (int, error)
	IncrViews(ctx context.Context, id string) error
	GetHotList(ctx context.Context) ([]string, error)
	SetHotList(ctx context.Context, ids []string, ttl time.Duration) error
}

// ArticleLikeRepo 文章点赞仓库。
type ArticleLikeRepo interface {
	HasLiked(ctx context.Context, articleSlug, userUUID string) (bool, error)
	Toggle(ctx context.Context, articleSlug, userUUID string) (liked bool, count int32, err error)
	Remove(ctx context.Context, articleSlug, userUUID string) (removed bool, count int32, err error)
	ListUserLikedArticles(ctx context.Context, userUUID string, offset, limit int) ([]*entity.Article, int64, error)
}

// ArticleViewRepo 文章浏览记录仓库。
type ArticleViewRepo interface {
	Record(ctx context.Context, view *entity.ArticleView) error
	IncrViews(ctx context.Context, articleSlug string) error
}

// ==================== 分类/标签 ====================

// CategoryRepo 分类数据仓库。
type CategoryRepo interface {
	List(ctx context.Context, offset, limit int, keyword *string, sortBy, order *string) ([]*entity.Category, int64, error)
	ListAll(ctx context.Context) ([]*entity.Category, error)
	GetByID(ctx context.Context, id int64) (*entity.Category, error)
	Create(ctx context.Context, category entity.Category) (int64, error)
	Update(ctx context.Context, category entity.Category) error
	Delete(ctx context.Context, id int64) error
}

// TagRepo 标签数据仓库。
type TagRepo interface {
	List(ctx context.Context, offset, limit int, keyword *string, sortBy, order *string) ([]*entity.Tag, int64, error)
	ListAll(ctx context.Context) ([]*entity.Tag, error)
	ListByIDs(ctx context.Context, ids []int64) ([]*entity.Tag, error) // 根据 ID 列表获取标签
	GetByID(ctx context.Context, id int64) (*entity.Tag, error)
	Create(ctx context.Context, tag entity.Tag) (int64, error)
	Update(ctx context.Context, tag entity.Tag) error
	Delete(ctx context.Context, id int64) error
}

// ==================== 评论相关 ====================

// CommentRepo 评论数据仓库。
type CommentRepo interface {
	ListByArticleSlug(ctx context.Context, articleSlug string, offset, limit int) ([]*entity.Comment, int64, error)
	ListAll(ctx context.Context, offset, limit int, keyword, articleSlug, userUUID *string) ([]*entity.Comment, int64, error)
	GetByID(ctx context.Context, id int64) (*entity.Comment, error)
	Create(ctx context.Context, comment *entity.Comment) (int64, error)
	Delete(ctx context.Context, id int64) error
	CountByArticleSlug(ctx context.Context, articleSlug string) (int64, error)
}

// ==================== 用户相关 ====================

// UserRepo 用户数据仓库。
type UserRepo interface {
	GetByUUID(ctx context.Context, uuid string) (*entity.User, error)
	GetByID(ctx context.Context, id int64) (*entity.User, error)
	GetRoleByUUID(ctx context.Context, uuid string) (int, error)
	Create(ctx context.Context, user *entity.User) (int64, error)
	CreateFromSSO(ctx context.Context, uuid, nickname, email, avatar, address, signature string, registerSource int) error
	Update(ctx context.Context, user *entity.User) error
	List(ctx context.Context, offset, limit int, keyword, status *string) ([]*entity.User, int64, error)
	ListLogins(ctx context.Context, offset, limit int, keyword *string) ([]*entity.Login, int64, error)
	UpdateStatus(ctx context.Context, uuid string, status string) error
	GetLoginCountsByDate(ctx context.Context, days int) (map[string]int, error)
	GetRegisterCountsByDate(ctx context.Context, days int) (map[string]int, error)
}

// ==================== AI 聊天 ====================

// ChatSessionRepo AI 聊天会话仓库。
type ChatSessionRepo interface {
	Create(ctx context.Context, session *entity.ChatSession) (int64, error)
	GetByID(ctx context.Context, id int64) (*entity.ChatSession, error)
	List(ctx context.Context, userUUID string, offset, limit int) ([]*entity.ChatSession, int64, error)
	ListAll(ctx context.Context, offset, limit int, keyword, userUUID *string) ([]*entity.ChatSession, int64, error)
	Delete(ctx context.Context, id int64) error
	UpdateTitle(ctx context.Context, id int64, title string) error
}

// ChatMessageRepo AI 聊天消息仓库。
type ChatMessageRepo interface {
	Create(ctx context.Context, msg *entity.ChatMessage) (int64, error)
	ListBySessionID(ctx context.Context, sessionID int64) ([]*entity.ChatMessage, error)
	ListAll(ctx context.Context, offset, limit int, sessionID *int64, role *string) ([]*entity.ChatMessage, int64, error)
}

// ==================== 反馈 ====================

// FeedbackRepo 反馈数据仓库。
type FeedbackRepo interface {
	Create(ctx context.Context, feedback *entity.Feedback) (int64, error)
	List(ctx context.Context, offset, limit int, status *string) ([]*entity.Feedback, int64, error)
	GetByID(ctx context.Context, id int64) (*entity.Feedback, error)
	UpdateStatus(ctx context.Context, id int64, status string) error
	Delete(ctx context.Context, id int64) error
	UpdateReply(ctx context.Context, id int64, reply string) error
}

// ==================== 友链 ====================

// LinkRepo 友链数据仓库。
type LinkRepo interface {
	List(ctx context.Context) ([]*entity.Link, error)
	GetByID(ctx context.Context, id int64) (*entity.Link, error)
	Create(ctx context.Context, link *entity.Link) (int64, error)
	Update(ctx context.Context, link *entity.Link) error
	Delete(ctx context.Context, id int64) error
}

// ==================== 文件 ====================

// FileRepo 文件数据仓库。
type FileRepo interface {
	Create(ctx context.Context, file *entity.File) (int64, error)
	GetByKey(ctx context.Context, key string) (*entity.File, error)
	GetByIDs(ctx context.Context, ids []int64) ([]*entity.File, error)
	List(ctx context.Context, offset, limit int, filename, mimeType *string) ([]*entity.File, int64, error)
	Delete(ctx context.Context, key string) error
	DeleteByIDs(ctx context.Context, ids []int64) error
}

// ResourceRepo 资源数据仓库。
type ResourceRepo interface {
	Create(ctx context.Context, resource *entity.Resource) (int64, error)
	GetByID(ctx context.Context, id int64) (*entity.Resource, error)
	GetByIDs(ctx context.Context, ids []int64) ([]*entity.Resource, error)
	GetByFileHash(ctx context.Context, fileHash, userUUID string) (*entity.Resource, error)
	GetByFileHashAny(ctx context.Context, fileHash string) (*entity.Resource, error)
	List(ctx context.Context, offset, limit int, userUUID *string, filename, mimeType *string) ([]*entity.Resource, int64, error)
	Delete(ctx context.Context, id int64) error
	DeleteByIDs(ctx context.Context, ids []int64) error
	UpdateTranscodeStatus(ctx context.Context, id int64, status entity.TranscodeStatus, transcodeKey, thumbnailKey string) error
	UpdateTranscodeStatusByFileKey(ctx context.Context, fileKey string, status entity.TranscodeStatus, transcodeKey, thumbnailKey string) error
}

// ResourceUploadTaskRepo 资源上传任务数据仓库。
type ResourceUploadTaskRepo interface {
	Create(ctx context.Context, task *entity.ResourceUploadTask) (int64, error)
	GetByTaskID(ctx context.Context, taskID, userUUID string) (*entity.ResourceUploadTask, error)
	GetByFileHash(ctx context.Context, fileHash, userUUID string, statuses []entity.TaskStatus) (*entity.ResourceUploadTask, error)
	UpdateStatus(ctx context.Context, taskID string, status entity.TaskStatus) error
	UpdateChunkContext(ctx context.Context, taskID string, chunkNumber int, context string, status entity.TaskStatus) error
}

// ==================== 对象存储 ====================

// ObjectStore 对象存储仓库 (七牛云)。
type ObjectStore interface {
	Upload(ctx context.Context, key string, data io.Reader, size int64, contentType string) (string, error)
	Delete(ctx context.Context, key string) error
	GetURL(key string) string
	// 分片上传相关
	UploadBlock(ctx context.Context, data io.Reader, size int64) (string, error)
	MergeBlocks(ctx context.Context, fileSize int64, fileKey string, contexts []string) error
	GenerateFileKey(fileName, fileHash string) string
}

// ==================== AI ====================

// LLMWebAPI LLM API 调用。
type LLMWebAPI interface {
	// ChatStream 流式对话，返回内容 channel
	ChatStream(ctx context.Context, messages []LLMMessage) (<-chan string, error)
}

// LLMMessage LLM 消息。
type LLMMessage struct {
	Role    string // system, user, assistant
	Content string
}

// AIModelRepo AI 模型数据仓库。
type AIModelRepo interface {
	Create(ctx context.Context, model *entity.AIModel) (int64, error)
	GetByID(ctx context.Context, id int64) (*entity.AIModel, error)
	List(ctx context.Context, offset, limit int, name, provider *string) ([]*entity.AIModel, int64, error)
	Update(ctx context.Context, model *entity.AIModel) error
	Delete(ctx context.Context, id int64) error
}

// ==================== 表情 ====================

// EmojiRepo 表情数据仓库。
type EmojiRepo interface {
	ListActive(ctx context.Context) ([]*entity.Emoji, error)
	ListActiveByGroupKeys(ctx context.Context, groupKeys []string) ([]*entity.Emoji, error)
	ListGroups(ctx context.Context) ([]*entity.EmojiGroup, error)
	GetGroupByKey(ctx context.Context, groupKey string) (*entity.EmojiGroup, error)
	UpdateGroupSpriteConfURL(ctx context.Context, groupKey, confURL string) error
	List(ctx context.Context, offset, limit int, keyword, groupKey string, spriteGroup *int) ([]*entity.Emoji, int64, error)
}

// EmojiSpriteRepo 雪碧图数据仓库。
type EmojiSpriteRepo interface {
	ListActive(ctx context.Context) ([]*entity.EmojiSprite, error)
	DeleteAll(ctx context.Context) error
	CreateBatch(ctx context.Context, sprites []*entity.EmojiSprite) error
}

// ==================== 广告 ====================

// AdvertisementRepo 广告数据仓库。
type AdvertisementRepo interface {
	ListActive(ctx context.Context) ([]*entity.Advertisement, int64, error)
	List(ctx context.Context, offset, limit int, title *string) ([]*entity.Advertisement, int64, error)
}

// ==================== 页脚链接 ====================

// FooterLinkRepo 页脚链接数据仓库。
type FooterLinkRepo interface {
	List(ctx context.Context) ([]*entity.FooterLink, error)
}
