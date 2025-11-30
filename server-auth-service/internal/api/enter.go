package api

import "auth-service/internal/service"

type ApiGroup struct {
	AuthApi
	CaptchaApi
	DeviceApi
	ManageApi
	OAuthApi
}

var ApiGroupApp = new(ApiGroup)

var authService = service.ServiceGroupApp.AuthService
var deviceService = service.ServiceGroupApp.DeviceService
var manageService = service.ServiceGroupApp.ManageService
