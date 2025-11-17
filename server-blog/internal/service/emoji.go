package service

import (
	"context"
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"

	"server/internal/model"
	"server/pkg/global"
	"server/pkg/utils"
	"server/pkg/utils/upload"
)

type EmojiService struct{}

// GetEmojiList 获取表情列表
func (s *EmojiService) GetEmojiList(page, pageSize int, groupKey, keyword string) (*model.EmojiListResponse, error) {
	var emojis []model.Emoji
	var total int64

	db := global.DB.Model(&model.Emoji{}).Where("status = ?", model.EmojiStatusActive)

	// 按组筛选
	if groupKey != "" {
		db = db.Where("group_key = ?", groupKey)
	}

	// 关键词搜索
	if keyword != "" {
		db = db.Where("(`key` LIKE ? OR filename LIKE ?)", "%"+keyword+"%", "%"+keyword+"%")
	}

	// 获取总数
	if err := db.Count(&total).Error; err != nil {
		return nil, err
	}

	// 分页查询
	offset := (page - 1) * pageSize
	if err := db.Order("id ASC").Offset(offset).Limit(pageSize).Find(&emojis).Error; err != nil {
		return nil, err
	}

	// 处理CDN URL拼接
	for i := range emojis {
		emojis[i].CdnUrl = utils.PublicURLFromDB(emojis[i].CdnUrl)
	}

	return &model.EmojiListResponse{
		List:     emojis,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}, nil
}

// GetEmojiGroups 获取表情组列表
func (s *EmojiService) GetEmojiGroups() ([]model.EmojiGroup, error) {
	var groups []model.EmojiGroup

	// 更新每个组的emoji数量
	if err := s.updateGroupEmojiCount(); err != nil {
		return nil, err
	}

	if err := global.DB.Where("status = ?", 1).Order("sort_order ASC, id ASC").Find(&groups).Error; err != nil {
		return nil, err
	}

	return groups, nil
}

// CreateEmojiGroup 创建表情组
func (s *EmojiService) CreateEmojiGroup(groupName, description string, createdBy uuid.UUID) error {
	group := &model.EmojiGroup{
		GroupName:   groupName,
		Description: description,
		SortOrder:   0,
		EmojiCount:  0,
		Status:      1,
		CreatedBy:   createdBy,
	}

	return global.DB.Create(group).Error
}

// UpdateEmojiGroup 更新表情组
func (s *EmojiService) UpdateEmojiGroup(id uint64, groupName, description string) error {
	return global.DB.Model(&model.EmojiGroup{}).Where("id = ?", id).Updates(map[string]interface{}{
		"group_name":  groupName,
		"description": description,
		"updated_at":  time.Now(),
	}).Error
}

// DeleteEmojiGroup 删除表情组（需要确保组内无表情）
func (s *EmojiService) DeleteEmojiGroup(id uint64) error {
	// 检查组内是否有表情
	var count int64
	if err := global.DB.Model(&model.Emoji{}).Where("group_key = (SELECT group_key FROM emoji_groups WHERE id = ?) AND status = ?", id, model.EmojiStatusActive).Count(&count).Error; err != nil {
		return err
	}

	if count > 0 {
		return fmt.Errorf("组内还有 %d 个表情，无法删除", count)
	}

	return global.DB.Delete(&model.EmojiGroup{}, id).Error
}

// DeleteEmoji 软删除表情
func (s *EmojiService) DeleteEmoji(id uint64) error {
	now := time.Now()
	return global.DB.Model(&model.Emoji{}).Where("id = ?", id).Updates(map[string]interface{}{
		"status":     model.EmojiStatusDeleted,
		"deleted_at": &now,
		"updated_at": now,
	}).Error
}

// RestoreEmoji 恢复表情
func (s *EmojiService) RestoreEmoji(id uint64) error {
	return global.DB.Model(&model.Emoji{}).Where("id = ?", id).Updates(map[string]interface{}{
		"status":     model.EmojiStatusActive,
		"deleted_at": nil,
		"updated_at": time.Now(),
	}).Error
}

