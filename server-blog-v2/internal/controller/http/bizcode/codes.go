package bizcode

// 业务状态码定义（AABB 四位字符串格式）
// AA: 模块编号 (00=全局, 01=文章, 02=评论, 03=用户, 04=AI聊天)
// BB: 具体错误
const (
	// 成功
	Success = "0000"
	// 通用失败
	Error = "9999"

	// 全局错误码 - 数据相关 (0001-0019)
	ErrorParam        = "0001" // 参数错误
	ErrorParamMissing = "0002" // 参数缺失
	ErrorParamFormat  = "0003" // 参数格式错误
	ErrorDataNotFound = "0004" // 数据不存在
	ErrorNotFound     = "0005" // 资源不存在

	// 全局错误码 - 认证授权 (0020-0039)
	ErrorUnauthorized     = "0020" // 未授权
	ErrorTokenInvalid     = "0021" // Token无效
	ErrorTokenExpired     = "0022" // Token已过期
	ErrorPermissionDenied = "0023" // 权限不足
	ErrorLoginRequired    = "0024" // 请先登录

	// 全局错误码 - 系统错误 (0060-0079)
	ErrorSystem          = "0060" // 系统错误
	ErrorDatabase        = "0061" // 数据库错误
	ErrorCache           = "0062" // 缓存错误
	ErrorThirdParty      = "0064" // 第三方服务错误
	ErrorConfigNotLoaded = "0065" // 配置未加载

	// 文章模块 (01xx)
	ErrorArticleNotFound   = "0101"
	ErrorArticleCreateFail = "0160"
	ErrorArticleUpdateFail = "0161"
	ErrorArticleDeleteFail = "0162"

	// 评论模块 (02xx)
	ErrorCommentNotFound   = "0201"
	ErrorCommentCreateFail = "0260"
	ErrorCommentDeleteFail = "0261"

	// 用户模块 (03xx)
	ErrorUserNotFound = "0301"

	// AI聊天模块 (04xx)
	ErrorSessionNotFound   = "0401"
	ErrorSessionCreateFail = "0460"
	ErrorMessageSendFail   = "0461"
)
