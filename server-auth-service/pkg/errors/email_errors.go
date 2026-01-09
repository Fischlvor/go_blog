package errors

import "fmt"

// CooldownError 冷却错误，包含剩余秒数
type CooldownError struct {
	RemainingSeconds int
}

func (e *CooldownError) Error() string {
	return fmt.Sprintf("请等待 %d 秒后再试", e.RemainingSeconds)
}

// NewCooldownError 创建冷却错误
func NewCooldownError(remainingSeconds int) *CooldownError {
	return &CooldownError{RemainingSeconds: remainingSeconds}
}
