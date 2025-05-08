package config

type MysqlConfig struct {
	Host         string `mapstructure:"host" json:"host" yaml:"host"`
	Port         int    `mapstructure:"port" json:"port" yaml:"port"`
	User         string `mapstructure:"user" json:"user" yaml:"user"`
	Password     string `mapstructure:"password" json:"password" yaml:"password"`
	DatabaseName string `mapstructure:"databaseName" json:"databaseName" yaml:"databaseName"`
}
