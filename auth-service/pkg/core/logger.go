package core

import (
	"auth-service/pkg/config"
	"os"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// InitLogger 初始化日志
func InitLogger(conf *config.Config) *zap.Logger {
	// 日志级别
	var level zapcore.Level
	switch conf.Zap.Level {
	case "debug":
		level = zapcore.DebugLevel
	case "info":
		level = zapcore.InfoLevel
	case "warn":
		level = zapcore.WarnLevel
	case "error":
		level = zapcore.ErrorLevel
	default:
		level = zapcore.InfoLevel
	}

	// 编码配置
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalColorLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	// 日志输出路径
	writer := &lumberjack.Logger{
		Filename:   conf.Zap.Filename,
		MaxSize:    conf.Zap.MaxSize,
		MaxBackups: conf.Zap.MaxBackups,
		MaxAge:     conf.Zap.MaxAge,
		Compress:   true,
	}

	// 创建core
	var cores []zapcore.Core

	// 文件输出
	fileCore := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		zapcore.AddSync(writer),
		level,
	)
	cores = append(cores, fileCore)

	// 控制台输出
	if conf.Zap.IsConsolePrint {
		consoleCore := zapcore.NewCore(
			zapcore.NewConsoleEncoder(encoderConfig),
			zapcore.AddSync(os.Stdout),
			level,
		)
		cores = append(cores, consoleCore)
	}

	// 创建logger
	core := zapcore.NewTee(cores...)
	logger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(0))

	return logger
}
