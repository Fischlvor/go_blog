package admin

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v3"

	"server-blog-v2/internal/controller/http/bizcode"
	"server-blog-v2/internal/controller/http/shared"
	"server-blog-v2/internal/usecase/input"
)

// listEmojis 表情列表。
// @Summary 表情列表（管理端）
// @Tags Admin.Emoji
// @Security BearerAuth
// @Produce json
// @Param page query int false "页码" default(1)
// @Param page_size query int false "分页大小" default(20)
// @Param keyword query string false "关键字"
// @Param group_key query string false "分组 Key"
// @Param sprite_group query int false "雪碧图分组"
// @Success 200 {object} shared.Envelope
// @Router /admin/emoji/list [get]
func (a *Admin) listEmojis(c fiber.Ctx) error {
	pq := shared.ParsePageQueryWithOptions(c, shared.WithAllowedFilters("keyword", "group_key", "sprite_group"))

	var spriteGroup *int
	if sg, ok := pq.Filters["sprite_group"]; ok && sg != "" {
		if v, err := strconv.Atoi(sg); err == nil {
			spriteGroup = &v
		}
	}

	result, err := a.emoji.List(c.Context(), input.ListEmojis{
		PageParams:  input.PageParams{Page: pq.Page, PageSize: pq.PageSize},
		Keyword:     pq.Filters["keyword"],
		GroupKey:    pq.Filters["group_key"],
		SpriteGroup: spriteGroup,
	})

	if err != nil {
		a.logger.Error(err, "http - admin - emoji - listEmojis")
		return shared.WriteError(c, http.StatusInternalServerError, bizcode.ErrorDatabase, "failed to list emojis")
	}

	return shared.WriteSuccess(c, shared.WithData(shared.NewPage(result.Items, result.Page, result.PageSize, result.Total)))
}

// listSprites 雪碧图列表。
// @Summary 雪碧图列表（管理端）
// @Tags Admin.Emoji
// @Security BearerAuth
// @Produce json
// @Success 200 {object} shared.Envelope
// @Router /admin/emoji/sprites [get]
func (a *Admin) listSprites(c fiber.Ctx) error {
	sprites, err := a.emoji.ListSprites(c.Context())
	if err != nil {
		a.logger.Error(err, "http - admin - emoji - listSprites")
		return shared.WriteError(c, http.StatusInternalServerError, bizcode.ErrorDatabase, "failed to list sprites")
	}
	return shared.WriteSuccess(c, shared.WithData(sprites))
}

// regenerateSprites 重新生成雪碧图（SSE流式响应）。
// @Summary 重新生成雪碧图
// @Description 重新生成雪碧图，支持指定表情组或全部生成
// @Tags Admin.Emoji
// @Accept json
// @Produce text/event-stream
// @Param body body object{group_keys=[]string} false "表情组keys数组，为空则生成全部"
// @Success 200 {object} shared.Envelope
// @Router /admin/emoji/regenerate [post]
func (a *Admin) regenerateSprites(c fiber.Ctx) error {
	// 解析请求体
	var req struct {
		GroupKeys []string `json:"group_keys"`
	}
	if err := c.Bind().JSON(&req); err != nil {
		// 如果没有body，则默认为空数组（生成全部）
		req.GroupKeys = []string{}
	}

	// 设置SSE响应头
	c.Set("Content-Type", "text/event-stream")
	c.Set("Cache-Control", "no-cache")
	c.Set("Connection", "keep-alive")
	c.Set("Access-Control-Allow-Origin", "*")

	// 使用 fasthttp 的流式写入
	c.RequestCtx().SetBodyStreamWriter(func(w *bufio.Writer) {
		sseWriter := &FiberSSEWriter{writer: w}

		// 执行雪碧图生成
		err := a.emoji.RegenerateSprites(context.Background(), req.GroupKeys, sseWriter)
		if err != nil {
			sseWriter.WriteSSEEvent("error", map[string]interface{}{
				"message": fmt.Sprintf("雪碧图生成失败: %v", err),
			})
		}
	})

	return nil
}

// FiberSSEWriter Fiber 的 SSE 写入器实现。
type FiberSSEWriter struct {
	writer *bufio.Writer
}

func (w *FiberSSEWriter) WriteSSEEvent(event string, data interface{}) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	// 写入 SSE 格式
	fmt.Fprintf(w.writer, "event: %s\n", event)
	fmt.Fprintf(w.writer, "data: %s\n\n", string(jsonData))
	return w.writer.Flush()
}

func (w *FiberSSEWriter) Flush() {
	w.writer.Flush()
}
