package router

import (
	"auth-service/internal/api"

	"github.com/gin-gonic/gin"
)

type OAuthRouter struct{}

func (o *OAuthRouter) InitOAuthRouter(apiGroup *gin.RouterGroup) {
	oauth := apiGroup.Group("/oauth")
	{
		oauthApi := api.ApiGroupApp.OAuthApi
		oauth.GET("/applications", oauthApi.GetPublicApplications)
		oauth.GET("/authorize", oauthApi.Authorize)
	}
}
