package router

import (
	"auth-service/internal/api"

	"github.com/gin-gonic/gin"
)

type AuthRouter struct{}

func (a *AuthRouter) InitAuthRouter(apiGroup *gin.RouterGroup) {
	auth := apiGroup.Group("/auth")
	{
		authApi := api.ApiGroupApp.AuthApi
		auth.POST("/register", authApi.Register)
		auth.POST("/login", authApi.Login)
		auth.POST("/token", authApi.RefreshToken)
		auth.POST("/oauth/qq/login", authApi.QQLogin)
		auth.GET("/oauth/qq/callback", authApi.QQCallback)
		auth.POST("/sendEmailVerificationCode", authApi.SendEmailVerificationCode)
		auth.POST("/forgotPassword", authApi.ForgotPassword)
	}
}
