package models

type GooseFS struct {
	Debug         bool           `yaml:"debug,omitempty"`
	Bin           *string        `yaml:"bin,omitempty"`
	OutputDir     *string        `yaml:"output_dir,omitempty"`
	DingtalkAlert *DingtalkAlert `yaml:"dingtalk_alert,omitempty"`
}

type DingtalkAlert struct {
	Token string `yaml:"token,omitempty"`
}
