package request

// RegisterRequest 注册请求
type RegisterRequest struct {
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=8,max=20"`
	Nickname  string `json:"nickname" binding:"required"`
	AppID     string `json:"app_id" binding:"required"`
	CaptchaID string `json:"captcha_id" binding:"required"`
	Captcha   string `json:"captcha" binding:"required,len=6"`
}

// LoginRequest 登录请求
type LoginRequest struct {
	Email            string `json:"email" binding:"required,email"`
	Password         string `json:"password"`                  // 密码登录时必填
	VerificationCode string `json:"verification_code"`         // 验证码登录时必填
	State            string `json:"state" binding:"required"`  // OAuth state参数（包含redirect_uri等）
	AppID            string `json:"app_id" binding:"required"` // 必填：应用ID
	RedirectURI      string `json:"redirect_uri"`              // 从state中提取
	DeviceID         string `json:"device_id"`                 // 从state中提取
	DeviceName       string `json:"device_name"`
	DeviceType       string `json:"device_type"`
	CaptchaID        string `json:"captcha_id"` // 密码登录时必填
	Captcha          string `json:"captcha"`    // 密码登录时必填
}

// RefreshTokenRequest 刷新Token请求
type RefreshTokenRequest struct {
	GrantType    string `json:"grant_type" binding:"required"`
	RefreshToken string `json:"refresh_token" binding:"required"`
	ClientID     string `json:"client_id" binding:"required"`
	ClientSecret string `json:"client_secret" binding:"required"`
}

// TokenExchangeRequest OAuth 2.0 token端点请求（支持authorization_code和refresh_token）
type TokenExchangeRequest struct {
	GrantType    string `json:"grant_type" binding:"required"`
	Code         string `json:"code"`          // authorization_code模式需要
	RefreshToken string `json:"refresh_token"` // refresh_token模式需要
	ClientID     string `json:"client_id" binding:"required"`
	ClientSecret string `json:"client_secret" binding:"required"`
	RedirectURI  string `json:"redirect_uri"` // authorization_code模式需要
}

// UpdatePasswordRequest 修改密码请求
type UpdatePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=8,max=20"`
}

// UpdateUserInfoRequest 更新用户信息请求
type UpdateUserInfoRequest struct {
	Nickname  string `json:"nickname"`
	Avatar    string `json:"avatar"`
	Address   string `json:"address"`
	Signature string `json:"signature"`
}

// QQLoginRequest QQ登录请求
type QQLoginRequest struct {
	Code        string `json:"code" binding:"required"`         // QQ授权码
	AppID       string `json:"app_id" binding:"required"`       // 应用ID
	RedirectURI string `json:"redirect_uri" binding:"required"` // 回调URI
	DeviceID    string `json:"device_id"`                       // 设备ID
	DeviceName  string `json:"device_name"`                     // 设备名称
	DeviceType  string `json:"device_type"`                     // 设备类型
}

// SendEmailVerificationCodeRequest 发送邮箱验证码请求
type SendEmailVerificationCodeRequest struct {
	Email     string `json:"email" binding:"required,email"`
	CaptchaID string `json:"captcha_id" binding:"required"`
	Captcha   string `json:"captcha" binding:"required,len=6"`
}

// ForgotPasswordRequest 忘记密码请求
type ForgotPasswordRequest struct {
	Email            string `json:"email" binding:"required,email"`
	VerificationCode string `json:"verification_code" binding:"required,len=6"`
	NewPassword      string `json:"new_password" binding:"required,min=8,max=20"`
}
