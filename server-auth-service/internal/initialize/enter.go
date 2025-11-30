package initialize

import "auth-service/pkg/global"

// Init 初始化所有模块
func Init() {
	global.DB = InitGorm()
	global.Redis = InitRedis()
}
