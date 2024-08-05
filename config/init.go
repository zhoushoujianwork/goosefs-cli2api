package config

import (
	"goosefs-cli2api/internal/models"
	"log"
	"os"

	"github.com/alibabacloud-go/tea/tea"
	"github.com/spf13/viper"
)

var Config models.GooseFS

// 当前执行目录下的 config/config.yaml 配置文件中获取配置
func init() {
	viper.SetConfigName("config") // 配置文件名(无扩展名)
	viper.SetConfigType("yaml")   // 如果配置文件名中没有扩展名，则需要指定配置文件格式，如 "yaml"
	viper.AddConfigPath(".")      // 设置配置文件的搜索目录，当前目录
	viper.AddConfigPath("config") // 多个搜索路径，这里多加一个项目中的 config 目录

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}

	if err := viper.Unmarshal(&Config); err != nil {
		log.Fatalf("Unable to decode into struct, %v", err)
	}
	FixConfigForGoosefs()
}

func FixConfigForGoosefs() {
	if Config.OutputDir == nil {
		log.Println("Config.OutputDir is nil, use /tmp")
		Config.OutputDir = tea.String("/tmp")
	}
	// 判断目录是否存在
	if _, err := os.Stat(*Config.OutputDir); os.IsNotExist(err) {
		log.Println("OutputDir does not exist, create it")
		os.MkdirAll(*Config.OutputDir, 0755)
	}

	if Config.Bin == nil {
		log.Panicf("Config.goosefs_bin is nil")
	}

	if _, err := os.Stat(*Config.Bin); os.IsNotExist(err) {
		log.Panicf("goosefs_bin is not exist, please check %s", *Config.Bin)
	}

}
