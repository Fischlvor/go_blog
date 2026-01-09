package service

type ServiceGroup struct {
	AuthService
	DeviceService
	ManageService
	QQService
	ApplicationService
}

var ServiceGroupApp = new(ServiceGroup)
