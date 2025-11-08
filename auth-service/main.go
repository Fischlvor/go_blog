package main

import (
	"auth-service/internal/initialize"
	"auth-service/internal/model/entity"
	"auth-service/internal/router"
	"auth-service/pkg/core"
	"auth-service/pkg/global"
	"auth-service/pkg/jwt"
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

	// 6. 自动迁移数据库表
	if err := global.DB.AutoMigrate(
		&entity.SSOUser{},
		&entity.SSOOAuthBinding{},
		&entity.SSOApplication{},
		&entity.UserAppRelation{},
		&entity.SSODevice{},
		&entity.SSOLoginLog{},
	); err != nil {
		global.Log.Fatal("数据库表迁移失败", zap.Error(err))
	}

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
