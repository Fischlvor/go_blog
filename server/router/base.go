package router

import (
	"github.com/gin-gonic/gin"
	"server/api"
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
