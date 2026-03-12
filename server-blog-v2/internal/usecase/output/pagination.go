package output

// ListResult 列表结果。
type ListResult[T any] struct {
	Items    []T
	Page     int
	PageSize int
	Total    int64
}

// AllResult 全量结果。
type AllResult[T any] struct {
	Items []T
	Total int64
}
