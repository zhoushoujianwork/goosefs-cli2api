package models

type GooseFS struct {
	Bin       *string `yaml:"bin,omitempty"`
	OutputDir *string `yaml:"output_dir,omitempty"`
}
