package database

import (
	"server/pkg/global"

	"github.com/gofrs/uuid"
)

// TranscodeStatus 转码状态枚举
type TranscodeStatus int8

const (
	TranscodeStatusNone       TranscodeStatus = 0 // 无需转码（非视频文件）
	TranscodeStatusProcessing TranscodeStatus = 1 // 转码中
	TranscodeStatusSuccess    TranscodeStatus = 2 // 转码成功
	TranscodeStatusFailed     TranscodeStatus = 3 // 转码失败
)

// Resource 资源表
// 存储已上传完成的资源信息，永久保存
// 秒传说明：相同 FileHash 的文件只存一份物理文件，但每个用户有独立记录
type Resource struct {
	global.MODEL
	FileKey         string          `json:"file_key" gorm:"size:500;not null;comment:七牛云Key(不含域名)"`
	FileName        string          `json:"file_name" gorm:"size:255;not null;comment:文件名"`
	FileHash        string          `json:"file_hash" gorm:"size:64;index;comment:文件MD5(秒传检测用)"`
	FileSize        int64           `json:"file_size" gorm:"not null;comment:文件大小(字节)"`
	MimeType        string          `json:"mime_type" gorm:"size:100;not null;comment:MIME类型"`
	UserUUID        uuid.UUID       `json:"user_uuid" gorm:"type:char(36);index;comment:用户UUID"`
	User            User            `json:"-" gorm:"foreignKey:UserUUID;references:UUID"`
	TranscodeStatus TranscodeStatus `json:"transcode_status" gorm:"default:0;comment:转码状态:0无需转码,1转码中,2成功,3失败"`
	TranscodeKey    string          `json:"transcode_key" gorm:"size:500;comment:转码后视频Key(不含域名)"`
	ThumbnailKey    string          `json:"thumbnail_key" gorm:"size:500;comment:缩略图Key(不含域名)"`
}

// TableName 表名
func (Resource) TableName() string {
	return "resources"
}

// 联合唯一索引：同一用户不能重复上传相同文件
// 需要在迁移时手动创建：
// CREATE UNIQUE INDEX idx_user_hash ON resources(user_uuid, file_hash)
