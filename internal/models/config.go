package models

type GooseFS struct {
	Debug         bool           `mapstructure:"debug"`
	Bin           *string        `mapstructure:"bin"`
	OutputDir     *string        `mapstructure:"output_dir"`
	DingtalkAlert *DingtalkAlert `mapstructure:"dingtalk_alert"`
}

type DingtalkAlert struct {
	Token string `mapstructure:"token"`
}
