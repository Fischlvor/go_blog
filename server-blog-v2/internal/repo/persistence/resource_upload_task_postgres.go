package persistence

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"server-blog-v2/internal/entity"
	"server-blog-v2/internal/repo"

	"gorm.io/gorm"
)

// ResourceUploadTask 资源上传任务数据库模型。
type ResourceUploadTask struct {
	ID            int64  `gorm:"column:id;primaryKey;autoIncrement"`
	TaskID        string `gorm:"column:task_id;type:varchar(36);uniqueIndex"`
	FileName      string `gorm:"column:file_name;type:varchar(255)"`
	FileSize      int64  `gorm:"column:file_size"`
	FileHash      string `gorm:"column:file_hash;type:varchar(64);index"`
	MimeType      string `gorm:"column:mime_type;type:varchar(100)"`
	ChunkSize     int    `gorm:"column:chunk_size"`
	TotalChunks   int    `gorm:"column:total_chunks"`
	Status        int8   `gorm:"column:status;default:0"`
	UserUUID      string `gorm:"column:user_uuid;type:uuid;index"`
	ExpiresAt     time.Time `gorm:"column:expires_at"`
	QiniuContexts string `gorm:"column:qiniu_contexts;type:jsonb"`
	CreatedAt     time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt     time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

func (ResourceUploadTask) TableName() string {
	return "resource_upload_tasks"
}

type resourceUploadTaskRepo struct {
	db *gorm.DB
}

// NewResourceUploadTaskRepo 创建资源上传任务仓库。
func NewResourceUploadTaskRepo(db *gorm.DB) repo.ResourceUploadTaskRepo {
	return &resourceUploadTaskRepo{db: db}
}

func (r *resourceUploadTaskRepo) Create(ctx context.Context, task *entity.ResourceUploadTask) (int64, error) {
	mt := toModelTask(task)
	if err := r.db.WithContext(ctx).Create(mt).Error; err != nil {
		return 0, err
	}
	return mt.ID, nil
}

func (r *resourceUploadTaskRepo) GetByTaskID(ctx context.Context, taskID, userUUID string) (*entity.ResourceUploadTask, error) {
	var mt ResourceUploadTask
	if err := r.db.WithContext(ctx).Where("task_id = ? AND user_uuid = ?", taskID, userUUID).First(&mt).Error; err != nil {
		return nil, err
	}
	return toEntityTask(&mt), nil
}

func (r *resourceUploadTaskRepo) GetByFileHash(ctx context.Context, fileHash, userUUID string, statuses []entity.TaskStatus) (*entity.ResourceUploadTask, error) {
	var mt ResourceUploadTask
	statusInts := make([]int8, len(statuses))
	for i, s := range statuses {
		statusInts[i] = int8(s)
	}
	if err := r.db.WithContext(ctx).Where("file_hash = ? AND user_uuid = ? AND status IN ?", fileHash, userUUID, statusInts).First(&mt).Error; err != nil {
		return nil, err
	}
	return toEntityTask(&mt), nil
}

func (r *resourceUploadTaskRepo) UpdateStatus(ctx context.Context, taskID string, status entity.TaskStatus) error {
	return r.db.WithContext(ctx).Model(&ResourceUploadTask{}).Where("task_id = ?", taskID).Update("status", int8(status)).Error
}

func (r *resourceUploadTaskRepo) UpdateMimeType(ctx context.Context, taskID string, mimeType string) error {
	return r.db.WithContext(ctx).Model(&ResourceUploadTask{}).Where("task_id = ?", taskID).Updates(map[string]interface{}{
		"mime_type":  mimeType,
		"updated_at": time.Now(),
	}).Error
}

func (r *resourceUploadTaskRepo) UpdateChunkContext(ctx context.Context, taskID string, chunkNumber int, context string, status entity.TaskStatus) error {
	// 使用参数化查询防止 SQL 注入
	// 添加条件检查：只有当该块为空且任务状态允许时才更新（并发安全 + 幂等性）
	sql := `UPDATE resource_upload_tasks
            SET qiniu_contexts = jsonb_set(qiniu_contexts::jsonb, $1, to_jsonb($2::text)),
                status = $3,
                updated_at = NOW()
            WHERE task_id = $4
              AND status IN (0, 1)
              AND (qiniu_contexts->$5 IS NULL OR qiniu_contexts->$5 = 'null'::jsonb OR qiniu_contexts->>$5 = '')`

	path := fmt.Sprintf("{%d}", chunkNumber)
	chunkNumberStr := strconv.Itoa(chunkNumber)
	result := r.db.WithContext(ctx).Exec(sql, path, context, int8(status), taskID, chunkNumberStr)

	if result.Error != nil {
		return result.Error
	}

	// 检查是否真的更新了行（并发控制）
	if result.RowsAffected == 0 {
		// 可能是该块已上传或任务状态不对，查询任务状态
		var task ResourceUploadTask
		if err := r.db.WithContext(ctx).Where("task_id = ?", taskID).First(&task).Error; err != nil {
			return fmt.Errorf("查询任务失败: %w", err)
		}

		// 检查任务状态
		if task.Status != 0 && task.Status != 1 {
			return fmt.Errorf("任务状态不允许上传: %d", task.Status)
		}

		// 如果任务状态正常但未更新，说明该块已有值（并发上传或重复上传），这是正常的
		// 不报错，让上层逻辑处理
	}

	return nil
}

func toModelTask(t *entity.ResourceUploadTask) *ResourceUploadTask {
	return &ResourceUploadTask{
		ID:            t.ID,
		TaskID:        t.TaskID,
		FileName:      t.FileName,
		FileSize:      t.FileSize,
		FileHash:      t.FileHash,
		MimeType:      t.MimeType,
		ChunkSize:     t.ChunkSize,
		TotalChunks:   t.TotalChunks,
		Status:        int8(t.Status),
		UserUUID:      t.UserUUID,
		ExpiresAt:     t.ExpiresAt,
		CreatedAt:     t.CreatedAt,
		UpdatedAt:     t.UpdatedAt,
		QiniuContexts: t.QiniuContexts,
	}
}

func toEntityTask(mt *ResourceUploadTask) *entity.ResourceUploadTask {
	return &entity.ResourceUploadTask{
		ID:            mt.ID,
		TaskID:        mt.TaskID,
		FileName:      mt.FileName,
		FileSize:      mt.FileSize,
		FileHash:      mt.FileHash,
		MimeType:      mt.MimeType,
		ChunkSize:     mt.ChunkSize,
		TotalChunks:   mt.TotalChunks,
		Status:        entity.TaskStatus(mt.Status),
		UserUUID:      mt.UserUUID,
		ExpiresAt:     mt.ExpiresAt,
		QiniuContexts: mt.QiniuContexts,
		CreatedAt:     mt.CreatedAt,
		UpdatedAt:     mt.UpdatedAt,
	}
}
