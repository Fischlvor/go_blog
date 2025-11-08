package database

import (
	"server/pkg/global"

	"github.com/gofrs/uuid"
)

// ArticleLike 文章收藏表
type ArticleLike struct {
	global.MODEL
	ArticleID string    `json:"article_id"`                     // 文章 ID
	UserUUID  uuid.UUID `json:"user_uuid" gorm:"type:char(36)"` // 用户 UUID
	User      User      `json:"-" gorm:"foreignKey:UserUUID;references:UUID"`
}
