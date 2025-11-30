package api

import (
	"auth-service/internal/model/response"
	"auth-service/pkg/global"

	"github.com/gin-gonic/gin"
	"github.com/mojocn/base64Captcha"
	"go.uber.org/zap"
)

type CaptchaApi struct {
}

var captchaStore = base64Captcha.DefaultMemStore

// GetCaptcha 生成验证码
func (h *CaptchaApi) GetCaptcha(c *gin.Context) {
	// 创建数字验证码驱动
	driver := base64Captcha.NewDriverDigit(
		global.Config.Captcha.Height,
		global.Config.Captcha.Width,
		global.Config.Captcha.Length,
		global.Config.Captcha.MaxSkew,
		global.Config.Captcha.DotCount,
	)

	// 创建验证码对象
	captcha := base64Captcha.NewCaptcha(driver, captchaStore)

	// 生成验证码
	id, b64s, _, err := captcha.Generate()
	if err != nil {
		global.Log.Error("生成验证码失败", zap.Error(err))
		response.Error(c, 5001, "生成验证码失败")
		return
	}

	response.Success(c, gin.H{
		"captcha_id": id,
		"pic_path":   b64s,
	})
}

// VerifyCaptcha 验证验证码
func (h *CaptchaApi) VerifyCaptcha(captchaID, captchaCode string) bool {
	return captchaStore.Verify(captchaID, captchaCode, true)
}

// GetStore 获取store供其他handler使用
func (h *CaptchaApi) GetStore() base64Captcha.Store {
	return captchaStore
}
