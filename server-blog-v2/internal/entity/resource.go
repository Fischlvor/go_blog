package entity

import "time"

// TranscodeStatus 转码状态枚举。
type TranscodeStatus int8

const (
	TranscodeStatusNone       TranscodeStatus = 0 // 无需转码（非视频文件）
	TranscodeStatusProcessing TranscodeStatus = 1 // 转码中
	TranscodeStatusSuccess    TranscodeStatus = 2 // 转码成功
	TranscodeStatusFailed     TranscodeStatus = 3 // 转码失败
)

// Resource 资源实体。
type Resource struct {
	ID              int64
	FileKey         string
	FileName        string
	FileHash        string
	FileSize        int64
	MimeType        string
	UserUUID        string
	TranscodeStatus TranscodeStatus
	TranscodeKey    string
	ThumbnailKey    string
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

// TaskStatus 任务状态枚举。
type TaskStatus int8

const (
	TaskStatusInit      TaskStatus = 0 // 初始化
	TaskStatusUploading TaskStatus = 1 // 上传中
	TaskStatusCompleted TaskStatus = 2 // 已完成
	TaskStatusCancelled TaskStatus = 3 // 已取消
	TaskStatusFailed    TaskStatus = 4 // 失败
)

// ResourceUploadTask 资源上传任务实体。
type ResourceUploadTask struct {
	ID            int64
	TaskID        string
	FileName      string
	FileSize      int64
	FileHash      string
	MimeType      string
	ChunkSize     int
	TotalChunks   int
	Status        TaskStatus
	UserUUID      string
	ExpiresAt     time.Time
	QiniuContexts string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
