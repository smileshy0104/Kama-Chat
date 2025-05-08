package config

type LogConfig struct {
	LogPath string `mapstructure:"log_path" json:"log_path" yaml:"log_path"`
}
