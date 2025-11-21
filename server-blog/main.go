package main

import (
	"server/internal/api"
	"server/internal/initialize"
	"server/pkg/core"
	"server/pkg/global"
	"server/pkg/utils"
	"server/scripts/flag"

	"go.uber.org/zap"
)

func main() {
	global.Config = core.InitConf()
	global.Log = core.InitLogger()
	initialize.OtherInit() // 解析 JWT 令牌时间

	// ✅ 加载SSO公钥（用于验证SSO颁发的JWT）
	if err := utils.LoadSSOPublicKey("./keys/public.pem"); err != nil {
		global.Log.Fatal("加载SSO公钥失败", zap.Error(err))
	}
	global.Log.Info("✓ SSO公钥加载成功")

	global.DB = initialize.InitGorm()
	global.Redis = initialize.ConnectRedis()
	global.ESClient = initialize.ConnectEs()
	global.RateLimiter = initialize.InitRateLimiter()

	// 初始化AI聊天服务（在数据库初始化完成后）
	api.InitAIChatService()

	defer global.Redis.Close()

	flag.InitFlag()

	initialize.InitCron()

	core.RunServer()
}
