package config

// Config 应用配置
type Config struct {
	Captcha Captcha `yaml:"captcha"`
	JWT     JWT     `yaml:"jwt"`
	MySQL   MySQL   `yaml:"mysql"`
	Redis   Redis   `yaml:"redis"`
	QQ      QQ      `yaml:"qq"`
	Email   Email   `yaml:"email"`
	Server  Server  `yaml:"server"`
	System  System  `yaml:"system"`
	Zap     Zap     `yaml:"zap"`
}

// Captcha 验证码配置
type Captcha struct {
	Height   int     `yaml:"height"`
	Width    int     `yaml:"width"`
	Length   int     `yaml:"length"`
	MaxSkew  float64 `yaml:"max_skew"`
	DotCount int     `yaml:"dot_count"`
}

// JWT 配置
type JWT struct {
	Algorithm              string `yaml:"algorithm"`
	PrivateKeyPath         string `yaml:"private_key_path"`
	PublicKeyPath          string `yaml:"public_key_path"`
	AccessTokenExpiryTime  string `yaml:"access_token_expiry_time"`
	RefreshTokenExpiryTime string `yaml:"refresh_token_expiry_time"`
	Issuer                 string `yaml:"issuer"`
}

// MySQL 配置
type MySQL struct {
	Host         string `yaml:"host"`
	Port         int    `yaml:"port"`
	Config       string `yaml:"config"`
	DBName       string `yaml:"db_name"`
	Username     string `yaml:"username"`
	Password     string `yaml:"password"`
	MaxIdleConns int    `yaml:"max_idle_conns"`
	MaxOpenConns int    `yaml:"max_open_conns"`
	LogMode      string `yaml:"log_mode"`
}

// Redis 配置
type Redis struct {
	Address  string `yaml:"address"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
}

// QQ 第三方登录配置
type QQ struct {
	Enable      bool   `yaml:"enable"`
	AppID       string `yaml:"app_id"`
	AppKey      string `yaml:"app_key"`
	RedirectURI string `yaml:"redirect_uri"`
}

// QQLoginURL 生成QQ登录URL
func (qq QQ) QQLoginURL() string {
	return "https://graph.qq.com/oauth2.0/authorize?" +
		"response_type=code&" +
		"client_id=" + qq.AppID + "&" +
		"redirect_uri=" + qq.RedirectURI
}

// Server 服务器配置
type Server struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
	Mode string `yaml:"mode"`
}

// System 系统配置
type System struct {
	RouterPrefix    string   `yaml:"router_prefix"`
	AllowAllOrigins bool     `yaml:"allow_all_origins"`
	AllowedOrigins  []string `yaml:"allowed_origins"`
}

// Email 邮箱配置
type Email struct {
	Host     string `json:"host" yaml:"host"`         // 邮件服务器地址，例如 smtp.qq.com
	Port     int    `json:"port" yaml:"port"`         // 邮件服务器端口，常见的如 587 (TLS) 或 465 (SSL)
	From     string `json:"from" yaml:"from"`         // 发件人邮箱地址
	Nickname string `json:"nickname" yaml:"nickname"` // 发件人昵称，用于显示在邮件中的发件人信息
	Secret   string `json:"secret" yaml:"secret"`     // 发件人邮箱的密码或应用专用密码，用于身份验证
	IsSSL    bool   `json:"is_ssl" yaml:"is_ssl"`     // 是否使用 SSL 加密连接，true 表示使用，false 表示不使用
}

// Zap 日志配置
type Zap struct {
	Level          string `yaml:"level"`
	Filename       string `yaml:"filename"`
	MaxSize        int    `yaml:"max_size"`
	MaxBackups     int    `yaml:"max_backups"`
	MaxAge         int    `yaml:"max_age"`
	IsConsolePrint bool   `yaml:"is_console_print"`
}
