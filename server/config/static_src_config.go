package config

type StaticSrcConfig struct {
	StaticAvatarPath string `mapstructure:"static_avatar_path" json:"static_avatar_path" yaml:"static_avatar_path"`
	StaticFilePath   string `mapstructure:"static_file_path" json:"static_file_path" yaml:"static_file_path"`
}