// RegenerateSprites 重新生成雪碧图
func (s *EmojiService) RegenerateSprites(createdBy uuid.UUID) (*model.EmojiTask, error) {
	// 创建任务记录
	task := &model.EmojiTask{
		TaskType:  model.TaskTypeRegenerateSprites,
		Status:    model.TaskStatusPending,
		Progress:  0,
		Message:   "准备重新生成雪碧图...",
		Params:    "{}", // 提供空的JSON对象而不是空字符串
		Result:    "{}", // 提供空的JSON对象而不是空字符串
		CreatedBy: createdBy,
	}

	if err := global.DB.Create(task).Error; err != nil {
		return nil, err
	}

	// 异步执行任务
	go s.executeRegenerateSprites(task.ID)

	return task, nil
}

// GetTaskStatus 获取任务状态
func (s *EmojiService) GetTaskStatus(taskId uint64) (*model.EmojiTask, error) {
	var task model.EmojiTask
	if err := global.DB.First(&task, taskId).Error; err != nil {
		return nil, err
	}
	return &task, nil
}

// GetEmojiConfig 获取前端配置
func (s *EmojiService) GetEmojiConfig() (*model.EmojiConfigResponse, error) {
	// 获取所有活跃的emoji
	var emojis []model.Emoji
	if err := global.DB.Where("status = ?", model.EmojiStatusActive).Order("id ASC").Find(&emojis).Error; err != nil {
		return nil, err
	}

	// 获取雪碧图信息
	var sprites []model.EmojiSprite
	if err := global.DB.Where("status = ?", 1).Order("sprite_group ASC").Find(&sprites).Error; err != nil {
		return nil, err
	}

	// 构建雪碧图信息
	var spriteInfos []model.EmojiSpriteInfo
	for _, sprite := range sprites {
		// 计算该雪碧图包含的emoji范围
		startIdx := sprite.SpriteGroup * 128
		endIdx := startIdx + sprite.EmojiCount - 1

		spriteInfos = append(spriteInfos, model.EmojiSpriteInfo{
			ID:       sprite.SpriteGroup,
			Filename: sprite.Filename,
			URL:      utils.PublicURLFromDB(sprite.CdnUrl),
			Range:    [2]int{startIdx, endIdx},
			Frozen:   true,                                        // 所有雪碧图都标记为冻结
			Size:     [2]int{sprite.Width / 2, sprite.Height / 2}, // 显示尺寸为实际尺寸的一半
		})
	}

	return &model.EmojiConfigResponse{
		Version:     fmt.Sprintf("v1.%d", time.Now().Unix()),
		TotalEmojis: int64(len(emojis)),
		Sprites:     spriteInfos,
		Mapping:     nil, // 不再需要 mapping，前端直接使用 e123 格式
		UpdatedAt:   time.Now(),
	}, nil
}

// GetSpriteList 获取雪碧图列表
func (s *EmojiService) GetSpriteList() ([]model.EmojiSprite, error) {
	var sprites []model.EmojiSprite

	err := global.DB.Where("status = ?", model.EmojiSpriteStatusActive).
		Order("sprite_group ASC").
		Find(&sprites).Error

	if err != nil {
		return nil, err
	}

	// 处理CDN URL拼接
	for i := range sprites {
		sprites[i].CdnUrl = utils.PublicURLFromDB(sprites[i].CdnUrl)
	}

	return sprites, nil
}

// 私有方法

// getNextEmojiKey 获取下一个emoji key
func (s *EmojiService) getNextEmojiKey() (string, error) {
	var maxKey string
	if err := global.DB.Model(&model.Emoji{}).Select("MAX(`key`)").Scan(&maxKey).Error; err != nil {
		return "", err
	}

	if maxKey == "" {
		return "e0", nil
	}

	// 解析最大key
	numStr := strings.TrimPrefix(maxKey, "e")
	num, err := strconv.Atoi(numStr)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("e%d", num+1), nil
}

