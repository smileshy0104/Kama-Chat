package config

import "time"

type KafkaConfig struct {
	MessageMode string        `mapstructure:"message_mode" json:"message_mode" yaml:"message_mode"`
	HostPort    string        `mapstructure:"host_port" json:"host_port" yaml:"host_port"`
	LoginTopic  string        `mapstructure:"login_topic" json:"login_topic" yaml:"login_topic"`
	LogoutTopic string        `mapstructure:"logout_topic" json:"logout_topic" yaml:"logout_topic"`
	ChatTopic   string        `mapstructure:"chat_topic" json:"chat_topic" yaml:"chat_topic"`
	Partition   int           `mapstructure:"partition" json:"partition" yaml:"partition"`
	Timeout     time.Duration `mapstructure:"timeout" json:"timeout" yaml:"timeout"`
}
