// Package logger Zerolog 日志封装。
package logger

import (
	"fmt"
	"os"
	"strings"

	"github.com/rs/zerolog"
)

// Interface 日志接口。
type Interface interface {
	Debug(message interface{}, fields ...interface{})
	Info(message string, fields ...interface{})
	Warn(message string, fields ...interface{})
	Error(message interface{}, fields ...interface{})
	Fatal(message interface{}, fields ...interface{})
}

// Logger 基于 Zerolog 的日志实现。
type Logger struct {
	logger *zerolog.Logger
}

var _ Interface = (*Logger)(nil)

// New 创建 Logger。
func New(level string) *Logger {
	var l zerolog.Level

	switch strings.ToLower(level) {
	case "error":
		l = zerolog.ErrorLevel
	case "warn":
		l = zerolog.WarnLevel
	case "info":
		l = zerolog.InfoLevel
	case "debug":
		l = zerolog.DebugLevel
	default:
		l = zerolog.InfoLevel
	}

	zerolog.SetGlobalLevel(l)

	skipFrameCount := 3
	logger := zerolog.New(os.Stdout).With().Timestamp().CallerWithSkipFrameCount(zerolog.CallerSkipFrameCount + skipFrameCount).Logger()

	return &Logger{
		logger: &logger,
	}
}

// Debug 记录 Debug 日志。
func (l *Logger) Debug(message interface{}, fields ...interface{}) {
	l.msg("debug", message, fields...)
}

// Info 记录 Info 日志。
func (l *Logger) Info(message string, fields ...interface{}) {
	l.log(message, fields...)
}

// Warn 记录 Warn 日志。
func (l *Logger) Warn(message string, fields ...interface{}) {
	l.logger.Warn().Msgf(message, fields...)
}

// Error 记录 Error 日志。
func (l *Logger) Error(message interface{}, fields ...interface{}) {
	if l.logger.GetLevel() == zerolog.DebugLevel {
		l.Debug(message, fields...)
	}

	l.msg("error", message, fields...)
}

// Fatal 记录 Fatal 日志并退出进程。
func (l *Logger) Fatal(message interface{}, fields ...interface{}) {
	l.msg("fatal", message, fields...)

	os.Exit(1)
}

func (l *Logger) log(message string, fields ...interface{}) {
	if len(fields) == 0 {
		l.logger.Info().Msg(message)
	} else {
		l.logger.Info().Msgf(message, fields...)
	}
}

func (l *Logger) msg(level string, message interface{}, fields ...interface{}) {
	switch msg := message.(type) {
	case error:
		l.log(msg.Error(), fields...)
	case string:
		l.log(msg, fields...)
	default:
		l.log(fmt.Sprintf("%s message %v has unknown type %v", level, message, msg), fields...)
	}
}
