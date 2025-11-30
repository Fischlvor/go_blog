package main

import (
	"auth-service/internal/initialize"
	"auth-service/internal/router"
	"auth-service/pkg/core"
	"auth-service/pkg/global"
	"auth-service/pkg/jwt"
	"auth-service/scripts/flag"
	"fmt"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {
	// 1. 初始化配置
	global.Config = core.InitConfig()

	// 2. 初始化日志
	global.Log = core.InitLogger(global.Config)

	// 3. 初始化数据库
	global.DB = initialize.InitGorm()

	// 4. 初始化Redis
	global.Redis = initialize.InitRedis()

	// 5. 加载RSA密钥
	privateKey, err := jwt.LoadPrivateKey(global.Config.JWT.PrivateKeyPath)
	if err != nil {
		global.Log.Fatal("加载私钥失败", zap.Error(err))
	}
	global.RSAPrivateKey = privateKey

	publicKey, err := jwt.LoadPublicKey(global.Config.JWT.PublicKeyPath)
	if err != nil {
		global.Log.Fatal("加载公钥失败", zap.Error(err))
	}
	global.RSAPublicKey = publicKey

	// 6. 处理命令行标志（如 --sql 进行数据库迁移）
	flag.InitFlag()

	// 7. 设置Gin模式
	if global.Config.Server.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	// 8. 创建路由
	r := gin.Default()
	router.Setup(r)

	// 9. 启动服务
	addr := fmt.Sprintf(":%d", global.Config.Server.Port)
	global.Log.Info("SSO服务启动成功", zap.String("addr", addr))

	if err := r.Run(addr); err != nil {
		global.Log.Fatal("服务启动失败", zap.Error(err))
	}
}
