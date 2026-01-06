package router

import (
	"server/internal/api"

	"github.com/gin-gonic/gin"
)

type ResourceRouter struct{}

// InitResourceRouter 初始化资源上传路由
// 所有接口都需要管理员权限
func (r *ResourceRouter) InitResourceRouter(AdminRouter *gin.RouterGroup) {
	resourceAdminRouter := AdminRouter.Group("resources")
	resourceApi := api.ApiGroupApp.ResourceApi

	{
		// 上传相关
		resourceAdminRouter.POST("check", resourceApi.Check)              // 检查文件（秒传/续传检测）
		resourceAdminRouter.POST("init", resourceApi.Init)                // 初始化上传任务
		resourceAdminRouter.POST("upload-chunk", resourceApi.UploadChunk) // 上传分片
		resourceAdminRouter.POST("complete", resourceApi.Complete)        // 完成上传
		resourceAdminRouter.POST("cancel", resourceApi.Cancel)            // 取消上传
		resourceAdminRouter.GET("progress", resourceApi.Progress)         // 查询上传进度

		// 资源管理
		resourceAdminRouter.GET("list", resourceApi.List)      // 资源列表
		resourceAdminRouter.POST("delete", resourceApi.Delete) // 删除资源
	}
}

// InitResourcePublicRouter 初始化资源公开路由（无需认证）
// 用于七牛云回调等第三方服务调用
func (r *ResourceRouter) InitResourcePublicRouter(PublicRouter *gin.RouterGroup) {
	qiniuRouter := PublicRouter.Group("qiniu")
	resourceApi := api.ApiGroupApp.ResourceApi

	{
		// 七牛云转码回调
		qiniuRouter.POST("callback", resourceApi.QiniuCallback)
	}
}
