package main

import (
	"server/core"
	"server/flag"
	"server/global"
	"server/initialize"
)

func main() {
	global.Config = core.InitConf()
	global.Log = core.InitLogger()
	initialize.OtherInit() // 解析 JWT 令牌时间
	global.DB = initialize.InitGorm()
	global.Redis = initialize.ConnectRedis()
	global.ESClient = initialize.ConnectEs()

	defer global.Redis.Close()

	flag.InitFlag()

	initialize.InitCron()

	core.RunServer()
}
