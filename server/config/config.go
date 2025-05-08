package config

type Config struct {
	MainConfig      MainConfig      `mapstructure:"main_config" json:"main_config" yaml:"main_config"`
	MysqlConfig     MysqlConfig     `mapstructure:"mysql_config" json:"mysql_config" yaml:"mysql_config"`
	RedisConfig     RedisConfig     `mapstructure:"redis_config" json:"redis_config" yaml:"redis_config"`
	AuthCodeConfig  AuthCodeConfig  `mapstructure:"auth_code_config" json:"auth_code_config" yaml:"auth_code_config"`
	LogConfig       LogConfig       `mapstructure:"log_config" json:"log_config" yaml:"log_config"`
	KafkaConfig     KafkaConfig     `mapstructure:"kafka_config" json:"kafka_config" yaml:"kafka_config"`
	StaticSrcConfig StaticSrcConfig `mapstructure:"static_src_config" json:"static_src_config" yaml:"static_src_config"`
}
