package shared

import (
	"strings"

	"github.com/gofiber/fiber/v3"
)

const (
	_defaultPage     = 1
	_defaultPageSize = 10
	_maxPageSize     = 100
)

// PageMeta 分页元信息。
type PageMeta struct {
	CurrentPage int   `json:"current_page"`
	PageSize    int   `json:"page_size"`
	TotalItems  int64 `json:"total_items"`
	TotalPages  int   `json:"total_pages"`
}

// Page 分页结果。
type Page[T any] struct {
	List []T `json:"list"`
	PageMeta
}

// NewPage 创建分页结果。
func NewPage[T any](list []T, page, size int, total int64) Page[T] {
	tp := int((total + int64(size) - 1) / int64(size))
	return Page[T]{List: list, PageMeta: PageMeta{CurrentPage: page, PageSize: size, TotalItems: total, TotalPages: tp}}
}

// PageQuery 分页查询参数。
type PageQuery struct {
	Page     int               `query:"page"`
	PageSize int               `query:"page_size"`
	Keyword  string            `query:"keyword"`
	SortBy   string            `query:"sort_by"`
	Order    string            `query:"order"`
	Filters  map[string]string `query:"-"`
}

// PageQueryOption 分页查询选项。
type PageQueryOption func(*pageQueryConfig)

type pageQueryConfig struct {
	allowedSortBy  []string
	allowedFilters []string
}

// WithAllowedSortBy 设置允许的排序字段。
func WithAllowedSortBy(fields ...string) PageQueryOption {
	return func(cfg *pageQueryConfig) {
		cfg.allowedSortBy = fields
	}
}

// WithAllowedFilters 设置允许的过滤字段。
func WithAllowedFilters(keys ...string) PageQueryOption {
	return func(cfg *pageQueryConfig) {
		cfg.allowedFilters = keys
	}
}

// Normalize 规范化分页参数。
func (pq *PageQuery) Normalize() {
	if pq.Page <= 0 {
		pq.Page = _defaultPage
	}
	if pq.PageSize <= 0 {
		pq.PageSize = _defaultPageSize
	}
	if _maxPageSize > 0 && pq.PageSize > _maxPageSize {
		pq.PageSize = _maxPageSize
	}
	if pq.Filters == nil {
		pq.Filters = map[string]string{}
	}
}

// Offset 计算偏移量。
func (pq PageQuery) Offset() int { return (pq.Page - 1) * pq.PageSize }

// Limit 返回每页数量。
func (pq PageQuery) Limit() int { return pq.PageSize }

// ParsePageQuery 解析分页查询参数。
func ParsePageQuery(ctx fiber.Ctx) PageQuery { return ParsePageQueryWithOptions(ctx) }

// ParsePageQueryWithOptions 解析分页查询参数（带选项）。
func ParsePageQueryWithOptions(ctx fiber.Ctx, opts ...PageQueryOption) PageQuery {
	cfg := pageQueryConfig{}
	for _, opt := range opts {
		opt(&cfg)
	}

	page := fiber.Query[int](ctx, "page", _defaultPage)
	pageSize := fiber.Query[int](ctx, "page_size", _defaultPageSize)

	sortBy := ctx.Query("sort_by", "")
	if len(cfg.allowedSortBy) > 0 && sortBy != "" {
		ok := false
		for _, s := range cfg.allowedSortBy {
			if sortBy == s {
				ok = true
				break
			}
		}
		if !ok {
			sortBy = ""
		}
	}

	order := ctx.Query("order", "")
	if sortBy != "" {
		l := strings.ToLower(order)
		if l != "asc" && l != "desc" {
			order = "desc"
		} else {
			order = l
		}
	}

	keyword := ctx.Query("keyword", "")

	pq := PageQuery{
		Page:     page,
		PageSize: pageSize,
		Keyword:  keyword,
		SortBy:   sortBy,
		Order:    order,
		Filters:  map[string]string{},
	}

	// 解析 filter.xxx 参数
	for k, v := range ctx.Queries() {
		if strings.HasPrefix(k, "filter.") {
			key := strings.TrimPrefix(k, "filter.")
			if key != "" && v != "" {
				if len(cfg.allowedFilters) > 0 {
					allowed := false
					for _, f := range cfg.allowedFilters {
						if key == f {
							allowed = true
							break
						}
					}
					if !allowed {
						continue
					}
				}
				pq.Filters[key] = v
			}
		}
	}

	pq.Normalize()
	return pq
}
