package input

// PageParams 分页参数。
type PageParams struct {
	Page     int
	PageSize int
}

// KeywordParams 关键字参数。
type KeywordParams struct {
	Keyword string
}

// SortParams 排序参数。
type SortParams struct {
	SortBy string
	Order  string
}

// IntFilterParam 整数过滤参数。
type IntFilterParam *int

// ParseIntFilterParam 解析整数过滤参数。
func ParseIntFilterParam(s string) IntFilterParam {
	// 简单实现，实际应该解析字符串
	return nil
}
