package service

type ServiceGroup struct {
	AuthService
	DeviceService
	ManageService
	QQService
}

var ServiceGroupApp = new(ServiceGroup)
