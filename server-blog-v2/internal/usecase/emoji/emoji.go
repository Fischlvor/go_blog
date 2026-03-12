package emoji

import (
	"context"
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"

	"server-blog-v2/config"
	"server-blog-v2/internal/entity"
	"server-blog-v2/internal/repo"
	"server-blog-v2/internal/usecase"
	"server-blog-v2/internal/usecase/input"
	"server-blog-v2/internal/usecase/output"
)

type useCase struct {
	cfg     *config.Config
	emojis  repo.EmojiRepo
	sprites repo.EmojiSpriteRepo
}

// New 创建 Emoji UseCase。
func New(cfg *config.Config, emojis repo.EmojiRepo, sprites repo.EmojiSpriteRepo) usecase.Emoji {
	return &useCase{
		cfg:     cfg,
		emojis:  emojis,
		sprites: sprites,
	}
}

func (u *useCase) GetConfig(ctx context.Context) (*output.EmojiConfig, error) {
	// 获取所有活跃的 emoji
	emojis, err := u.emojis.ListActive(ctx)
	if err != nil {
		return nil, err
	}

	// 获取雪碧图信息
	sprites, err := u.sprites.ListActive(ctx)
	if err != nil {
		return nil, err
	}

	// 构建雪碧图信息
	spriteInfos := make([]output.EmojiSpriteInfo, len(sprites))
	for i, sprite := range sprites {
		// 计算该雪碧图包含的 emoji 范围
		startIdx := sprite.SpriteGroup * 128
		endIdx := startIdx + sprite.EmojiCount - 1

		// 拼接完整 URL
		url := sprite.CdnURL
		if u.cfg.Qiniu.Domain != "" && url != "" && url[0] != 'h' {
			protocol := "http://"
			if u.cfg.Qiniu.UseHTTPS {
				protocol = "https://"
			}
			url = protocol + u.cfg.Qiniu.Domain + "/" + url
		}

		spriteInfos[i] = output.EmojiSpriteInfo{
			ID:       sprite.SpriteGroup,
			Filename: sprite.Filename,
			URL:      url,
			Range:    [2]int{startIdx, endIdx},
			Frozen:   true,
			Size:     [2]int{sprite.Width / 2, sprite.Height / 2},
		}
	}

	return &output.EmojiConfig{
		Version:     fmt.Sprintf("v1.%d", time.Now().Unix()),
		TotalEmojis: int64(len(emojis)),
		Sprites:     spriteInfos,
		Mapping:     nil,
		UpdatedAt:   time.Now(),
	}, nil
}

func (u *useCase) ListGroups(ctx context.Context) ([]*output.EmojiGroup, error) {
	groups, err := u.emojis.ListGroups(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]*output.EmojiGroup, len(groups))
	for i, g := range groups {
		result[i] = &output.EmojiGroup{
			ID:          g.ID,
			GroupName:   g.GroupName,
			GroupKey:    g.GroupKey,
			Description: g.Description,
			SortOrder:   g.SortOrder,
			EmojiCount:  g.EmojiCount,
			Status:      int(g.Status),
			CreatedAt:   g.CreatedAt,
			UpdatedAt:   g.UpdatedAt,
		}
	}
	return result, nil
}

func (u *useCase) List(ctx context.Context, params input.ListEmojis) (*output.ListResult[output.EmojiInfo], error) {
	offset := (params.Page - 1) * params.PageSize

	emojis, total, err := u.emojis.List(ctx, offset, params.PageSize, params.Keyword, params.GroupKey, params.SpriteGroup)
	if err != nil {
		return nil, err
	}

	items := make([]output.EmojiInfo, len(emojis))
	for i, e := range emojis {
		// 拼接完整 URL
		cdnURL := e.CdnURL
		if u.cfg.Qiniu.Domain != "" && cdnURL != "" && cdnURL[0] != 'h' {
			protocol := "http://"
			if u.cfg.Qiniu.UseHTTPS {
				protocol = "https://"
			}
			cdnURL = protocol + u.cfg.Qiniu.Domain + "/" + cdnURL
		}

		items[i] = output.EmojiInfo{
			ID:              e.ID,
			Key:             e.Key,
			Filename:        e.Filename,
			GroupKey:        e.GroupKey,
			SpriteGroup:     e.SpriteGroup,
			SpritePositionX: e.SpritePositionX,
			SpritePositionY: e.SpritePositionY,
			FileSize:        int64(e.FileSize),
			CdnURL:          cdnURL,
			UploadTime:      e.UploadTime,
			Status:          int(e.Status),
			UpdatedAt:       e.UpdatedAt,
		}
	}

	return &output.ListResult[output.EmojiInfo]{
		Items:    items,
		Page:     params.Page,
		PageSize: params.PageSize,
		Total:    total,
	}, nil
}

