package config

// SSO SSO配置
type SSO struct {
	ServiceURL   string `json:"service_url" yaml:"service_url"`     // SSO服务后端地址
	WebURL       string `json:"web_url" yaml:"web_url"`             // SSO前端登录页地址
	ClientID     string `json:"client_id" yaml:"client_id"`         // 客户端ID
	ClientSecret string `json:"client_secret" yaml:"client_secret"` // 客户端密钥
}
