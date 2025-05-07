package config

type MainConfig struct {
	AppName string `json:"app_name" yaml:"app_name"`
	Host    string `json:"host" yaml:"host"`
	Port    int    `json:"port" yaml:"port"`
}
