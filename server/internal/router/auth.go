package router

import (
	"server/internal/api"

	"github.com/gin-gonic/gin"
)

type AuthRouter struct{}

func (a *AuthRouter) InitAuthRouter(PublicRouter *gin.RouterGroup) {
	authPublicRouter := PublicRouter.Group("auth")
	authApi := api.ApiGroupApp.AuthApi
	{
		// 获取SSO登录URL
		authPublicRouter.GET("sso_login_url", authApi.GetSSOLoginURL)
		// SSO回调接口（后端用code换token，refresh_token存session）
		authPublicRouter.GET("callback", authApi.SSOCallback)
		// ❌ 删除了前端刷新接口，刷新由中间件自动完成！
	}
}
