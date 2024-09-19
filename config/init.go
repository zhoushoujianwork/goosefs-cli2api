package config

import (
	"goosefs-cli2api/internal/models"
	"goosefs-cli2api/pkg/db"
	"os"

	"github.com/xops-infra/noop/log"

	"github.com/alibabacloud-go/tea/tea"
	"github.com/spf13/viper"
)

var Config models.GooseFS

var DB db.DB

// 当前执行目录下的 config/config.yaml 配置文件中获取配置
func Init(debug bool) {

	if debug {
		Config.Debug = true
	}
	InitLog(Config.Debug)

	viper.SetConfigName("config") // 配置文件名(无扩展名)
	viper.SetConfigType("yaml")   // 如果配置文件名中没有扩展名，则需要指定配置文件格式，如 "yaml"
	viper.AddConfigPath(".")      // 设置配置文件的搜索目录，当前目录
	viper.AddConfigPath("/etc/goosefs-cli2api")
	viper.AddConfigPath("config") // 多个搜索路径，这里多加一个项目中的 config 目录

	if err := viper.ReadInConfig(); err != nil {
		log.Warnf("Error reading config file faild, %s, also support ENV", err)
	}

	if err := viper.Unmarshal(&Config); err != nil {
		log.Panicf("Unable to decode into struct, %v", err)
	}

	// 支持环境变量方式载入配置
	loadFromEnv()

	// 初始化数据库文件
	var dbfile string
	if Config.DBFile == nil || *Config.DBFile == "" {
		dbfile = "/opt/goosefs-cli2api/goosefs.db"
		log.Warnf("Config.db_file is nil, use %s", dbfile)
	} else {
		dbfile = *Config.DBFile
	}
	DB = db.NewSqliteDB(dbfile, Config.Debug)
}

// 优先级高于配置文件
func loadFromEnv() {
	if os.Getenv("GOOSEFS_BIN") != "" {
		if Config.Bin != nil && *Config.Bin != "" {
			log.Warnf("ENV GOOSEFS_BIN is set, ignore Config.goosefs_bin")
		}
		Config.Bin = tea.String(os.Getenv("GOOSEFS_BIN"))
		log.Infof("env set bin: %s", *Config.Bin)
	}

	if os.Getenv("GOOSEFS_OUTPUT_DIR") != "" {
		if Config.OutputDir != nil && *Config.OutputDir != "" {
			log.Warnf("ENV GOOSEFS_OUTPUT_DIR is set, ignore Config.output_dir")
		}
		Config.OutputDir = tea.String(os.Getenv("GOOSEFS_OUTPUT_DIR"))
		log.Infof("env set output dir: %s", *Config.OutputDir)
	}

	if os.Getenv("GOOSEFS_DINGTALK_ALERT_TOKEN") != "" {
		if Config.DingtalkAlert != nil && Config.DingtalkAlert.Token != "" {
			log.Warnf("ENV GOOSEFS_DINGTALK_ALERT_TOKEN is set, ignore Config.dingtalk_alert.token")
		}
		Config.DingtalkAlert = &models.DingtalkAlert{
			Token: os.Getenv("GOOSEFS_DINGTALK_ALERT_TOKEN"),
		}
		log.Infof("env set dingtalk token: %s", Config.DingtalkAlert.Token)
	}
	fixConfigForGoosefs()
}

func fixConfigForGoosefs() {
	if Config.OutputDir == nil {
		log.Warnf("Config.OutputDir is nil, use /tmp")
		Config.OutputDir = tea.String("/tmp")
	}
	// 判断目录是否存在
	if _, err := os.Stat(*Config.OutputDir); os.IsNotExist(err) {
		log.Infof("OutputDir does not exist, create it")
		os.MkdirAll(*Config.OutputDir, 0755)
	}

	if Config.Bin == nil {
		log.Panicf("Config.goosefs_bin is nil")
	}

	if _, err := os.Stat(*Config.Bin); os.IsNotExist(err) {
		log.Panicf("goosefs_bin is not exist, please check %s", *Config.Bin)
	}

}
