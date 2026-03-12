package input

// ListComments 评论列表参数。
type ListComments struct {
	PageParams
}

// ListAllComments 评论列表参数（管理端）。
type ListAllComments struct {
	PageParams
	Keyword     *KeywordParams
	ArticleSlug *string
	UserUUID    *string
}

// CreateComment 创建评论参数。
type CreateComment struct {
	ArticleSlug string
	UserUUID    string
	Content     string
	ParentID    *int64
}
