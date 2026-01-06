package request

// ResourceCheck 检查文件（秒传/续传检测）
type ResourceCheck struct {
	FileHash string `json:"file_hash" binding:"required"` // 文件MD5（前端计算）
	FileSize int64  `json:"file_size" binding:"required"` // 文件大小（字节）
	FileName string `json:"file_name" binding:"required"` // 文件名
}

// ResourceInit 初始化上传任务
type ResourceInit struct {
	FileHash string `json:"file_hash" binding:"required"` // 文件MD5
	FileSize int64  `json:"file_size" binding:"required"` // 文件大小（字节）
	FileName string `json:"file_name" binding:"required"` // 文件名
	MimeType string `json:"mime_type" binding:"required"` // MIME类型
}

// ResourceComplete 完成上传
type ResourceComplete struct {
	TaskID string `json:"task_id" binding:"required"` // 任务ID
}

// ResourceCancel 取消上传
type ResourceCancel struct {
	TaskID string `json:"task_id" binding:"required"` // 任务ID
}

// ResourceProgress 查询上传进度
type ResourceProgress struct {
	TaskID string `form:"task_id" binding:"required"` // 任务ID
}

// ResourceList 资源列表查询
type ResourceList struct {
	Page     int    `form:"page" binding:"required,min=1"`      // 页码
	PageSize int    `form:"page_size" binding:"required,min=1"` // 每页数量
	FileName string `form:"file_name"`                          // 文件名（模糊搜索）
	MimeType string `form:"mime_type"`                          // MIME类型筛选
}

// ResourceDelete 删除资源
type ResourceDelete struct {
	IDs []uint `json:"ids" binding:"required"` // 资源ID列表
}

// QiniuCallbackItem 七牛云回调中的单个处理结果
type QiniuCallbackItem struct {
	Cmd       string `json:"cmd"`       // 处理命令
	Code      int    `json:"code"`      // 状态码：0=成功
	Desc      string `json:"desc"`      // 状态描述
	Error     string `json:"error"`     // 错误信息
	Hash      string `json:"hash"`      // 输出文件hash
	Key       string `json:"key"`       // 输出文件key
	ReturnOld int    `json:"returnOld"` // 是否返回旧文件
}

// QiniuCallback 七牛云转码回调请求
type QiniuCallback struct {
	ID           string              `json:"id"`           // 任务ID (persistentId)
	Pipeline     string              `json:"pipeline"`     // 使用的队列
	Code         int                 `json:"code"`         // 整体状态码：0=成功, 1=等待, 2=处理中, 3=失败
	Desc         string              `json:"desc"`         // 状态描述
	Reqid        string              `json:"reqid"`        // 请求ID
	InputBucket  string              `json:"inputBucket"`  // 源文件bucket
	InputKey     string              `json:"inputKey"`     // 源文件key（用于匹配数据库）
	Items        []QiniuCallbackItem `json:"items"`        // 处理结果数组
	CreationDate string              `json:"creationDate"` // 任务创建时间
}
