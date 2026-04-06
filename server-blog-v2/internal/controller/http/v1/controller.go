package v1

import (
	"github.com/go-playground/validator/v10"

	"server-blog-v2/config"
	"server-blog-v2/internal/controller/http/middleware"
	"server-blog-v2/internal/usecase"
	"server-blog-v2/pkg/logger"
)

// V1 公开 API 控制器。
type V1 struct {
	cfg            *config.Config
	logger         logger.Interface
	validate       *validator.Validate
	content        usecase.Content
	comment        usecase.Comment
	aiChat         usecase.AIChat
	feedback       usecase.Feedback
	link           usecase.Link
	file           usecase.File
	user           usecase.User
	setting        usecase.Setting
	website        usecase.Website
	emoji          usecase.Emoji
	advertisement  usecase.Advertisement
	sessionManager *middleware.SessionManager
}

// New 创建 V1 控制器。
func New(
	cfg *config.Config,
	l logger.Interface,
	content usecase.Content,
	comment usecase.Comment,
	aiChat usecase.AIChat,
	feedback usecase.Feedback,
	link usecase.Link,
	file usecase.File,
	user usecase.User,
	setting usecase.Setting,
	website usecase.Website,
	emoji usecase.Emoji,
	advertisement usecase.Advertisement,
	sessionManager *middleware.SessionManager,
) *V1 {
	return &V1{
		cfg:            cfg,
		logger:         l,
		validate:       validator.New(),
		content:        content,
		comment:        comment,
		aiChat:         aiChat,
		feedback:       feedback,
		link:           link,
		file:           file,
		user:           user,
		setting:        setting,
		website:        website,
		emoji:          emoji,
		advertisement:  advertisement,
		sessionManager: sessionManager,
	}
}


