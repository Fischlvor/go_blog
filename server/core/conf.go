package core

import (
	"gopkg.in/yaml.v3"
	"log"
	"server/config"
	"server/utils"
)

// InitConf 从 YAML 文件加载配置
func InitConf() *config.Config {
	conf := &config.Config{}
	yamlConf, err := utils.LoadYAML()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	if err = yaml.Unmarshal(yamlConf, conf); err != nil {
		log.Fatalf("Failed to unmarshal YAML configuration: %v", err)
	}
	return conf
}
