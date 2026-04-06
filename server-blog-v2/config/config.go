package config

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/viper"
)

type (
	Config struct {
		App      App      `mapstructure:"app"`
		Log      Log      `mapstructure:"log"`
		HTTP     HTTP     `mapstructure:"http"`
		Postgres Postgres `mapstructure:"postgres"`
		Redis    Redis    `mapstructure:"redis"`
		ES       ES       `mapstructure:"elasticsearch"`
		Qiniu    Qiniu    `mapstructure:"qiniu"`
		SSO      SSO      `mapstructure:"sso"`
		AI       AI       `mapstructure:"ai"`
		Swagger  Swagger  `mapstructure:"swagger"`
		Website  Website  `mapstructure:"website"`
		Gaode    Gaode    `mapstructure:"gaode"`
		System   System   `mapstructure:"system"`
		Email    Email    `mapstructure:"email"`
		QQ       QQ       `mapstructure:"qq"`
		Jwt      Jwt      `mapstructure:"jwt"`
	}

	App struct {
		Name    string `mapstructure:"name"`
		Version string `mapstructure:"version"`
	}

	Log struct {
		Level string `mapstructure:"level"`
	}

	HTTP struct {
		Port           int  `mapstructure:"port"`
		UsePreforkMode bool `mapstructure:"use_prefork_mode"`
	}

	Postgres struct {
		Host         string `mapstructure:"host"`
		Port         int    `mapstructure:"port"`
		User         string `mapstructure:"user"`
		Password     string `mapstructure:"password"`
		DBName       string `mapstructure:"dbname"`
		SSLMode      string `mapstructure:"sslmode"`
		TimeZone     string `mapstructure:"time_zone"`
		MaxIdleConns int    `mapstructure:"max_idle_conns"`
		MaxOpenConns int    `mapstructure:"max_open_conns"`
	}

	Redis struct {
		Host     string `mapstructure:"host"`
		Port     int    `mapstructure:"port"`
		Password string `mapstructure:"password"`
		DB       int    `mapstructure:"db"`
	}

	ES struct {
		Addresses []string `mapstructure:"addresses"`
		Username  string   `mapstructure:"username"`
		Password  string   `mapstructure:"password"`
	}

	Qiniu struct {
		Zone          string `mapstructure:"zone"`
		AccessKey     string `mapstructure:"access_key"`
		SecretKey     string `mapstructure:"secret_key"`
		Bucket        string `mapstructure:"bucket"`
		Domain        string `mapstructure:"domain"`
		UseHTTPS      bool   `mapstructure:"use_https"`
		UseCDNDomains bool   `mapstructure:"use_cdn_domains"`
		PathPrefix    string `mapstructure:"path_prefix"`
	}

	SSO struct {
		ServiceURL    string `mapstructure:"service_url"`
		WebURL        string `mapstructure:"web_url"`
		ClientID      string `mapstructure:"client_id"`
		ClientSecret  string `mapstructure:"client_secret"`
		PublicKeyPath string `mapstructure:"public_key_path"`
	}

	AI struct {
		APIKey      string        `mapstructure:"api_key"`
		BaseURL     string        `mapstructure:"base_url"`
		Model       string        `mapstructure:"model"`
		MaxTokens   int           `mapstructure:"max_tokens"`
		Temperature float32       `mapstructure:"temperature"`
		Timeout     time.Duration `mapstructure:"timeout"`
	}

	Swagger struct {
		Enabled bool `mapstructure:"enabled"`
	}

	Website struct {
		Avatar      string `mapstructure:"avatar"`
		Title       string `mapstructure:"title"`
		Description string `mapstructure:"description"`
		Version     string `mapstructure:"version"`
		CreatedAt   string `mapstructure:"created_at"`
		ICPFiling   string `mapstructure:"icp_filing"`
		BilibiliURL string `mapstructure:"bilibili_url"`
		GithubURL   string `mapstructure:"github_url"`
		SteamURL    string `mapstructure:"steam_url"`
		Name        string `mapstructure:"name"`
		Job         string `mapstructure:"job"`
		Address     string `mapstructure:"address"`
		Email       string `mapstructure:"email"`
	}

	Gaode struct {
		Enable bool   `mapstructure:"enable"`
		Key    string `mapstructure:"key"`
	}

	System struct {
		UseMultipoint  bool   `mapstructure:"use_multipoint"`
		SessionsSecret string `mapstructure:"sessions_secret"`
		OssType        string `mapstructure:"oss_type"`
	}

	Email struct {
		Host     string `mapstructure:"host"`
		Port     int    `mapstructure:"port"`
		From     string `mapstructure:"from"`
		Nickname string `mapstructure:"nickname"`
		Secret   string `mapstructure:"secret"`
		IsSSL    bool   `mapstructure:"is_ssl"`
	}

	QQ struct {
		Enable      bool   `mapstructure:"enable"`
		AppID       string `mapstructure:"app_id"`
		AppKey      string `mapstructure:"app_key"`
		RedirectURI string `mapstructure:"redirect_uri"`
	}

	Jwt struct {
		AccessTokenSecret      string `mapstructure:"access_token_secret"`
		RefreshTokenSecret     string `mapstructure:"refresh_token_secret"`
		AccessTokenExpiryTime  string `mapstructure:"access_token_expiry_time"`
		RefreshTokenExpiryTime string `mapstructure:"refresh_token_expiry_time"`
		Issuer                 string `mapstructure:"issuer"`
	}
)

// getConfigFile 根据环境变量获取配置文件路径
func getConfigFile() string {
	env := os.Getenv("APP_ENV")
	if env == "prod" || env == "production" {
		return "configs/config.prod.yaml"
	}
	return "configs/config.yaml"
}

func NewConfig() (*Config, error) {
	cfg := &Config{}
	configFile := getConfigFile()
	viper.SetConfigFile(configFile)

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("config error: %w", err)
	}

	if err := viper.Unmarshal(cfg); err != nil {
		return nil, fmt.Errorf("config error: %w", err)
	}

	return cfg, nil
}
