package admin

import (
	"github.com/go-playground/validator/v10"

	"server-blog-v2/config"
	"server-blog-v2/internal/usecase"
	"server-blog-v2/pkg/logger"
)

// Admin 管理后台 API 控制器。
type Admin struct {
	cfg           *config.Config
	logger        logger.Interface
	validate      *validator.Validate
	content       usecase.Content
	comment       usecase.Comment
	feedback      usecase.Feedback
	link          usecase.Link
	file          usecase.File
	resource      usecase.Resource
	user          usecase.User
	setting       usecase.Setting
	emoji         usecase.Emoji
	aiChat        usecase.AIChat
	aiModel       usecase.AIModel
	website       usecase.Website
	advertisement usecase.Advertisement
}

// New 创建 Admin 控制器。
func New(
	cfg *config.Config,
	l logger.Interface,
	content usecase.Content,
	comment usecase.Comment,
	feedback usecase.Feedback,
	link usecase.Link,
	file usecase.File,
	resource usecase.Resource,
	user usecase.User,
	setting usecase.Setting,
	emoji usecase.Emoji,
	aiChat usecase.AIChat,
	aiModel usecase.AIModel,
	website usecase.Website,
	advertisement usecase.Advertisement,
) *Admin {
	return &Admin{
		cfg:           cfg,
		logger:        l,
		validate:      validator.New(),
		content:       content,
		comment:       comment,
		feedback:      feedback,
		link:          link,
		file:          file,
		resource:      resource,
		user:          user,
		setting:       setting,
		emoji:         emoji,
		aiChat:        aiChat,
		aiModel:       aiModel,
		website:       website,
		advertisement: advertisement,
	}
}

