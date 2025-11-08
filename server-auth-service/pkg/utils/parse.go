package utils

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// ParseDuration 解析时间字符串（如 "24h", "7d", "30d"）
func ParseDuration(s string) (time.Duration, error) {
	if s == "" {
		return 0, errors.New("空的时间字符串")
	}

	// 处理天（d）和周（w）
	if strings.HasSuffix(s, "d") {
		days, err := strconv.Atoi(strings.TrimSuffix(s, "d"))
		if err != nil {
			return 0, fmt.Errorf("无效的天数: %w", err)
		}
		return time.Duration(days) * 24 * time.Hour, nil
	}

	if strings.HasSuffix(s, "w") {
		weeks, err := strconv.Atoi(strings.TrimSuffix(s, "w"))
		if err != nil {
			return 0, fmt.Errorf("无效的周数: %w", err)
		}
		return time.Duration(weeks) * 7 * 24 * time.Hour, nil
	}

	// 其他格式使用标准库解析（如 "2h30m"）
	return time.ParseDuration(s)
}
