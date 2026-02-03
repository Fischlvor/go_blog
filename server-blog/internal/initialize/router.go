package initialize

import (
	"net/http"
	"strconv"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"

	"server/internal/middleware"
	"server/internal/router"
	"server/pkg/global"
)

// InitRouter 初始化路由
func InitRouter() *gin.Engine {
	// 设置gin模式
	gin.SetMode(global.Config.System.Env)
	Router := gin.Default()

	// 使用日志记录中间件
	Router.Use(middleware.GinLogger(), middleware.GinRecovery(true))

	// 注意：限流中间件移到路由组级别，以支持用户级别限流

	// 使用gin会话路由 - 复用全局Redis配置参数
	// 注意：这里复用global.Config.Redis的配置，但gin-contrib/sessions/redis使用不同的Redis库
	store, err := redis.NewStoreWithDB(
		10,                                   // 连接池大小
		"tcp",                                // 网络类型
		global.Config.Redis.Address,          // 复用全局Redis地址
		global.Config.Redis.Password,         // 复用全局Redis密码
		strconv.Itoa(global.Config.Redis.DB), // 复用全局Redis数据库
		[]byte(global.Config.System.SessionsSecret), // 加密密钥
	)
	if err != nil {
		panic(err)
	}
	Router.Use(sessions.Sessions("blog_session", store))

	// 将指定目录下的文件提供给客户端
	// "uploads" 是URL路径前缀，http.Dir("uploads")是实际文件系统中存储文件的目录
	Router.StaticFS(global.Config.Upload.Path, http.Dir(global.Config.Upload.Path))

	// 健康检查接口
	Router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// 创建路由组
	routerGroup := router.RouterGroupApp

	// 公开接口：按 IP / 设备ID 限流
	publicGroup := Router.Group(global.Config.System.RouterPrefix)
	publicGroup.Use(middleware.RateLimitMiddleware(global.RateLimiter))

	// 私有接口：先 SSO JWT 认证，再按用户 UUID 限流
	privateGroup := Router.Group(global.Config.System.RouterPrefix)
	privateGroup.Use(middleware.SSOJWTAuth())
	privateGroup.Use(middleware.RateLimitMiddleware(global.RateLimiter))

	// 管理员接口：先 SSO JWT + Admin 认证，再按用户 UUID 限流
	adminGroup := Router.Group(global.Config.System.RouterPrefix)
	adminGroup.Use(middleware.SSOJWTAuth())
	adminGroup.Use(middleware.AdminAuth())
	adminGroup.Use(middleware.RateLimitMiddleware(global.RateLimiter))

	// 基础 路由
	{
		routerGroup.InitBaseRouter(publicGroup)
		routerGroup.InitAuthRouter(publicGroup)
		routerGroup.InitPublicEmojiRouter(publicGroup)
		routerGroup.InitResourcePublicRouter(publicGroup) // 七牛云回调等公开接口
	}
	// 功能 路由
	{
		routerGroup.InitUserRouter(privateGroup, publicGroup, adminGroup)
		routerGroup.InitArticleRouter(privateGroup, publicGroup, adminGroup)
		routerGroup.InitCommentRouter(privateGroup, publicGroup, adminGroup)
		routerGroup.InitFeedbackRouter(privateGroup, publicGroup, adminGroup)
		routerGroup.InitAIChatRouter(privateGroup, publicGroup, adminGroup)
	}
	// 管理员 路由
	{
		routerGroup.InitImageRouter(adminGroup)
		routerGroup.InitAdvertisementRouter(adminGroup, publicGroup)
		routerGroup.InitFriendLinkRouter(adminGroup, publicGroup)
		routerGroup.InitWebsiteRouter(adminGroup, publicGroup)
		routerGroup.InitConfigRouter(adminGroup)
		routerGroup.InitAIManagementRouter(adminGroup)
		routerGroup.InitEmojiRouter(adminGroup)
		routerGroup.InitResourceRouter(adminGroup)
	}
	return Router
}
