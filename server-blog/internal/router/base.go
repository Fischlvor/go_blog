package router

import (
	"server/internal/api"

	"github.com/gin-gonic/gin"
)

type BaseRouter struct {
}

func (b *BaseRouter) InitBaseRouter(Router *gin.RouterGroup) {
	baseRouter := Router.Group("base")

	BaseApi := api.ApiGroupApp.BaseApi

	{
		baseRouter.POST("captcha", BaseApi.Captcha)
		baseRouter.POST("sendEmailVerificationCode", BaseApi.SendEmailVerificationCode)
		baseRouter.GET("qqLoginURL", BaseApi.QQLoginURL)
		baseRouter.GET("guestWeather", BaseApi.GuestWeather)
	}

}
