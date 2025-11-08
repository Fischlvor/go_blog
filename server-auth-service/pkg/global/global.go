package global

import (
	"auth-service/pkg/config"
	"crypto/rsa"

	"github.com/go-redis/redis"
	"github.com/songzhibin97/gkit/cache/local_cache"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var (
	Config        *config.Config
	Log           *zap.Logger
	DB            *gorm.DB
	Redis         *redis.Client
	BlackCache    local_cache.Cache
	RSAPrivateKey *rsa.PrivateKey // RSA私钥（签名）
	RSAPublicKey  *rsa.PublicKey  // RSA公钥（验证）
)
