package main

import (
	"server/internal/api"
	"server/internal/initialize"
	"server/pkg/core"
	"server/pkg/global"
	"server/scripts/flag"
)

func main() {
	global.Config = core.InitConf()
	global.Log = core.InitLogger()
	initialize.OtherInit() // 解析 JWT 令牌时间
	global.DB = initialize.InitGorm()
	global.Redis = initialize.ConnectRedis()
	global.ESClient = initialize.ConnectEs()

	// 初始化AI聊天服务（在数据库初始化完成后）
	api.InitAIChatService()

	defer global.Redis.Close()

	flag.InitFlag()

	initialize.InitCron()

	core.RunServer()
}
