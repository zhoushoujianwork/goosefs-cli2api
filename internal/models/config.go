package models

type GooseFS struct {
	Debug         bool           `mapstructure:"debug"`
	Bin           *string        `mapstructure:"bin"`
	OutputDir     *string        `mapstructure:"output_dir"` // 保留执行结果
	DingtalkAlert *DingtalkAlert `mapstructure:"dingtalk_alert"`
	DBFile        *string        `mapstructure:"db_file"` // sqlite 数据库文件
}

type DingtalkAlert struct {
	Token string `mapstructure:"token"`
}
