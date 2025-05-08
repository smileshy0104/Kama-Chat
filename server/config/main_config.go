package config

type MainConfig struct {
	AppName string `mapstructure:"app_name" json:"app_name" yaml:"app_name"`
	Host    string `mapstructure:"host" json:"host" yaml:"host"`
	Port    int    `mapstructure:"port" json:"port" yaml:"port"`
}
