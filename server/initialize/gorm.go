package initialize

import (
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"os"
	"server/global"
)

func InitGorm() *gorm.DB {
	mysqlCfg := global.Config.Mysql
	cfg := &gorm.Config{
		Logger: logger.Default.LogMode(mysqlCfg.LogLevel()),
	}

	db, err := gorm.Open(mysql.Open(mysqlCfg.Dsn()), cfg)
	if err != nil {
		global.Log.Error("Failed to connect to MySQL:", zap.Error(err))
		os.Exit(1)
	}

	// 获取底层的 SQL 数据库连接对象
	sqlDB, _ := db.DB()
	// 设置数据库连接池中的最大空闲连接数
	sqlDB.SetMaxIdleConns(mysqlCfg.MaxIdleConns)
	// 设置数据库的最大打开连接数
	sqlDB.SetMaxOpenConns(mysqlCfg.MaxOpenConns)

	return db
}
