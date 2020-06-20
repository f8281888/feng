package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

//NodeConf ..
var NodeConf NodeConfig

//NodeConfig 配置表信息
type NodeConfig struct {
	Version        uint64 `json:"version"`
	VersionStr     string `json:"versionstr"`
	FullVersionStr string `json:"fullversionstr"`
}

//InitConfig 初始化配置
func InitConfig(configPath, pre string, value interface{}) {
	viper.AddConfigPath(configPath)
	viper.SetConfigType("json")
	configName := "config"
	if pre != "" {
		configName = pre + "-" + configName
	}

	viper.SetConfigName(configName)

	if err := viper.ReadInConfig(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if err := viper.Unmarshal(&value); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
