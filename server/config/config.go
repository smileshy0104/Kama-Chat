package config

type Config struct {
	MainConfig      `json:"main_config" yaml:"main_config"`
	MysqlConfig     `json:"mysql_config" yaml:"mysql_config"`
	RedisConfig     `json:"redis_config" yaml:"redis_config"`
	AuthCodeConfig  `json:"auth_code_config" yaml:"auth_code_config"`
	LogConfig       `json:"log_config" yaml:"log_config"`
	KafkaConfig     `json:"kafka_config" yaml:"kafka_config"`
	StaticSrcConfig `json:"static_src_config" yaml:"static_src_config"`
}
