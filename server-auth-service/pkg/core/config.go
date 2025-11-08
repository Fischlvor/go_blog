package core

import (
	"auth-service/pkg/config"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// InitConfig 初始化配置
func InitConfig() *config.Config {
	configFile := "configs/config.yaml"

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
