package utils

import (
	"io/fs"
	"os"
	"server/pkg/global"

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

// LoadYAML 从文件中读取 YAML 数据并返回字节数组
func LoadYAML() ([]byte, error) {
	configFile := getConfigFile()
	return os.ReadFile(configFile)
}

// SaveYAML 将全局配置对象保存为 YAML 格式到文件
func SaveYAML() error {
	byteData, err := yaml.Marshal(global.Config)
	if err != nil {
		return err
	}
	configFile := getConfigFile()
	return os.WriteFile(configFile, byteData, fs.ModePerm)
}
