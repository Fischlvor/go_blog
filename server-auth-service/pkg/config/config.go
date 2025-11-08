package config

// Config 应用配置
type Config struct {
	Captcha Captcha `yaml:"captcha"`
	JWT     JWT     `yaml:"jwt"`
	MySQL   MySQL   `yaml:"mysql"`
	Redis   Redis   `yaml:"redis"`
	QQ      QQ      `yaml:"qq"`
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

// Zap 日志配置
type Zap struct {
	Level          string `yaml:"level"`
	Filename       string `yaml:"filename"`
	MaxSize        int    `yaml:"max_size"`
	MaxBackups     int    `yaml:"max_backups"`
	MaxAge         int    `yaml:"max_age"`
	IsConsolePrint bool   `yaml:"is_console_print"`
}