// updateGroupEmojiCount 更新组的emoji数量
func (s *EmojiService) updateGroupEmojiCount() error {
	return global.DB.Exec(`
		UPDATE emoji_groups 
		SET emoji_count = (
			SELECT COUNT(*) 
			FROM emojis 
			WHERE emojis.group_key = emoji_groups.group_key 
			AND emojis.status = ?
		)
	`, model.EmojiStatusActive).Error
}

// DoRegenerateSpritesWithSSE 直接执行雪碧图生成并通过SSE流式返回进度
func (s *EmojiService) DoRegenerateSpritesWithSSE(userUUID uuid.UUID, groupKeys []string, sseWriter SSEWriter) error {
	// 使用临时的 taskId（0表示非数据库任务）
	return s.doRegenerateSprites(0, groupKeys, sseWriter)
}

// regenerateSprites 重新生成雪碧图（异步）
func (s *EmojiService) regenerateSprites() {
	// 创建雪碧图生成任务
	task := &model.EmojiTask{
		TaskType:  model.TaskTypeRegenerateSprites,
		Status:    model.TaskStatusPending,
		Progress:  0,
		Message:   "等待生成雪碧图...",
		Params:    "{}",                                                         // 提供空的JSON对象而不是空字符串
		Result:    "{}",                                                         // 提供空的JSON对象而不是空字符串
		CreatedBy: uuid.FromStringOrNil("37395c61-a2ec-464e-9567-ce6fa92630f7"), // 系统任务
	}

	if err := global.DB.Create(task).Error; err != nil {
		fmt.Printf("创建雪碧图生成任务失败: %v\n", err)
		return
	}

	fmt.Printf("已创建雪碧图生成任务，任务ID: %d\n", task.ID)

	// 异步执行任务
	go s.executeRegenerateSprites(task.ID)
}

// DBTaskSSEWriter 数据库任务的SSE写入器实现
type DBTaskSSEWriter struct {
	taskId uint64
}

func (w *DBTaskSSEWriter) WriteSSEEvent(event string, data interface{}) error {
	// 对于数据库任务，我们只更新进度信息
	if event == "progress" || event == "group_start" || event == "sprite_progress" {
		if dataMap, ok := data.(map[string]interface{}); ok {
			message := ""
			if msg, ok := dataMap["message"]; ok {
				message = fmt.Sprintf("%v", msg)
			}
			if message != "" {
				global.DB.Model(&model.EmojiTask{}).Where("id = ?", w.taskId).Update("message", message)
			}
		}
	}
	return nil
}

func (w *DBTaskSSEWriter) Flush() {
	// 数据库写入器不需要flush
}

// executeRegenerateSprites 执行重新生成雪碧图任务
func (s *EmojiService) executeRegenerateSprites(taskId uint64) {
	defer func() {
		if r := recover(); r != nil {
			// 任务执行失败
			now := time.Now()
			global.DB.Model(&model.EmojiTask{}).Where("id = ?", taskId).Updates(map[string]interface{}{
				"status":       model.TaskStatusFailed,
				"message":      fmt.Sprintf("任务执行失败: %v", r),
				"completed_at": &now,
			})
		}
	}()

	// 更新任务状态为运行中
	global.DB.Model(&model.EmojiTask{}).Where("id = ?", taskId).Updates(map[string]interface{}{
		"status":     model.TaskStatusRunning,
		"progress":   0,
		"message":    "正在生成雪碧图...",
		"started_at": time.Now(),
	})

	// 创建数据库任务的SSE写入器
	sseWriter := &DBTaskSSEWriter{taskId: taskId}

	// 执行雪碧图生成（空的groupKeys表示生成所有）
	err := s.doRegenerateSprites(taskId, []string{}, sseWriter)

	if err != nil {
		// 任务失败
		now := time.Now()
		global.DB.Model(&model.EmojiTask{}).Where("id = ?", taskId).Updates(map[string]interface{}{
			"status":       model.TaskStatusFailed,
			"progress":     0,
			"message":      fmt.Sprintf("雪碧图生成失败: %v", err),
			"completed_at": &now,
		})
		return
	}

	// 完成任务
	now := time.Now()
	global.DB.Model(&model.EmojiTask{}).Where("id = ?", taskId).Updates(map[string]interface{}{
		"status":       model.TaskStatusCompleted,
		"progress":     100,
		"message":      "雪碧图生成完成",
		"result":       `{"success": true, "sprites_generated": true}`,
		"completed_at": &now,
	})
}

