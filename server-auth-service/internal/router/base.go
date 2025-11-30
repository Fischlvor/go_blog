package router

import (
	"auth-service/internal/api"

	"github.com/gin-gonic/gin"
)

type BaseRouter struct{}

func (b *BaseRouter) InitBaseRouter(apiGroup *gin.RouterGroup) {
	base := apiGroup.Group("/base")
	{
		captchaApi := api.ApiGroupApp.CaptchaApi
		base.GET("/captcha", captchaApi.GetCaptcha)

		authApi := api.ApiGroupApp.AuthApi
		base.GET("/qqLoginURL", authApi.QQLoginURL)
	}
}
