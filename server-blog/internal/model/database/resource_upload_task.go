package database

import (
	"server/pkg/global"
	"time"

	"github.com/gofrs/uuid"
)

// TaskStatus 任务状态枚举
type TaskStatus int8

const (
	TaskStatusInit      TaskStatus = 0 // 初始化
	TaskStatusUploading TaskStatus = 1 // 上传中
	TaskStatusCompleted TaskStatus = 2 // 已完成
	TaskStatusCancelled TaskStatus = 3 // 已取消
	TaskStatusFailed    TaskStatus = 4 // 失败
)

// ResourceUploadTask 资源上传任务表
// 用于存储分片上传的临时状态，上传完成后可清理
type ResourceUploadTask struct {
	global.MODEL
	TaskID        string     `json:"task_id" gorm:"size:64;uniqueIndex;not null;comment:任务ID(UUID)"`
	FileName      string     `json:"file_name" gorm:"size:255;not null;comment:文件名"`
	FileSize      int64      `json:"file_size" gorm:"not null;comment:文件大小(字节)"`
	FileHash      string     `json:"file_hash" gorm:"size:64;index;comment:文件MD5(前端计算)"`
	MimeType      string     `json:"mime_type" gorm:"size:100;not null;comment:MIME类型"`
	ChunkSize     int        `json:"chunk_size" gorm:"default:4194304;comment:块大小(默认4MB)"`
	TotalChunks   int        `json:"total_chunks" gorm:"not null;comment:总块数"`
	Status        TaskStatus `json:"status" gorm:"default:0;comment:状态(0初始化 1上传中 2已完成 3已取消 4失败)"`
	UserUUID      uuid.UUID  `json:"user_uuid" gorm:"type:char(36);index;comment:用户UUID"`
	User          User       `json:"-" gorm:"foreignKey:UserUUID;references:UUID"`
	ExpiresAt     time.Time  `json:"expires_at" gorm:"index;comment:过期时间"`
	QiniuContexts string     `json:"qiniu_contexts" gorm:"type:json;comment:七牛云Context数组(JSON)"`
	// QiniuContexts 格式说明：
	// - 初始化时：["", "", "", ...] (长度=TotalChunks)
	// - 上传块后：["ctx0", "ctx1", "", "ctx3", ...] (非空=已上传)
	// - 并发更新：使用 MySQL JSON_SET 原子操作
	//   UPDATE ... SET qiniu_contexts = JSON_SET(qiniu_contexts, '$[2]', 'ctx2')
}

// TableName 表名
func (ResourceUploadTask) TableName() string {
	return "resource_upload_tasks"
}