func (u *useCase) ListSprites(ctx context.Context) ([]*output.SpriteInfo, error) {
	sprites, err := u.sprites.ListActive(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]*output.SpriteInfo, len(sprites))
	for i, s := range sprites {
		// 拼接完整 URL
		cdnURL := s.CdnURL
		if u.cfg.Qiniu.Domain != "" && cdnURL != "" && cdnURL[0] != 'h' {
			protocol := "http://"
			if u.cfg.Qiniu.UseHTTPS {
				protocol = "https://"
			}
			cdnURL = protocol + u.cfg.Qiniu.Domain + "/" + cdnURL
		}

		result[i] = &output.SpriteInfo{
			ID:          s.ID,
			SpriteGroup: s.SpriteGroup,
			Filename:    s.Filename,
			CdnURL:      cdnURL,
			Width:       s.Width,
			Height:      s.Height,
			EmojiCount:  s.EmojiCount,
			FileSize:    s.FileSize,
			Status:      int(s.Status),
			CreatedAt:   s.CreatedAt,
			UpdatedAt:   s.UpdatedAt,
		}
	}
	return result, nil
}

// RegenerateSprites 重新生成雪碧图。
func (u *useCase) RegenerateSprites(ctx context.Context, groupKeys []string, sseWriter usecase.SSEWriter) error {
	// 1. 获取表情列表
	emojis, err := u.emojis.ListActiveByGroupKeys(ctx, groupKeys)
	if err != nil {
		return fmt.Errorf("获取表情列表失败: %v", err)
	}

	if len(emojis) == 0 {
		return fmt.Errorf("没有活跃的表情")
	}

	// 2. 按 group_key 分组
	targetGroups := make(map[string][]*entity.Emoji)
	for _, emoji := range emojis {
		targetGroups[emoji.GroupKey] = append(targetGroups[emoji.GroupKey], emoji)
	}

	if len(targetGroups) == 0 {
		return fmt.Errorf("没有找到要处理的表情组")
	}

	// 3. 发送开始事件
	sseWriter.WriteSSEEvent("start", map[string]interface{}{
		"total_groups": len(targetGroups),
		"message":      "开始生成雪碧图...",
	})

	// 4. 创建临时目录
	tempDir := filepath.Join(os.TempDir(), fmt.Sprintf("emoji_sprites_%d", time.Now().UnixNano()))
	if err := os.MkdirAll(tempDir, 0755); err != nil {
		return fmt.Errorf("创建临时目录失败: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// 5. 为每个组生成雪碧图
	spriteResults := make(map[string][]*entity.EmojiSprite)
	totalGroups := len(targetGroups)
	processedGroups := 0

	for groupKey, groupEmojis := range targetGroups {
		processedGroups++

		// 获取 group_name
		group, _ := u.emojis.GetGroupByKey(ctx, groupKey)
		groupName := groupKey
		if group != nil {
			groupName = group.GroupName
		}

		// 发送组处理开始事件
		sseWriter.WriteSSEEvent("group_start", map[string]interface{}{
			"group_key":  groupKey,
			"group_name": groupName,
			"current":    processedGroups,
			"total":      totalGroups,
			"message":    fmt.Sprintf("正在处理表情组: %s (%d/%d)...", groupName, processedGroups, totalGroups),
		})

		// 生成雪碧图
		sprites, err := u.generateSpritesForGroup(ctx, groupKey, groupName, groupEmojis, tempDir, sseWriter)
		if err != nil {
			sseWriter.WriteSSEEvent("error", map[string]interface{}{
				"message": fmt.Sprintf("生成表情组 %s 的雪碧图失败: %v", groupName, err),
			})
			return err
		}

		spriteResults[groupKey] = sprites

		// 发送组处理完成事件
		sseWriter.WriteSSEEvent("group_complete", map[string]interface{}{
			"group_key":    groupKey,
			"group_name":   groupName,
			"current":      processedGroups,
			"total":        totalGroups,
			"sprite_count": len(sprites),
		})
	}

	sseWriter.WriteSSEEvent("progress", map[string]interface{}{
		"message": "正在清理旧数据...",
	})

	// 6. 清理旧的雪碧图记录
	if err := u.sprites.DeleteAll(ctx); err != nil {
		return fmt.Errorf("清理旧雪碧图记录失败: %v", err)
	}

	// 7. 保存新的雪碧图记录
	allSprites := make([]*entity.EmojiSprite, 0)
	for _, sprites := range spriteResults {
		allSprites = append(allSprites, sprites...)
	}

	if err := u.sprites.CreateBatch(ctx, allSprites); err != nil {
		return fmt.Errorf("保存雪碧图记录失败: %v", err)
	}

	sseWriter.WriteSSEEvent("progress", map[string]interface{}{
		"message": "正在生成配置文件...",
	})

	// 8. 为每个组生成配置文件并上传
	for groupKey, sprites := range spriteResults {
		confPath, err := u.generateAndUploadSpriteConfig(ctx, groupKey, sprites)
		if err != nil {
			return fmt.Errorf("生成配置文件失败: %v", err)
		}

		// 更新 emoji_groups 的 sprite_conf_url
		if err := u.emojis.UpdateGroupSpriteConfURL(ctx, groupKey, confPath); err != nil {
			return fmt.Errorf("更新配置URL失败: %v", err)
		}
	}

	sseWriter.WriteSSEEvent("complete", map[string]interface{}{
		"message":       "雪碧图生成完成",
		"groups_count":  len(spriteResults),
		"sprites_count": len(allSprites),
	})

	return nil
}

// generateSpritesForGroup 为一个表情组生成雪碧图。
func (u *useCase) generateSpritesForGroup(ctx context.Context, groupKey, groupName string, emojis []*entity.Emoji, tempDir string, sseWriter usecase.SSEWriter) ([]*entity.EmojiSprite, error) {
	const (
		targetSize      = 64  // 每个emoji的目标尺寸
		spritesPerRow   = 16  // 每行16个emoji
		emojisPerSprite = 128 // 每个雪碧图128个emoji
	)

	sprites := make([]*entity.EmojiSprite, 0)
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
		sseWriter.WriteSSEEvent("sprite_progress", map[string]interface{}{
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
			emojiImg, err := u.downloadAndResizeEmoji(emoji.CdnURL, targetSize)
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

		if err := png.Encode(file, spriteImg); err != nil {
			file.Close()
			return nil, fmt.Errorf("编码雪碧图失败: %v", err)
		}
		file.Close()

		// 上传到七牛云
		uploadedPath, err := u.uploadFileToQiniu(spritePath, "emoji/sprite")
		if err != nil {
			return nil, fmt.Errorf("上传雪碧图失败: %v", err)
		}

		// 获取文件大小
		fileInfo, err := os.Stat(spritePath)
		if err != nil {
			return nil, fmt.Errorf("获取雪碧图文件大小失败: %v", err)
		}

		// 创建 EmojiSprite 记录
		sprite := &entity.EmojiSprite{
			SpriteGroup: spriteSeq,
			Filename:    spriteFilename,
			CdnURL:      uploadedPath,
			Width:       spriteWidth,
			Height:      spriteHeight,
			EmojiCount:  len(spriteEmojis),
			FileSize:    int(fileInfo.Size()),
			Status:      1, // active
		}

		sprites = append(sprites, sprite)
	}

	return sprites, nil
}

// downloadAndResizeEmoji 下载并调整emoji尺寸。
func (u *useCase) downloadAndResizeEmoji(cdnURL string, targetSize int) (image.Image, error) {
	// 拼接完整URL
	fullURL := cdnURL
	if u.cfg.Qiniu.Domain != "" && cdnURL != "" && cdnURL[0] != 'h' {
		protocol := "http://"
		if u.cfg.Qiniu.UseHTTPS {
			protocol = "https://"
		}
		fullURL = protocol + u.cfg.Qiniu.Domain + "/" + cdnURL
	}

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
			srcX := bounds.Min.X + int(float64(x)/scaleX)
			srcY := bounds.Min.Y + int(float64(y)/scaleY)

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

// uploadFileToQiniu 上传文件到七牛云。
func (u *useCase) uploadFileToQiniu(filePath, pathPrefix string) (string, error) {
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
	putPolicy := storage.PutPolicy{Scope: u.cfg.Qiniu.Bucket}
	mac := qbox.NewMac(u.cfg.Qiniu.AccessKey, u.cfg.Qiniu.SecretKey)
	upToken := putPolicy.UploadToken(mac)
	cfg := &storage.Config{
		UseHTTPS:      u.cfg.Qiniu.UseHTTPS,
		UseCdnDomains: u.cfg.Qiniu.UseCDNDomains,
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

// generateAndUploadSpriteConfig 生成并上传雪碧图配置文件。
func (u *useCase) generateAndUploadSpriteConfig(ctx context.Context, groupKey string, sprites []*entity.EmojiSprite) (string, error) {
	// 获取该组的所有表情
	emojis, err := u.emojis.ListActiveByGroupKeys(ctx, []string{groupKey})
	if err != nil {
		return "", fmt.Errorf("获取表情失败: %v", err)
	}

	// 构建配置结构
	configData := map[string]interface{}{
		"group_key":  groupKey,
		"version":    fmt.Sprintf("v1.%d", time.Now().Unix()),
		"updated_at": time.Now().Format(time.RFC3339),
		"sprites":    sprites,
		"emojis":     emojis,
	}

	// 转换为JSON
	configJSON, err := json.MarshalIndent(configData, "", "  ")
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
	uploadedPath, err := u.uploadFileToQiniu(confPath, "emoji/sprite")
	if err != nil {
		return "", fmt.Errorf("上传配置文件失败: %v", err)
	}

	return uploadedPath, nil
}