// SSEWriter SSE响应写入器接口
type SSEWriter interface {
	WriteSSEEvent(event string, data interface{}) error
	Flush()
}

// doRegenerateSprites 实际执行雪碧图生成
func (s *EmojiService) doRegenerateSprites(taskId uint64, groupKeys []string, sseWriter SSEWriter) error {
	// 1. 构建查询条件
	var emojis []model.Emoji
	db := global.DB.Where("status = ?", model.EmojiStatusActive).Order("group_key ASC, id ASC")

	// 如果指定了 groupKeys，则只查询这些 group 的表情
	if len(groupKeys) > 0 {
		db = db.Where("group_key IN ?", groupKeys)
	}

	if err := db.Find(&emojis).Error; err != nil {
		return fmt.Errorf("获取表情列表失败: %v", err)
	}

	if len(emojis) == 0 {
		return fmt.Errorf("没有活跃的表情")
	}

	// 2. 按 group_key 分组
	targetGroups := make(map[string][]model.Emoji)
	for _, emoji := range emojis {
		targetGroups[emoji.GroupKey] = append(targetGroups[emoji.GroupKey], emoji)
	}

	if len(targetGroups) == 0 {
		return fmt.Errorf("没有找到要处理的表情组")
	}

	// 4. 发送开始事件
	sseWriter.WriteSSEEvent("start", gin.H{
		"total_groups": len(targetGroups),
		"message":      "开始生成雪碧图...",
	})

	// 5. 创建临时目录存放生成的雪碧图
	tempDir := filepath.Join(os.TempDir(), fmt.Sprintf("emoji_sprites_%d", taskId))
	if err := os.MkdirAll(tempDir, 0755); err != nil {
		return fmt.Errorf("创建临时目录失败: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// 6. 为每个组生成雪碧图
	oss := upload.NewOss()
	spriteResults := make(map[string][]*model.EmojiSprite) // group_key -> []*EmojiSprite

	totalGroups := len(targetGroups)
	processedGroups := 0

	for groupKey, groupEmojis := range targetGroups {
		processedGroups++

		// 获取 group_name
		var group model.EmojiGroup
		global.DB.Where("group_key = ?", groupKey).First(&group)

		// 发送组处理开始事件
		sseWriter.WriteSSEEvent("group_start", gin.H{
			"group_key":  groupKey,
			"group_name": group.GroupName,
			"current":    processedGroups,
			"total":      totalGroups,
			"message":    fmt.Sprintf("正在处理表情组: %s (%d/%d)...", group.GroupName, processedGroups, totalGroups),
		})

		// 下载表情图片并生成雪碧图
		sprites, err := s.generateSpritesForGroupWithSSE(groupKey, group.GroupName, groupEmojis, tempDir, oss, sseWriter)
		if err != nil {
			sseWriter.WriteSSEEvent("error", gin.H{
				"message": fmt.Sprintf("生成表情组 %s 的雪碧图失败: %v", group.GroupName, err),
			})
			return err
		}

		spriteResults[groupKey] = sprites

		// 发送组处理完成事件
		sseWriter.WriteSSEEvent("group_complete", gin.H{
			"group_key":    groupKey,
			"group_name":   group.GroupName,
			"current":      processedGroups,
			"total":        totalGroups,
			"sprite_count": len(sprites),
		})
	}

	sseWriter.WriteSSEEvent("progress", gin.H{
		"message": "正在清理旧数据...",
	})

	// 7. 清理旧的雪碧图记录
	if err := global.DB.Where("1=1").Delete(&model.EmojiSprite{}).Error; err != nil {
		return fmt.Errorf("清理旧雪碧图记录失败: %v", err)
	}

	// 8. 上传新的雪碧图并保存记录
	allSprites := make([]model.EmojiSprite, 0)
	for _, sprites := range spriteResults {
		for _, sprite := range sprites {
			allSprites = append(allSprites, *sprite)
		}
	}

	if err := global.DB.CreateInBatches(allSprites, 100).Error; err != nil {
		return fmt.Errorf("保存雪碧图记录失败: %v", err)
	}

	sseWriter.WriteSSEEvent("progress", gin.H{
		"message": "正在生成配置文件...",
	})

	// 9. 为每个组生成配置文件并上传
	for groupKey, sprites := range spriteResults {
		var group model.EmojiGroup
		global.DB.Where("group_key = ?", groupKey).First(&group)

		confPath, err := s.generateAndUploadSpriteConfig(groupKey, sprites, oss)
		if err != nil {
			return fmt.Errorf("生成配置文件失败: %v", err)
		}

		// 更新 emoji_groups 的 sprite_conf_url
		if err := global.DB.Model(&model.EmojiGroup{}).Where("group_key = ?", groupKey).Update("sprite_conf_url", confPath).Error; err != nil {
			return fmt.Errorf("更新配置URL失败: %v", err)
		}
	}

	sseWriter.WriteSSEEvent("complete", gin.H{
		"message":       "雪碧图生成完成",
		"groups_count":  len(spriteResults),
		"sprites_count": len(allSprites),
	})

	return nil
}

// generateSpritesForGroupWithSSE 为一个表情组生成雪碧图（带SSE进度报告）
func (s *EmojiService) generateSpritesForGroupWithSSE(groupKey string, groupName string, emojis []model.Emoji, tempDir string, oss upload.OSS, sseWriter SSEWriter) ([]*model.EmojiSprite, error) {
	const (
		targetSize      = 64  // 每个emoji的目标尺寸
		spritesPerRow   = 16  // 每行16个emoji
		emojisPerSprite = 128 // 每个雪碧图128个emoji
	)

	sprites := make([]*model.EmojiSprite, 0)
	totalSprites := (len(emojis) + emojisPerSprite - 1) / emojisPerSprite

	// 按 emojisPerSprite 分组
	for spriteSeq := 0; spriteSeq*emojisPerSprite < len(emojis); spriteSeq++ {
		startIdx := spriteSeq * emojisPerSprite
		endIdx := startIdx + emojisPerSprite
		if endIdx > len(emojis) {
			endIdx = len(emojis)
		}

		spriteEmojis := emojis[startIdx:endIdx]

		// 发送雪碧图处理进度
		sseWriter.WriteSSEEvent("sprite_progress", gin.H{
			"group_key":  groupKey,
			"group_name": groupName,
			"sprite_seq": spriteSeq,
			"total":      totalSprites,
			"current":    spriteSeq + 1,
			"message":    fmt.Sprintf("正在生成雪碧图 %d/%d...", spriteSeq+1, totalSprites),
		})

		// 计算雪碧图尺寸
		rows := (len(spriteEmojis) + spritesPerRow - 1) / spritesPerRow
		spriteWidth := spritesPerRow * targetSize
		spriteHeight := rows * targetSize

		// 创建空白雪碧图
		spriteImg := image.NewRGBA(image.Rect(0, 0, spriteWidth, spriteHeight))

		// 下载并粘贴每个emoji
		for i, emoji := range spriteEmojis {
			row := i / spritesPerRow
			col := i % spritesPerRow

			x := col * targetSize
			y := row * targetSize

			// 下载表情图片
			emojiImg, err := s.downloadAndResizeEmoji(emoji.CdnUrl, targetSize)
			if err != nil {
				return nil, fmt.Errorf("下载表情 %s 失败: %v", emoji.Key, err)
			}

			// 粘贴到雪碧图
			draw.Draw(spriteImg, image.Rect(x, y, x+targetSize, y+targetSize), emojiImg, image.Point{}, draw.Over)
		}

		// 保存雪碧图
		spriteFilename := fmt.Sprintf("sprite_%s_%d.png", groupKey, spriteSeq)
		spritePath := filepath.Join(tempDir, spriteFilename)

		file, err := os.Create(spritePath)
		if err != nil {
			return nil, fmt.Errorf("创建雪碧图文件失败: %v", err)
		}
		defer file.Close()

		if err := png.Encode(file, spriteImg); err != nil {
			return nil, fmt.Errorf("编码雪碧图失败: %v", err)
		}

		// 上传到七牛云
		pathPrefix := fmt.Sprintf("emoji/sprite")
		uploadedPath, err := s.uploadFileToOSS(spritePath, pathPrefix, oss)
		if err != nil {
			return nil, fmt.Errorf("上传雪碧图失败: %v", err)
		}

		// 获取文件大小
		fileInfo, err := os.Stat(spritePath)
		if err != nil {
			return nil, fmt.Errorf("获取雪碧图文件大小失败: %v", err)
		}

		// 创建 EmojiSprite 记录
		sprite := &model.EmojiSprite{
			SpriteGroup: spriteSeq,
			Filename:    spriteFilename,
			CdnUrl:      uploadedPath,
			Width:       spriteWidth,
			Height:      spriteHeight,
			EmojiCount:  len(spriteEmojis),
			FileSize:    int(fileInfo.Size()),
			Status:      model.EmojiSpriteStatusActive,
		}

		sprites = append(sprites, sprite)
	}

	return sprites, nil
}

// generateSpritesForGroup 为一个表情组生成雪碧图
func (s *EmojiService) generateSpritesForGroup(groupKey string, emojis []model.Emoji, tempDir string, oss upload.OSS) ([]*model.EmojiSprite, error) {
	const (
		targetSize      = 64  // 每个emoji的目标尺寸
		spritesPerRow   = 16  // 每行16个emoji
		emojisPerSprite = 128 // 每个雪碧图128个emoji
	)

	sprites := make([]*model.EmojiSprite, 0)

	// 按 emojisPerSprite 分组
	for spriteSeq := 0; spriteSeq*emojisPerSprite < len(emojis); spriteSeq++ {
		startIdx := spriteSeq * emojisPerSprite
		endIdx := startIdx + emojisPerSprite
		if endIdx > len(emojis) {
			endIdx = len(emojis)
		}

		spriteEmojis := emojis[startIdx:endIdx]

		// 计算雪碧图尺寸
		rows := (len(spriteEmojis) + spritesPerRow - 1) / spritesPerRow
		spriteWidth := spritesPerRow * targetSize
		spriteHeight := rows * targetSize

		// 创建空白雪碧图
		spriteImg := image.NewRGBA(image.Rect(0, 0, spriteWidth, spriteHeight))

		// 下载并粘贴每个emoji
		for i, emoji := range spriteEmojis {
			row := i / spritesPerRow
			col := i % spritesPerRow

			x := col * targetSize
			y := row * targetSize

			// 下载表情图片
			emojiImg, err := s.downloadAndResizeEmoji(emoji.CdnUrl, targetSize)
			if err != nil {
				return nil, fmt.Errorf("下载表情 %s 失败: %v", emoji.Key, err)
			}

			// 粘贴到雪碧图
			draw.Draw(spriteImg, image.Rect(x, y, x+targetSize, y+targetSize), emojiImg, image.Point{}, draw.Over)
		}

		// 保存雪碧图
		spriteFilename := fmt.Sprintf("sprite_%s_%d.png", groupKey, spriteSeq)
		spritePath := filepath.Join(tempDir, spriteFilename)

		file, err := os.Create(spritePath)
		if err != nil {
			return nil, fmt.Errorf("创建雪碧图文件失败: %v", err)
		}
		defer file.Close()

		if err := png.Encode(file, spriteImg); err != nil {
			return nil, fmt.Errorf("编码雪碧图失败: %v", err)
		}

		// 上传到七牛云
		fileHeader, err := os.Open(spritePath)
		if err != nil {
			return nil, fmt.Errorf("打开雪碧图文件失败: %v", err)
		}
		defer fileHeader.Close()

		// 使用 OSS 上传（需要转换为 multipart.FileHeader）
		// 这里我们直接上传二进制数据
		pathPrefix := fmt.Sprintf("emoji/sprite")
		uploadedPath, err := s.uploadFileToOSS(spritePath, pathPrefix, oss)
		if err != nil {
			return nil, fmt.Errorf("上传雪碧图失败: %v", err)
		}

		// 获取文件大小
		fileInfo, err := os.Stat(spritePath)
		if err != nil {
			return nil, fmt.Errorf("获取雪碧图文件大小失败: %v", err)
		}

		// 创建 EmojiSprite 记录
		sprite := &model.EmojiSprite{
			SpriteGroup: spriteSeq,
			Filename:    spriteFilename,
			CdnUrl:      uploadedPath,
			Width:       spriteWidth,
			Height:      spriteHeight,
			EmojiCount:  len(spriteEmojis),
			FileSize:    int(fileInfo.Size()),
			Status:      model.EmojiSpriteStatusActive,
		}

		sprites = append(sprites, sprite)
	}

	return sprites, nil
}

// downloadAndResizeEmoji 下载并调整emoji尺寸
func (s *EmojiService) downloadAndResizeEmoji(cdnUrl string, targetSize int) (image.Image, error) {
	// 拼接完整URL
	fullURL := utils.PublicURLFromDB(cdnUrl)

	resp, err := http.Get(fullURL)
	if err != nil {
		return nil, fmt.Errorf("下载emoji失败: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("下载emoji返回状态码: %d", resp.StatusCode)
	}

	// 解码图片
	img, err := png.Decode(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("解码emoji图片失败: %v", err)
	}

	// 获取原始图片尺寸
	bounds := img.Bounds()
	srcWidth := bounds.Dx()
	srcHeight := bounds.Dy()

	// 创建目标大小的图片
	resized := image.NewRGBA(image.Rect(0, 0, targetSize, targetSize))

	// 直接缩放到目标分辨率
	scaleX := float64(targetSize) / float64(srcWidth)
	scaleY := float64(targetSize) / float64(srcHeight)

	for y := 0; y < targetSize; y++ {
		for x := 0; x < targetSize; x++ {
			// 计算源图片中的坐标
			srcX := bounds.Min.X + int(float64(x)/scaleX)
			srcY := bounds.Min.Y + int(float64(y)/scaleY)

			// 确保坐标在有效范围内
			if srcX >= bounds.Min.X && srcX < bounds.Max.X && srcY >= bounds.Min.Y && srcY < bounds.Max.Y {
				r, g, b, a := img.At(srcX, srcY).RGBA()
				resized.SetRGBA(x, y, color.RGBA{
					R: uint8(r >> 8),
					G: uint8(g >> 8),
					B: uint8(b >> 8),
					A: uint8(a >> 8),
				})
			}
		}
	}

	return resized, nil
}

// uploadFileToOSS 上传文件到OSS
func (s *EmojiService) uploadFileToOSS(filePath string, pathPrefix string, oss upload.OSS) (string, error) {
	// 根据 OSS 类型选择上传方式
	if global.Config.System.OssType == "qiniu" {
		return s.uploadFileToQiniu(filePath, pathPrefix)
	}
	// 本地存储
	return s.uploadFileToLocal(filePath, pathPrefix)
}

// uploadFileToQiniu 上传文件到七牛云
func (s *EmojiService) uploadFileToQiniu(filePath string, pathPrefix string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("打开文件失败: %v", err)
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return "", fmt.Errorf("获取文件信息失败: %v", err)
	}

	// 生成文件 key
	filename := filepath.Base(filePath)
	ext := filepath.Ext(filename)
	name := strings.TrimSuffix(filename, ext)
	fileKey := fmt.Sprintf("%s/%s-%d%s", pathPrefix, name, time.Now().UnixNano(), ext)

	// 上传到七牛云
	putPolicy := storage.PutPolicy{Scope: global.Config.Qiniu.Bucket}
	mac := qbox.NewMac(global.Config.Qiniu.AccessKey, global.Config.Qiniu.SecretKey)
	upToken := putPolicy.UploadToken(mac)
	cfg := &storage.Config{
		UseHTTPS:      true,
		UseCdnDomains: true,
	}
	formUploader := storage.NewFormUploader(cfg)
	putRet := storage.PutRet{}
	putExtra := storage.PutExtra{Params: map[string]string{}}

	err = formUploader.Put(context.Background(), &putRet, upToken, fileKey, file, fileInfo.Size(), &putExtra)
	if err != nil {
		return "", fmt.Errorf("上传到七牛云失败: %v", err)
	}

	return putRet.Key, nil
}

// uploadFileToLocal 上传文件到本地存储
func (s *EmojiService) uploadFileToLocal(filePath string, pathPrefix string) (string, error) {
	// 本地存储：复制文件到指定目录
	uploadDir := filepath.Join(global.Config.Upload.Path, pathPrefix)
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		return "", fmt.Errorf("创建上传目录失败: %v", err)
	}

	filename := filepath.Base(filePath)
	destPath := filepath.Join(uploadDir, filename)

	// 复制文件
	source, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("打开源文件失败: %v", err)
	}
	defer source.Close()

	dest, err := os.Create(destPath)
	if err != nil {
		return "", fmt.Errorf("创建目标文件失败: %v", err)
	}
	defer dest.Close()

	if _, err := io.Copy(dest, source); err != nil {
		return "", fmt.Errorf("复制文件失败: %v", err)
	}

	// 返回相对路径
	return filepath.Join(pathPrefix, filename), nil
}

