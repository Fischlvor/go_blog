package core

import (
	"auth-service/pkg/config"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// getConfigFile 根据环境变量获取配置文件路径
func getConfigFile() string {
	env := os.Getenv("APP_ENV")
	if env == "prod" || env == "production" {
		return "configs/config.prod.yaml"
	}
	return "configs/config.yaml"
}

// InitConfig 初始化配置
func InitConfig() *config.Config {
	configFile := getConfigFile()

	// 读取配置文件
	content, err := os.ReadFile(configFile)
	if err != nil {
		panic(fmt.Sprintf("读取配置文件失败: %v", err))
	}

	// 解析配置
	var conf config.Config
	if err := yaml.Unmarshal(content, &conf); err != nil {
		panic(fmt.Sprintf("解析配置文件失败: %v", err))
	}

	return &conf
}
