package handler

import (
	"auth-service/pkg/global"
	"auth-service/pkg/utils"

	"github.com/gin-gonic/gin"
	"github.com/mojocn/base64Captcha"
	"go.uber.org/zap"
)

type CaptchaHandler struct {
	store base64Captcha.Store
}

func NewCaptchaHandler() *CaptchaHandler {
	return &CaptchaHandler{
		store: base64Captcha.DefaultMemStore,
	}
}

// GetCaptcha 生成验证码
func (h *CaptchaHandler) GetCaptcha(c *gin.Context) {
	// 创建数字验证码驱动
	driver := base64Captcha.NewDriverDigit(
		global.Config.Captcha.Height,
		global.Config.Captcha.Width,
		global.Config.Captcha.Length,
		global.Config.Captcha.MaxSkew,
		global.Config.Captcha.DotCount,
	)

	// 创建验证码对象
	captcha := base64Captcha.NewCaptcha(driver, h.store)

	// 生成验证码
	id, b64s, _, err := captcha.Generate()
	if err != nil {
		global.Log.Error("生成验证码失败", zap.Error(err))
		utils.Error(c, 5001, "生成验证码失败")
		return
	}

	utils.Success(c, gin.H{
		"captcha_id": id,
		"pic_path":   b64s,
	})
}

// VerifyCaptcha 验证验证码
func (h *CaptchaHandler) VerifyCaptcha(captchaID, captchaCode string) bool {
	return h.store.Verify(captchaID, captchaCode, true)
}

// GetStore 获取store供其他handler使用
func (h *CaptchaHandler) GetStore() base64Captcha.Store {
	return h.store
}