// generateAndUploadSpriteConfig 生成并上传雪碧图配置文件
func (s *EmojiService) generateAndUploadSpriteConfig(groupKey string, sprites []*model.EmojiSprite, oss upload.OSS) (string, error) {
	// 获取该组的所有表情
	var emojis []model.Emoji
	if err := global.DB.Where("group_key = ? AND status = ?", groupKey, model.EmojiStatusActive).Order("id ASC").Find(&emojis).Error; err != nil {
		return "", fmt.Errorf("获取表情失败: %v", err)
	}

	// 构建配置结构
	config := map[string]interface{}{
		"group_key":  groupKey,
		"version":    fmt.Sprintf("v1.%d", time.Now().Unix()),
		"updated_at": time.Now().Format(time.RFC3339),
		"sprites":    sprites,
		"emojis":     emojis,
	}

	// 转换为JSON
	configJSON, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return "", fmt.Errorf("生成配置JSON失败: %v", err)
	}

	// 保存到临时文件
	tempDir := os.TempDir()
	confFilename := fmt.Sprintf("sprite_%s.conf", groupKey)
	confPath := filepath.Join(tempDir, confFilename)

	if err := os.WriteFile(confPath, configJSON, 0644); err != nil {
		return "", fmt.Errorf("写入配置文件失败: %v", err)
	}

	// 上传到七牛云
	uploadedPath, err := s.uploadFileToOSS(confPath, "emoji/sprite", oss)
	if err != nil {
		return "", fmt.Errorf("上传配置文件失败: %v", err)
	}

	return uploadedPath, nil
}

// updateTaskProgress 更新任务进度
func (s *EmojiService) updateTaskProgress(taskId uint64, progress int, message string) {
	global.DB.Model(&model.EmojiTask{}).Where("id = ?", taskId).Updates(map[string]interface{}{
		"progress": progress,
		"message":  message,
	})
}
